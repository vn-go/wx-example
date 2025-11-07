import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const config = {
    preprocess: vitePreprocess({
        script: true, // cho phép xử lý <script lang="ts"> có decorators

    }),

    kit: {
        adapter: adapter({
            fallback: 'index.html' // 👈 cần thiết cho SPA
        })
    }
};

export default config;
