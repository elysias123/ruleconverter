package target

import "strings"

func MihomoDomainConvert(content []string) []string {

	var list []string
	list = append(list, "payload:")
	for _, line := range content {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = "  - '" + line + "'"
		list = append(list, line)
	}
	return list
}
