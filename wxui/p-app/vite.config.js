// import preact from '@preact/preset-vite'
// import { defineConfig } from 'vite'

// // https://vite.dev/config/
// export default defineConfig({
//   plugins: [preact(), tailwindcss(),],
// })
import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'vite'
export default defineConfig({
  plugins: [preact(), tailwindcss()],

})
