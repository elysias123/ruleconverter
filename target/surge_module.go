package target

import "strings"

func TargetSurgeModule(content []string) []string {

	var list []string
	list = append(list, "[Rule]")
	for _, line := range content {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = "DOMAIN," + line + ",REJECT,extended-matching,pre-matching"
		list = append(list, line)
	}
	return list
}
