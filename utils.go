package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"ruleconverter/origin"
	"ruleconverter/target"
	"runtime"
	"syscall"
	"time"
)

// 私有 / 保留地址段，禁止服务端向其发起请求，防止 SSRF
var blockedCIDRs = func() []*net.IPNet {
	cidrs := []string{
		"0.0.0.0/8",          // 当前网络
		"10.0.0.0/8",         // 私有
		"100.64.0.0/10",      // CGNAT
		"127.0.0.0/8",        // 回环
		"169.254.0.0/16",     // 链路本地（含云元数据 169.254.169.254）
		"172.16.0.0/12",      // 私有
		"192.0.0.0/24",       // IETF 协议分配
		"192.0.2.0/24",       // TEST-NET-1
		"192.168.0.0/16",     // 私有
		"198.18.0.0/15",      // 基准测试
		"198.51.100.0/24",    // TEST-NET-2
		"203.0.113.0/24",     // TEST-NET-3
		"240.0.0.0/4",        // 保留
		"255.255.255.255/32", // 广播
		"::1/128",            // IPv6 回环
		"::/128",             // 未指定地址
		"64:ff9b::/96",       // IPv4/IPv6 转换
		"100::/64",           // 丢弃前缀
		"fc00::/7",           // 唯一本地地址
		"fe80::/10",          // 链路本地
	}
	nets := make([]*net.IPNet, 0, len(cidrs))
	for _, c := range cidrs {
		if _, n, err := net.ParseCIDR(c); err == nil {
			nets = append(nets, n)
		}
	}
	return nets
}()

// isBlockedIP 判断 IP 是否落在禁止访问的地址段
func isBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	for _, n := range blockedCIDRs {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// validateConnAddr 是 dialer 的 Control 钩子使用的校验函数。
// 它在实际 connect 之前对已解析出的目标 IP 做检查，
// 因此能同时防御 DNS rebinding 与跳转/重定向到内网。
func validateConnAddr(address string) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid address %q: %w", address, err)
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("address %q is not a literal IP", address)
	}
	if isBlockedIP(ip) {
		return fmt.Errorf("access to non-public address %s is blocked", ip)
	}
	return nil
}

// validateRequestURL 在发起请求前校验 URL 的 scheme，
// 仅允许 http / https，拒绝 file:// gopher:// 等危险协议。
func validateRequestURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	switch u.Scheme {
	case "http", "https":
	default:
		return fmt.Errorf("unsupported URL scheme %q", u.Scheme)
	}
	if u.Hostname() == "" {
		return fmt.Errorf("URL missing host")
	}
	return nil
}

func unique[T comparable](items []T) []T {
	seen := make(map[T]struct{}, len(items))
	out := make([]T, 0, len(items))
	for _, v := range items {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// maxResponseBytes 限制单个 URL 响应体大小，防止内存放大型 DoS
const maxResponseBytes = 16 * 1024 * 1024 // 16 MiB

func GetUrlContent(url string) ([]string, error) {
	// 先做 scheme / host 基础校验，拒绝 file:// 等危险协议
	if err := validateRequestURL(url); err != nil {
		return nil, err
	}

	// Control 钩子在 connect 之前对已解析出的目标 IP 做校验，
	// 能同时防御 DNS rebinding 与重定向/跳转到内网（每跳都会重新拨号）。
	control := func(network, address string, c syscall.RawConn) error {
		return validateConnAddr(address)
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
			Control:   control,
		}).DialContext,
	}

	// 在 Android 上使用自定义解析器（例如 223.5.5.5）
	if runtime.GOOS == "android" {
		resolver := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := &net.Dialer{Timeout: 3 * time.Second}
				return d.DialContext(ctx, "udp", "223.5.5.5:53")
			},
		}
		transport.DialContext = (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
			Resolver:  resolver,
			Control:   control,
		}).DialContext
	}

	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: transport,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Println("http get error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("bad status:", resp.Status)
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// 限制读取大小，超出即视为异常
	limited := io.LimitReader(resp.Body, maxResponseBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		log.Println("read body error:", err)
		return nil, err
	}
	if int64(len(body)) > maxResponseBytes {
		return nil, fmt.Errorf("response body exceeds %d bytes limit", maxResponseBytes)
	}

	scanner := bufio.NewScanner(bytes.NewReader(body))
	// 如需处理超长行，可增大容量
	buf := make([]byte, 0, 64*1024)
	const maxCapacity = 10 * 1024 * 1024
	scanner.Buffer(buf, maxCapacity)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println("scan error:", err)
		return nil, err
	}
	return lines, nil
}

func RuleConverter(content []string, originType string, targetType string) []string {

	content = unique(content) // 去重
	domains := origin.ParseRuleOrigin(content, originType)

	ret := target.ParseRuleTarget(domains, targetType)
	return ret
}
