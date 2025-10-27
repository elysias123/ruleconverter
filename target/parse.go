package target

func ParseRuleTarget(content []string, target string) []string {
	switch target {
	case "mihomo":
		return TargetMihomoDomain(content)
	case "mihomo_domain":
		return TargetMihomoDomain(content)
	case "mihomo_mrs":
		return TargetMihomoMrs(content)
	case "surge":
		return TargetSurgeModule(content)
	case "surge_module":
		return TargetSurgeModule(content)
	case "surge_ruleset":
		return TargetSurgeRuleset(content)
	default:
		return []string{"nil"}
	}
}
