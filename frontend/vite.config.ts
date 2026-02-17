import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [sveltekit(), tailwindcss()],
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, '')
			},
			'/auth': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/projects': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/boards': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/tasks': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/epics': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/members': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/wiki': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/files': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/admin': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/health': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/ready': {
				target: 'http://localhost:8080',
				changeOrigin: true
			}
		}
	},
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}', 'tests/**/*.{test,spec}.{js,ts}'],
		environment: 'jsdom',
		globals: true,
		setupFiles: ['./tests/setup.ts']
	}
});
