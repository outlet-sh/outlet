import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { mdsvex } from 'mdsvex';
import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';
import rehypePrettyCode from 'rehype-pretty-code';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const config = {
	extensions: ['.svelte', '.md'],
	preprocess: [
		vitePreprocess(),
		mdsvex({
			extensions: ['.md'],
			layout: {
				_: resolve(__dirname, './src/lib/layouts/markdown.svelte')
			},
			rehypePlugins: [
				[
					rehypePrettyCode,
					{
						theme: 'github-dark',
						keepBackground: true
					}
				]
			]
		})
	],
	kit: {
		alias: {
			$content: resolve(__dirname, './src/content')
		},
		adapter: adapter({
			pages: 'build',
			assets: 'build',
			fallback: '200.html', // SPA fallback for non-prerendered routes
			precompress: false,
			strict: false
		}),
		prerender: {
			entries: ['*'], // Prerender all discoverable routes
			crawl: true, // Automatically crawl links to find routes
			handleHttpError: 'ignore',
			handleMissingId: 'ignore',
			handleUnseenRoutes: 'ignore' // Ignore dynamic routes without data
		}
	}
};

export default config;
