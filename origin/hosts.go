package origin

import (
	"net/url"
	"strings"
)

func OriginHostsRule(content []string) []string {
	var domains []string
	seen := make(map[string]struct{})

	for _, line := range content {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
			if line == "" {
				continue
			}
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		for i := 1; i < len(parts); i++ {
			host := strings.TrimSpace(parts[i])
			if host == "" {
				continue
			}
			if host == "0.0.0.0" || host == "127.0.0.1" || host == "::1" {
				continue
			}
			if u, err := url.Parse("http://" + host); err == nil && u.Host != "" {

				h := u.Host
				if strings.Contains(h, ":") {
					h = strings.Split(h, ":")[0]
				}
				if _, ok := seen[h]; !ok {
					seen[h] = struct{}{}
					domains = append(domains, h)
				}
			}
		}
	}

	return domains
}
