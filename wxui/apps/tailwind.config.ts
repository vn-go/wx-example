// tailwind.config.ts
import type { Config } from "tailwindcss";

export default {
    content: [
        "./src/**/*.{html,js,svelte,ts}",
        "./src/app.css"        // ← QUAN TRỌNG NHẤT ĐÂY
    ],
} satisfies Config;