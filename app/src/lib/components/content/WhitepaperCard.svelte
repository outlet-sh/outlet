<script lang="ts">
	import type { ContentMeta } from '$lib/content/loader';
	import { FileText, Download } from 'lucide-svelte';

	interface Props {
		meta: ContentMeta;
	}

	let { meta }: Props = $props();
</script>

<article class="whitepaper-card">
	<a href="/resources/whitepapers/{meta.slug}" class="whitepaper-card-link">
		{#if meta.image}
			<div class="whitepaper-card-image">
				<img src={meta.image} alt={meta.title} />
				{#if meta.gated}
					<div class="whitepaper-gated-badge">
						<Download class="w-4 h-4" />
					</div>
				{/if}
			</div>
		{:else}
			<div class="whitepaper-card-placeholder">
				<FileText class="w-12 h-12 text-blue-600" />
			</div>
		{/if}

		<div class="whitepaper-card-content">
			<h3 class="whitepaper-card-title">{meta.title}</h3>
			<p class="whitepaper-card-excerpt">{meta.excerpt}</p>

			<div class="whitepaper-card-meta">
				{#if meta.pageCount}
					<span class="whitepaper-meta-item">
						<FileText class="w-4 h-4" />
						{meta.pageCount} pages
					</span>
				{/if}
				{#if meta.gated}
					<span class="whitepaper-meta-item">
						<Download class="w-4 h-4" />
						Free Download
					</span>
				{/if}
			</div>

			{#if meta.tags && meta.tags.length > 0}
				<div class="whitepaper-card-tags">
					{#each meta.tags.slice(0, 3) as tag}
						<span class="whitepaper-tag">{tag}</span>
					{/each}
				</div>
			{/if}

			<div class="whitepaper-card-cta">
				<span class="whitepaper-cta-text">Download Now</span>
				<Download class="w-4 h-4" />
			</div>
		</div>
	</a>
</article>

<style>
	@import 'tailwindcss/theme' reference;
	.whitepaper-card {
		@apply bg-white rounded-2xl border border-gray-200 overflow-hidden hover:shadow-xl hover:border-blue-300 transition-all duration-300;
	}

	.whitepaper-card-link {
		@apply block text-inherit no-underline;
	}

	.whitepaper-card-image {
		@apply relative w-full h-48 overflow-hidden bg-gradient-to-br from-indigo-50 to-blue-50;
	}

	.whitepaper-card-image img {
		@apply w-full h-full object-cover transition-transform duration-300;
	}

	.whitepaper-card:hover .whitepaper-card-image img {
		@apply scale-105;
	}

	.whitepaper-gated-badge {
		@apply absolute top-4 right-4 bg-blue-600 text-white p-2 rounded-lg shadow-lg;
	}

	.whitepaper-card-placeholder {
		@apply w-full h-48 flex items-center justify-center bg-gradient-to-br from-indigo-50 to-blue-50;
	}

	.whitepaper-card-content {
		@apply p-6;
	}

	.whitepaper-card-title {
		@apply text-xl font-bold text-gray-900 mb-3 leading-tight;
	}

	.whitepaper-card-excerpt {
		@apply text-gray-600 mb-4 line-clamp-2;
	}

	.whitepaper-card-meta {
		@apply flex items-center gap-4 text-sm text-gray-500 mb-3;
	}

	.whitepaper-meta-item {
		@apply flex items-center gap-1;
	}

	.whitepaper-card-tags {
		@apply flex flex-wrap gap-2 mb-4;
	}

	.whitepaper-tag {
		@apply text-xs px-2 py-1 bg-gray-100 text-gray-600 rounded;
	}

	.whitepaper-card-cta {
		@apply flex items-center justify-between text-blue-600 font-semibold pt-4 border-t border-gray-100;
	}

	.whitepaper-card:hover .whitepaper-card-cta {
		@apply text-blue-700;
	}

	.whitepaper-cta-text {
		@apply transition-transform duration-300;
	}

	.whitepaper-card:hover .whitepaper-cta-text {
		@apply translate-x-1;
	}
</style>
