import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 前后端分离：前端作为独立静态站点构建，输出到 dist/
export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: 'dist',
  },
  server: {
    // 开发时把 /rule、/ping 代理到本地 go 后端
    proxy: {
      '/rule': 'http://localhost:30000',
      '/ping': 'http://localhost:30000',
    },
  },
})
