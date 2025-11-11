import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { resolveAlias } from "./vite.resolveAlias.config";
export default defineConfig({
	plugins: [

		sveltekit()
	],
	resolve: {
		alias: resolveAlias,
	}
});