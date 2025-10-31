// src/vite/base.config.ts
import { sveltekit } from '@sveltejs/kit/vite';
import path from 'path';
import { defineConfig } from 'vite';

export const baseConfig = defineConfig({
    plugins: [sveltekit()],
    resolve: {
        alias: {
            '@components': path.resolve('./src/lib/components'),
            '@lib': path.resolve('./src/lib'),
            '@routes': path.resolve('./src/routes'),
            '@store': path.resolve('./src/lib/store'),
            '@utils': path.resolve('./src/lib/utils'),
        }
    }
});