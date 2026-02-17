import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [sveltekit(), tailwindcss()],
	server: {
		proxy: {
			'/auth': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
			},
			'/projects': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
			},
			'/boards': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
			},
			'/tasks': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
			},
			'/epics': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
			},
			'/members': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
			},
			'/wiki': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
			},
			'/files': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
			},
			'/admin': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
			},
			'/health': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
			},
			'/ready': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
			},
			'/labels': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				cookieDomainRewrite: 'localhost'
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
