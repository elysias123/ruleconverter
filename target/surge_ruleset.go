package target

import "strings"

func TargetSurgeRuleset(content []string) []string {

	var list []string
	for _, line := range content {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = "DOMAIN," + line
		list = append(list, line)
	}
	return list
}
