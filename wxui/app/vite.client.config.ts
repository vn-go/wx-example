// src/vite/client.config.ts
import { defineConfig } from 'vitest/config';
import { baseConfig } from './vite.base.config';

export const clientTestConfig = defineConfig({
	...baseConfig,
	test: {
		name: 'client',
		environment: 'browser',
		expect: { requireAssertions: true },
		browser: {
			enabled: true,
			provider: 'playwright',
			instances: [{ browser: 'chromium' }]
		},
		include: ['src/**/*.svelte.{test,spec}.{js,ts}'],
		exclude: ['src/lib/server/**'],
		setupFiles: ['./vitest-setup-client.ts']
	}
});