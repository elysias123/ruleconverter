<script setup lang="ts">
import { reactive, ref, computed, onMounted, watch } from "vue";
import Icon from "./components/Icon.vue";
import Toast from "./components/Toast.vue";
import { useTheme } from "./composables/useTheme";
import { useToast } from "./composables/useToast";
import { useConverter, type ConvertParams } from "./composables/useConverter";
import { usePing } from "./composables/usePing";
import {
  ORIGIN_OPTIONS,
  TARGET_OPTIONS,
  BACKEND_OPTIONS,
  CUSTOM_BACKEND,
} from "./constants";

const { theme, init: initTheme, toggle: toggleTheme } = useTheme();
const { show: toast } = useToast();
const { resultUrl, generate } = useConverter();
const { status: pingStatus, message: pingMessage, probe, probeDebounced } =
  usePing();

// 表单状态
const form = reactive({
  rawUrl: "",
  origin: ORIGIN_OPTIONS[0].value,
  target: TARGET_OPTIONS[0].value,
});

// 后端地址：下拉选择，末项为自定义输入
const backendSelect = ref(BACKEND_OPTIONS[0].value);
const customBackend = ref("");
const isCustomBackend = computed(() => backendSelect.value === CUSTOM_BACKEND);
// 实际使用的后端地址（自定义输入不带 /rule，buildUrl 会自动补）
const effectiveBackend = computed(() =>
  isCustomBackend.value ? customBackend.value.trim() : backendSelect.value
);

// 组装传给转换器的参数
const params = (): ConvertParams => ({
  rawUrl: form.rawUrl,
  origin: form.origin,
  target: form.target,
  backend: effectiveBackend.value,
});

onMounted(() => {
  initTheme();
  probe(effectiveBackend.value); // 挂载时探测默认后端
});

// 后端变化时测试连通：下拉切换立即探测，自定义输入用防抖
watch(backendSelect, (v) => {
  if (v !== CUSTOM_BACKEND) probe(v);
  else if (customBackend.value.trim()) probeDebounced(customBackend.value);
  else {
    pingStatus.value = "idle";
    pingMessage.value = "";
  }
});
watch(customBackend, (v) => {
  if (isCustomBackend.value) probeDebounced(v);
});

function onGenerate() {
  if (isCustomBackend.value && !customBackend.value.trim())
    return toast("请填写自定义后端地址", true);
  const r = generate(params());
  if (!r.ok) return toast("请先填写规则链接", true);
  toast("订阅链接已生成");
}

async function copy(text: string) {
  if (!text) return toast("没有可复制的内容", true);
  try {
    await navigator.clipboard.writeText(text);
  } catch {
    const ta = document.createElement("textarea");
    ta.value = text;
    document.body.appendChild(ta);
    ta.select();
    document.execCommand("copy");
    document.body.removeChild(ta);
  }
  toast("已复制到剪贴板");
}
</script>

<template>
  <div class="card">
    <header class="header">
      <div class="icons">
        <a
          href="https://github.com/elysias123/ruleconverter"
          target="_blank"
          rel="noopener"
          title="GitHub"
          aria-label="GitHub"
        >
          <Icon name="github" :size="22" />
        </a>
      </div>
      <h1 class="title">规则转换</h1>
      <div class="icons"></div>
    </header>

    <div class="body">
      <div class="row">
        <label for="url">规则链接:</label>
        <div class="field">
          <textarea
            id="url"
            v-model="form.rawUrl"
            placeholder="支持 adblock / hosts / surge ruleset 等规则链接，多个链接每行一个或用 , 分隔"
          ></textarea>
        </div>
      </div>

      <div class="row">
        <label for="origin">源类型:</label>
        <div class="field">
          <select id="origin" v-model="form.origin">
            <option v-for="o in ORIGIN_OPTIONS" :key="o.value" :value="o.value">
              {{ o.label }}
            </option>
          </select>
        </div>
      </div>

      <div class="row">
        <label for="target">目标类型:</label>
        <div class="field">
          <select id="target" v-model="form.target">
            <option v-for="o in TARGET_OPTIONS" :key="o.value" :value="o.value">
              {{ o.label }}
            </option>
          </select>
        </div>
      </div>

      <div class="row">
        <label for="backend">后端地址:</label>
        <div class="field">
          <select id="backend" v-model="backendSelect">
            <option v-for="o in BACKEND_OPTIONS" :key="o.value" :value="o.value">
              {{ o.label }}
            </option>
          </select>
          <input
            v-if="isCustomBackend"
            v-model="customBackend"
            type="text"
            class="custom-backend"
            placeholder="输入自定义后端地址，如 https://your-host:8081（不带 /rule）"
          />
          <div class="ping-status" :class="pingStatus">
            <span class="ping-dot"></span>
            <span class="ping-text">{{ pingMessage || "未测试" }}</span>
            <button
              type="button"
              class="ping-retry"
              :disabled="pingStatus === 'checking'"
              @click="probe(effectiveBackend)"
            >
              重新测试
            </button>
          </div>
        </div>
      </div>

      <div class="theme-divider">
        <button
          class="theme-btn"
          title="切换主题"
          aria-label="切换主题"
          @click="toggleTheme"
        >
          <Icon :name="theme === 'dark' ? 'sun' : 'moon'" :size="16" />
        </button>
      </div>

      <div class="row">
        <label>订阅链接:</label>
        <div class="field">
          <div class="result-row">
            <input
              v-model="resultUrl"
              type="text"
              readonly
              placeholder="点击下方按钮生成转换链接"
            />
            <button class="copy-btn" @click="copy(resultUrl)">
              <Icon name="copy" :size="15" /> 复制
            </button>
          </div>
        </div>
      </div>

      <div class="actions">
        <button class="btn btn-pink btn-main" @click="onGenerate">
          生成订阅链接
        </button>
      </div>
    </div>
  </div>

  <Toast />
</template>

<style scoped>
.card {
  width: 100%;
  max-width: 1000px;
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: 16px;
  box-shadow: 0 18px 50px rgba(0, 0, 0, 0.35);
  overflow: hidden;
  backdrop-filter: blur(var(--blur)) saturate(1.4);
  -webkit-backdrop-filter: blur(var(--blur)) saturate(1.4);
}

.header {
  background: var(--header-bg);
  border-bottom: 1px solid var(--card-border);
  color: #fff;
  padding: 16px 22px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.title {
  flex: 1;
  text-align: center;
  font-family: inherit;
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 6px;
  margin: 0;
  text-shadow: 0 1px 8px rgba(0, 0, 0, 0.3);
}
.icons {
  display: flex;
  gap: 14px;
}
.icons a {
  color: #fff;
  opacity: 0.92;
  display: inline-flex;
  transition: opacity 0.2s, transform 0.2s;
}
.icons a:hover {
  opacity: 1;
  transform: translateY(-1px);
}

.body {
  padding: 26px 30px 34px;
}

.row {
  display: flex;
  align-items: flex-start;
  margin-bottom: 18px;
  gap: 16px;
}
.row > label {
  width: 86px;
  flex-shrink: 0;
  text-align: right;
  color: var(--label);
  font-size: 14px;
  font-weight: 600;
  padding-top: 11px;
}
.field {
  flex: 1;
  min-width: 0;
}

textarea,
select,
input[type="text"] {
  width: 100%;
  border: 1px solid var(--field-border);
  background: var(--field-bg);
  color: var(--field-text);
  border-radius: 8px;
  padding: 11px 13px;
  font-size: 14px;
  outline: none;
  font-family: inherit;
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  transition: border-color 0.2s, box-shadow 0.2s, background 0.2s;
}
textarea {
  min-height: 86px;
  resize: vertical;
  line-height: 1.5;
}
.custom-backend {
  margin-top: 10px;
}

/* 后端连通状态指示 */
.ping-status {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  font-size: 13px;
  color: var(--label);
}
.ping-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: #9aa0a6;
  flex-shrink: 0;
  transition: background 0.2s, box-shadow 0.2s;
}
.ping-status.checking .ping-dot {
  background: #f5b301;
  box-shadow: 0 0 0 3px rgba(245, 179, 1, 0.2);
  animation: pingPulse 1s ease-in-out infinite;
}
.ping-status.ok .ping-dot {
  background: #2ecc71;
  box-shadow: 0 0 0 3px rgba(46, 204, 113, 0.22);
}
.ping-status.fail .ping-dot {
  background: #ff5252;
  box-shadow: 0 0 0 3px rgba(255, 82, 82, 0.22);
}
.ping-status.ok .ping-text {
  color: #2ecc71;
}
.ping-status.fail .ping-text {
  color: #ff7676;
}
.ping-text {
  white-space: nowrap;
}
.ping-retry {
  margin-left: auto;
  background: var(--field-bg);
  border: 1px solid var(--field-border);
  color: var(--label);
  border-radius: 6px;
  padding: 4px 10px;
  font-size: 12px;
  cursor: pointer;
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  transition: filter 0.2s;
}
.ping-retry:hover {
  filter: brightness(1.08);
}
.ping-retry:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
@keyframes pingPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.45; }
}
select {
  appearance: none;
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='14' height='14' viewBox='0 0 24 24' fill='none' stroke='%23999' stroke-width='2'><polyline points='6 9 12 15 18 9'/></svg>");
  background-repeat: no-repeat;
  background-position: right 13px center;
  padding-right: 36px;
  cursor: pointer;
}
/* 下拉展开后的选项列表：浅底深字，避免系统深色模式下看不清 */
select option {
  background: #ffffff;
  color: #2e2740;
}
select option:checked {
  background: #2e2740;
  color: #ffffff;
}
textarea::placeholder,
input::placeholder {
  color: var(--placeholder);
}
textarea:focus,
select:focus,
input:focus {
  border-color: #4caf50;
  box-shadow: 0 0 0 3px rgba(76, 175, 80, 0.15);
}

.advanced-toggle {
  background: var(--field-bg);
  border: 1px solid var(--field-border);
  border-radius: 8px;
  padding: 11px 13px;
  text-align: center;
  color: var(--label);
  font-size: 14px;
  cursor: pointer;
  user-select: none;
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}
.advanced-panel {
  margin-top: 12px;
  padding: 14px;
  border: 1px dashed var(--field-border);
  border-radius: 8px;
  display: none;
}
.advanced-panel.open {
  display: block;
}
.adv-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  font-size: 14px;
  color: var(--label);
}
.adv-item:last-child {
  margin-bottom: 0;
}
.adv-item label {
  cursor: pointer;
}
.adv-note {
  opacity: 0.75;
  font-size: 13px;
}

.theme-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 26px 0;
}
.theme-divider::before,
.theme-divider::after {
  content: "";
  flex: 1;
  height: 1px;
  background: var(--divider);
}
.theme-btn {
  margin: 0 14px;
  width: 56px;
  height: 28px;
  border-radius: 16px;
  border: 1px solid var(--field-border);
  background: var(--field-bg);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--label);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

.result-row {
  display: flex;
  gap: 10px;
  align-items: stretch;
}
.result-row input,
.result-row textarea {
  flex: 1;
}
.copy-btn {
  background: var(--field-bg);
  border: 1px solid var(--field-border);
  color: var(--field-text);
  border-radius: 8px;
  padding: 0 16px;
  font-size: 14px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  transition: background 0.2s, border-color 0.2s;
}
.copy-btn:hover {
  background: rgba(255, 255, 255, 0.22);
}

.actions {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  margin-top: 30px;
}
.btn-group {
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
  justify-content: center;
}
.btn {
  border: none;
  border-radius: 8px;
  padding: 11px 26px;
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  transition: filter 0.2s, transform 0.1s;
}
.btn:hover {
  filter: brightness(1.05);
}
.btn:active {
  transform: translateY(1px);
}
.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.btn-pink {
  background: #f3a0a8;
}
.btn-main {
  width: 100%;
  max-width: 460px;
  justify-content: center;
}

@media (max-width: 640px) {
  .row {
    flex-direction: column;
    gap: 6px;
  }
  .row > label {
    width: auto;
    text-align: left;
    padding-top: 0;
  }
}
</style>
