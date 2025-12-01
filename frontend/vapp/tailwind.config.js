/** @type {import('tailwindcss').Config} */
import PrimeUI from 'tailwindcss-primeui';
//'./node_modules/primevue/**/*.{vue,js,ts}'
module.exports = {
    content: [
        "./index.html",
        "./src/**/*.{vue,js,ts}",
        "./node_modules/primevue/**/*.{vue,js,ts}"
    ],
    theme: {
        extend: {
            fontFamily: {
                sans: ['Inter', 'ui-sans-serif', 'system-ui'],
            },
        },
    },
    plugins: [
        PrimeUI,
        require('tailwindcss-primeui')
    ],
}
