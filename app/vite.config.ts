import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { fileURLToPath } from 'url';
import path from 'path';
/// <reference types="vitest" />

const __dirname = path.dirname(fileURLToPath(import.meta.url));

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],

	resolve: {
		alias: {
			$src: path.resolve(__dirname, 'src')
		}
	},

	optimizeDeps: {
		include: []
	},

	build: {
		rollupOptions: {}
	},

	define: {
		global: 'globalThis'
	},



	server: {
		fs: {
			allow: ['..']
		},
		proxy: {
			'/api': {
				target: 'http://localhost:9888',
				changeOrigin: true,
				ws: true
			}
		}
	},

	test: {
		environment: 'jsdom',
		globals: true,
		setupFiles: ['src/test/setup.ts']
	}
});
