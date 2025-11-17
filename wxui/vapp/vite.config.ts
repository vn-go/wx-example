import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],

  resolve: {
    alias: {
      "@core": path.resolve(__dirname, 'src/core'),
      "@app-layouts": path.resolve(__dirname, 'src/app-layouts'),
      "@widgets": path.resolve(__dirname, 'src/widgets'),
      "@views": path.resolve(__dirname, 'src/views')
    }
  }
})
