import { ref } from "vue";

const message = ref("");
const visible = ref(false);
const isError = ref(false);
let timer: number | undefined;

export function useToast() {
  const show = (msg: string, err = false) => {
    message.value = msg;
    isError.value = err;
    visible.value = true;
    clearTimeout(timer);
    timer = window.setTimeout(() => (visible.value = false), 2200);
  };
  return { message, visible, isError, show };
}
