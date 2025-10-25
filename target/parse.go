package target

func ParseRuleTarget(content []string, target string) []string {
	switch target {
	case "mihomo":
		return MihomoDomainConvert(content)
	case "mihomo_domain":
		return MihomoDomainConvert(content)
	case "mihomo_mrs":
		return MihomoMrsConvert(content)
	default:
		return []string{"nil"}
	}
}
