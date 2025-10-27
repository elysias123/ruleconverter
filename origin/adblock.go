package origin

import (
	"net"
	"net/url"
	"strings"
)

func OriginAdblockRule(content []string) []string {
	var domains []string
	seen := make(map[string]struct{})

	for _, line := range content {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "!") || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "##") {
			continue
		}

		if strings.HasPrefix(line, "@@") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			continue
		}

		if strings.HasPrefix(line, "/") && strings.HasSuffix(line, "/") {
			continue
		}

		var domainToProcess string

		// ||example.com^
		if strings.HasPrefix(line, "||") {
			domainToProcess = line[2:]
		} else if parts := strings.Fields(line); len(parts) == 2 && (parts[0] == "0.0.0.0" || parts[0] == "127.0.0.1" || parts[0] == "::1") {
			// 0.0.0.0 example.com or 127.0.0.1 example.com or ::1 example.com
			domainToProcess = parts[1]
		}

		if domainToProcess != "" {
			separatorIndex := strings.IndexAny(domainToProcess, "^/$")
			if separatorIndex != -1 {
				domainToProcess = domainToProcess[:separatorIndex]
			}

			if u, err := url.Parse("http://" + domainToProcess); err == nil && u.Host != "" {
				host := u.Host
				if h, _, err := net.SplitHostPort(host); err == nil {
					host = h
				}

				host = strings.ToLower(strings.TrimSpace(host))

				if net.ParseIP(host) != nil {
					continue
				}

				if !isValidDomain(host) {
					continue
				}

				if _, ok := seen[host]; !ok {
					seen[host] = struct{}{}
					domains = append(domains, host)
				}
			}
		}
	}

	return domains
}

// isValidDomain 简单校验域名格式：长度限制、标签规则、TLD 不全为数字
func isValidDomain(host string) bool {
	if host == "" {
		return false
	}
	if len(host) > 253 {
		return false
	}

	labels := strings.Split(host, ".")
	if len(labels) < 2 {
		// 要求至少有一个点来判断为域名（如 example.com）
		return false
	}

	for _, label := range labels {
		l := len(label)
		if l == 0 || l > 63 {
			return false
		}
		// 首尾不能是 '-'，且只能包含字母数字和 '-'
		for i, ch := range label {
			if !(ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == '-') {
				// 这里假定 host 已经是小写
				return false
			}
			if (i == 0 || i == l-1) && ch == '-' {
				return false
			}
		}
	}

	// TLD 不能全部是数字
	last := labels[len(labels)-1]
	allDigit := true
	for _, ch := range last {
		if ch < '0' || ch > '9' {
			allDigit = false
			break
		}
	}
	return !allDigit
}
