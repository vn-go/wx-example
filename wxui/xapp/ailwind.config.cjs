// tailwind.config.cjs
/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        './src/**/*.{html,js,svelte,ts}', // Đảm bảo có .svelte
    ],
    theme: {
        extend: {},
    },
    plugins: [],
}