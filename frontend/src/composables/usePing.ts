import { ref } from "vue";

// 后端连通状态
export type PingStatus = "idle" | "checking" | "ok" | "fail";

// 计算用于 ping 的 base：留空则用当前站点
function resolveBase(backend: string): string {
  return (backend.trim() || location.origin).replace(/\/+$/, "");
}

export function usePing() {
  const status = ref<PingStatus>("idle");
  const message = ref("");
  let timer: number | undefined;
  let seq = 0; // 防止旧请求覆盖新结果

  // 立即探测一次
  async function probe(backend: string) {
    const base = resolveBase(backend);
    if (!base) {
      status.value = "idle";
      message.value = "";
      return;
    }
    const my = ++seq;
    status.value = "checking";
    message.value = "测试中…";

    const ctrl = new AbortController();
    const to = window.setTimeout(() => ctrl.abort(), 6000);
    try {
      const resp = await fetch(`${base}/ping`, { signal: ctrl.signal });
      const text = (await resp.text()).trim();
      if (my !== seq) return; // 已有更新的探测，丢弃本次结果
      if (resp.ok && text === "pong") {
        status.value = "ok";
        message.value = "后端连通";
      } else {
        status.value = "fail";
        message.value = `异常响应：${resp.status}`;
      }
    } catch (e) {
      if (my !== seq) return;
      status.value = "fail";
      message.value =
        (e as Error).name === "AbortError" ? "连接超时" : "无法连接";
    } finally {
      window.clearTimeout(to);
    }
  }

  // 防抖探测，供输入框 watch 调用
  function probeDebounced(backend: string, delay = 500) {
    status.value = "checking";
    message.value = "";
    window.clearTimeout(timer);
    timer = window.setTimeout(() => probe(backend), delay);
  }

  return { status, message, probe, probeDebounced };
}
