import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import path from 'path';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	resolve: {
		alias: {
			'@components': path.resolve('./src/lib/components'),
			'@lib': path.resolve('./src/lib'),
			'@routes': path.resolve('./src/routes'),
			'@store': path.resolve('./src/lib/store'),
			'@utils': path.resolve('./src/lib/utils'),
			'@layouts': path.resolve('./src/layouts'),
		}
	}

});
