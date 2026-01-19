import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],

	resolve: {
		alias: {
			$src: 'src'
		}
	},

	server: {
		host: '0.0.0.0',
		proxy: {
			// Proxy API requests to Go backend during development
			'/api': {
				target: 'http://localhost:20202',
				changeOrigin: true
			},
			// Proxy health check
			'/health': {
				target: 'http://localhost:20202',
				changeOrigin: true
			},
			// Proxy subscription forms
			'/s/': {
				target: 'http://localhost:20202',
				changeOrigin: true
			},
			// Proxy subscription forms
			'/confirm/': {
				target: 'http://localhost:20202',
				changeOrigin: true
			},
			// Proxy subscription forms
			'/u/': {
				target: 'http://localhost:20202',
				changeOrigin: true
			},

			// Proxy WebSocket connections
			'/ws': {
				target: 'ws://localhost:20202',
				ws: true,
				changeOrigin: true
			}
		}
	}
});
