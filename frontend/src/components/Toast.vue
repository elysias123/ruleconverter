<script setup lang="ts">
import { useToast } from "../composables/useToast";

const { message, visible, isError } = useToast();
</script>

<template>
  <div class="toast" :class="{ show: visible, err: isError }">{{ message }}</div>
</template>

<style scoped>
.toast {
  position: fixed;
  bottom: 30px;
  left: 50%;
  transform: translateX(-50%) translateY(20px);
  background: rgba(30, 30, 30, 0.92);
  color: #fff;
  padding: 11px 22px;
  border-radius: 8px;
  font-size: 14px;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.25s, transform 0.25s;
  z-index: 999;
}
.toast.show {
  opacity: 1;
  transform: translateX(-50%) translateY(0);
}
.toast.err {
  background: rgba(200, 50, 50, 0.95);
}
</style>
