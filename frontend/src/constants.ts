// 源类型与目标类型选项，集中维护
export interface SelectOption {
  value: string;
  label: string;
}

export const ORIGIN_OPTIONS: SelectOption[] = [
  { value: "adblock", label: "Adblock / AdGuard" },
  { value: "hosts", label: "Hosts" },
  { value: "surge_ruleset", label: "Surge Ruleset" },
];

export const TARGET_OPTIONS: SelectOption[] = [
  { value: "mihomo", label: "Mihomo Domain (yaml)" },
  { value: "mihomo_mrs", label: "Mihomo MRS" },
  { value: "surge_module", label: "Surge Module" },
  { value: "surge_ruleset", label: "Surge Ruleset" },
];

// 自定义输入的标记值
export const CUSTOM_BACKEND = "__custom__";

// 后端地址：首项为本地默认，末项为自定义输入
export const BACKEND_OPTIONS: SelectOption[] = [
  { value: "http://localhost:30000", label: "http://localhost:30000" },
  { value: CUSTOM_BACKEND, label: "自定义输入" },
];

// 各目标类型对应的下载文件扩展名
export const EXT_MAP: Record<string, string> = {
  mihomo: "yaml",
  mihomo_domain: "yaml",
  mihomo_mrs: "txt",
  surge_module: "sgmodule",
  surge: "sgmodule",
  surge_ruleset: "list",
};

export const HELP_TEXT =
  "接口: GET /rule?target=&origin=&url=\n\n" +
  "origin 源类型: adblock / adguard / hosts / surge_ruleset\n" +
  "target 目标类型: mihomo(_domain) / mihomo_mrs / surge(_module) / surge_ruleset\n\n" +
  "多个 url 用英文逗号分隔。";
