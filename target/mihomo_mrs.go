package target

import (
	"bytes"
	"strings"

	ruleprovider "github.com/metacubex/mihomo/constant/provider"
	rulesprovider "github.com/metacubex/mihomo/rules/provider"
)

func TargetMihomoMrs(content []string) []string {
	domain := TargetMihomoDomain(content)

	var buf bytes.Buffer
	err := rulesprovider.ConvertToMrs(
		[]byte(strings.Join(domain, "\n")),
		ruleprovider.Domain,
		ruleprovider.TextRule,
		&buf,
	)
	if err != nil {
		panic(err)
	}

	out := strings.Split(buf.String(), "\n")
	return out
}
