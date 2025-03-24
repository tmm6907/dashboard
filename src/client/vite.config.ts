import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
	plugins: [
		sveltekit(),
		tailwindcss(),
	],
	server: {
		port: 3030,
		proxy: {
			'/api': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				secure: false,
			}
		},
		watch: {
			usePolling: true,
		}
	},
});
