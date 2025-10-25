package origin

func ParseRuleOrigin(content []string, origin string) []string {
	switch origin {
	case "adblock":
		return AdblockRuleConvert(content)
	case "hosts":
		return HostsRuleConvert(content)
	case "adguard":
		return AdblockRuleConvert(content)
	default:
		return []string{"nil"}
	}
}
