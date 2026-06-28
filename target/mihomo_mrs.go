package target

import (
	"bytes"
	"log"
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
		log.Printf("ConvertToMrs failed: %v", err)
		return nil
	}

	out := strings.Split(buf.String(), "\n")
	return out
}
