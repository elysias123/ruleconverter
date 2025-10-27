package origin

func ParseRuleOrigin(content []string, origin string) []string {
	switch origin {
	case "adblock":
		return OriginAdblockRule(content)
	case "hosts":
		return OriginHostsRule(content)
	case "adguard":
		return OriginAdblockRule(content)
	default:
		return []string{"nil"}
	}
}
