import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'
import { defineConfig } from 'vite'
// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte(), tailwindcss()],
  resolve: {
    alias: {
      "@core": path.resolve(__dirname, 'src/core'),
      "@template-layout": path.resolve(__dirname, 'src/template-layout'),
      "@components": path.resolve(__dirname, 'src/components'),
      "@widgets": path.resolve(__dirname, 'src/widgets'),
      "@views": path.resolve(__dirname, 'src/views'),
    }
  }
})
