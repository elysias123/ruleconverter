import { ref } from "vue";

export type Theme = "light" | "dark";

const theme = ref<Theme>("light");

function apply(t: Theme) {
  theme.value = t;
  document.documentElement.setAttribute("data-theme", t);
  localStorage.setItem("rc-theme", t);
}

export function useTheme() {
  const init = () => apply((localStorage.getItem("rc-theme") as Theme) || "light");
  const toggle = () => apply(theme.value === "dark" ? "light" : "dark");
  return { theme, init, toggle };
}
