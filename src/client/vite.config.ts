import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite'
import { VitePWA } from 'vite-plugin-pwa';


export default defineConfig({
	plugins: [
		sveltekit(),
		tailwindcss(),
		VitePWA({
			srcDir: 'src',
			filename: 'sw.js',
			registerType: 'autoUpdate',
			manifest: {
				name: 'My Svelte PWA',
				short_name: 'SveltePWA',
				description: 'A Progressive Web App built with Svelte',
				theme_color: '#ffffff',
				background_color: '#ffffff',
				display: 'standalone',
				icons: [
					{
						src: "/public/legislation (6).svg",
						sizes: "192x192",
						type: "image/svg+xml"
					},
				]
			},
			workbox: {
				globPatterns: ['**/*.{js,css,html,png,svg,ico,json}'],
			}
		})
	],
	server: {
		host: "0.0.0.0",
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
