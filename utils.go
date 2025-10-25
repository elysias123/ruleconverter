package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"ruleconverter/origin"
	"ruleconverter/target"
	"runtime"
	"time"
)

func GetUrlContent(url string) []string {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// 在 Android 上使用自定义解析器（例如 223.5.5.5），先解析主机名为 IP 再连接
	if runtime.GOOS == "android" {
		resolver := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := &net.Dialer{Timeout: 3 * time.Second}
				return d.DialContext(ctx, "udp", "223.5.5.5:53")
			},
		}

		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					// 无法拆分，回退到默认拨号
					d := &net.Dialer{Timeout: 10 * time.Second}
					return d.DialContext(ctx, network, addr)
				}

				ips, err := resolver.LookupIPAddr(ctx, host)
				if err != nil || len(ips) == 0 {
					// 解析失败，回退到默认拨号
					d := &net.Dialer{Timeout: 10 * time.Second}
					return d.DialContext(ctx, network, addr)
				}

				d := &net.Dialer{Timeout: 10 * time.Second}
				// 依次尝试解析到的 IP
				var lastErr error
				for _, ip := range ips {
					ipAddr := net.JoinHostPort(ip.IP.String(), port)
					conn, err := d.DialContext(ctx, network, ipAddr)
					if err == nil {
						return conn, nil
					}
					lastErr = err
				}
				if lastErr == nil {
					lastErr = err
				}
				return nil, lastErr
			},
		}
		client.Transport = transport
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Println("http get error:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("bad status:", resp.Status)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("read body error:", err)
		return nil
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
		return nil
	}
	return lines
}

func RuleConverter(content []string, originType string, targetType string) []string {
	domains := origin.ParseRuleOrigin(content, originType)

	ret := target.ParseRuleTarget(domains, targetType)
	return ret
}
