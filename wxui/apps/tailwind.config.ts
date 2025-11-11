// tailwind.config.ts
import forms from '@tailwindcss/forms';
import typography from '@tailwindcss/typography';
import type { Config } from "tailwindcss";
export default {
    content: [
        "./src/**/*.{html,js,svelte,ts}",
        "./src/app.css"        // ← QUAN TRỌNG NHẤT ĐÂY
    ],
    theme: {
        extend: {},
    },
    plugins: [forms, typography],
} satisfies Config;