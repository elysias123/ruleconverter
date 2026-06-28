import { ref } from "vue";
import { EXT_MAP } from "../constants";

export interface ConvertParams {
  rawUrl: string;
  origin: string;
  target: string;
  backend: string;
}

// 把多行 / 逗号分隔的输入规整为逗号分隔的单串
function normalizeUrls(raw: string): string[] {
  return raw
    .split(/[\n,]+/)
    .map((s) => s.trim())
    .filter(Boolean);
}

// 根据参数构造后端转换链接
export function buildUrl(p: ConvertParams): string | null {
  const urls = normalizeUrls(p.rawUrl);
  if (urls.length === 0) return null;
  const base = (p.backend.trim() || location.origin).replace(/\/+$/, "");
  const params = new URLSearchParams({
    target: p.target,
    origin: p.origin,
    url: urls.join(","),
  });
  return `${base}/rule?${params.toString()}`;
}

export function useConverter() {
  const resultUrl = ref("");
  const resultText = ref("");
  const loading = ref(false);

  // 生成订阅链接（不发请求）
  const generate = (p: ConvertParams): { ok: boolean; url?: string } => {
    const url = buildUrl(p);
    if (!url) return { ok: false };
    resultUrl.value = url;
    return { ok: true, url };
  };

  // 预览：拉取后端转换结果
  const preview = async (
    p: ConvertParams
  ): Promise<{ ok: boolean; lines?: number; status?: number; msg?: string }> => {
    const url = buildUrl(p);
    if (!url) return { ok: false, msg: "empty" };
    resultUrl.value = url;
    loading.value = true;
    try {
      const resp = await fetch(url);
      const text = await resp.text();
      resultText.value = text;
      if (!resp.ok) return { ok: false, status: resp.status };
      const lines = text.split("\n").filter(Boolean).length;
      return { ok: true, lines };
    } catch (e) {
      return { ok: false, msg: (e as Error).message };
    } finally {
      loading.value = false;
    }
  };

  // 下载当前结果
  const download = (target: string): boolean => {
    if (!resultText.value) return false;
    const ext = EXT_MAP[target] || "txt";
    const blob = new Blob([resultText.value], { type: "text/plain;charset=utf-8" });
    const a = document.createElement("a");
    a.href = URL.createObjectURL(blob);
    a.download = `ruleset.${ext}`;
    a.click();
    URL.revokeObjectURL(a.href);
    return true;
  };

  const clear = () => {
    resultUrl.value = "";
    resultText.value = "";
  };

  return { resultUrl, resultText, loading, generate, preview, download, clear };
}
