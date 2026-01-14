<script lang="ts">
	import type { ContentMeta } from '$lib/content/loader';

	interface Props {
		meta: ContentMeta;
	}

	let { meta }: Props = $props();
</script>

<article class="article-card">
	<a href="/resources/articles/{meta.slug}" class="article-card-link">
		{#if meta.image}
			<div class="article-card-image">
				<img src={meta.image} alt={meta.title} />
			</div>
		{/if}

		<div class="article-card-content">
			{#if meta.featured}
				<span class="article-badge">Featured</span>
			{/if}

			<h3 class="article-card-title">{meta.title}</h3>
			<p class="article-card-excerpt">{meta.excerpt}</p>

			<div class="article-card-meta">
				<time datetime={meta.date}>{new Date(meta.date).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}</time>
				{#if meta.readTime}
					<span class="article-card-divider">•</span>
					<span>{meta.readTime}</span>
				{/if}
			</div>

			{#if meta.tags && meta.tags.length > 0}
				<div class="article-card-tags">
					{#each meta.tags.slice(0, 3) as tag}
						<span class="article-tag">{tag}</span>
					{/each}
				</div>
			{/if}
		</div>
	</a>
</article>

<style>
	@import 'tailwindcss/theme' reference;
	.article-card {
		@apply bg-white rounded-2xl border border-gray-200 overflow-hidden hover:shadow-xl hover:border-blue-300 transition-all duration-300;
	}

	.article-card-link {
		@apply block text-inherit no-underline;
	}

	.article-card-image {
		@apply w-full h-48 overflow-hidden bg-gradient-to-br from-blue-50 to-indigo-50;
	}

	.article-card-image img {
		@apply w-full h-full object-cover transition-transform duration-300;
	}

	.article-card:hover .article-card-image img {
		@apply scale-105;
	}

	.article-card-content {
		@apply p-6;
	}

	.article-badge {
		@apply inline-block px-3 py-1 text-xs font-semibold text-blue-600 bg-blue-50 rounded-full mb-3;
	}

	.article-card-title {
		@apply text-xl font-bold text-gray-900 mb-3 leading-tight;
	}

	.article-card-excerpt {
		@apply text-gray-600 mb-4 line-clamp-2;
	}

	.article-card-meta {
		@apply flex items-center text-sm text-gray-500 mb-3;
	}

	.article-card-divider {
		@apply mx-2;
	}

	.article-card-tags {
		@apply flex flex-wrap gap-2;
	}

	.article-tag {
		@apply text-xs px-2 py-1 bg-gray-100 text-gray-600 rounded;
	}
</style>
