package target

func ParseRuleTarget(content []string, target string) []string {
	switch target {
	case "mihomo":
		return MihomoDomainConvert(content)
	case "mihomo_domain":
		return MihomoDomainConvert(content)
	case "mihomo_mrs":
		return MihomoMrsConvert(content)
	case "surge":
		return SurgeModuleConvert(content)
	case "surge_module":
		return SurgeModuleConvert(content)
	case "surge_ruleset":
		return SurgeRulesetConvert(content)
	default:
		return []string{"nil"}
	}
}
