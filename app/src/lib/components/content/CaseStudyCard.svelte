<script lang="ts">
	import type { ContentMeta } from '$lib/content/loader';
	import { TrendingUp, Clock, DollarSign } from 'lucide-svelte';

	interface Props {
		meta: ContentMeta & {
			industry?: string;
			companySize?: string;
			timeline?: string;
			roi?: string;
		};
	}

	let { meta }: Props = $props();
</script>

<article class="case-study-card">
	<a href="/case-studies/{meta.slug}" class="case-study-card-link">
		{#if meta.image}
			<div class="case-study-card-image">
				<img src={meta.image} alt={meta.title} />
				{#if meta.industry}
					<div class="case-study-industry-badge">
						{meta.industry}
					</div>
				{/if}
			</div>
		{/if}

		<div class="case-study-card-content">
			{#if meta.featured}
				<span class="case-study-badge">Featured Case Study</span>
			{/if}

			<h3 class="case-study-card-title">{meta.title}</h3>
			<p class="case-study-card-excerpt">{meta.excerpt}</p>

			<div class="case-study-stats">
				{#if meta.roi}
					<div class="case-study-stat">
						<TrendingUp class="w-4 h-4 text-green-600" />
						<div>
							<div class="stat-label">ROI</div>
							<div class="stat-value">{meta.roi}</div>
						</div>
					</div>
				{/if}

				{#if meta.timeline}
					<div class="case-study-stat">
						<Clock class="w-4 h-4 text-blue-600" />
						<div>
							<div class="stat-label">Timeline</div>
							<div class="stat-value">{meta.timeline}</div>
						</div>
					</div>
				{/if}

				{#if meta.companySize}
					<div class="case-study-stat">
						<div class="stat-icon">👥</div>
						<div>
							<div class="stat-label">Company Size</div>
							<div class="stat-value">{meta.companySize}</div>
						</div>
					</div>
				{/if}
			</div>

			{#if meta.tags && meta.tags.length > 0}
				<div class="case-study-card-tags">
					{#each meta.tags.slice(0, 3) as tag}
						<span class="case-study-tag">{tag}</span>
					{/each}
				</div>
			{/if}

			<div class="case-study-card-cta">
				<span class="case-study-cta-text">Read Case Study</span>
				<span class="case-study-arrow">→</span>
			</div>
		</div>
	</a>
</article>

<style>
	@import 'tailwindcss/theme' reference;
	.case-study-card {
		@apply bg-white rounded-2xl border border-gray-200 overflow-hidden hover:shadow-xl hover:border-blue-300 transition-all duration-300;
	}

	.case-study-card-link {
		@apply block text-inherit no-underline;
	}

	.case-study-card-image {
		@apply relative w-full h-56 overflow-hidden bg-gradient-to-br from-blue-50 to-indigo-50;
	}

	.case-study-card-image img {
		@apply w-full h-full object-cover transition-transform duration-300;
	}

	.case-study-card:hover .case-study-card-image img {
		@apply scale-105;
	}

	.case-study-industry-badge {
		@apply absolute top-4 left-4 bg-white/95 backdrop-blur-sm text-gray-900 px-3 py-1 rounded-full text-sm font-semibold shadow-md;
	}

	.case-study-card-content {
		@apply p-6;
	}

	.case-study-badge {
		@apply inline-block px-3 py-1 text-xs font-semibold text-green-600 bg-green-50 rounded-full mb-3;
	}

	.case-study-card-title {
		@apply text-xl font-bold text-gray-900 mb-3 leading-tight;
	}

	.case-study-card-excerpt {
		@apply text-gray-600 mb-4 line-clamp-2;
	}

	.case-study-stats {
		@apply grid grid-cols-2 md:grid-cols-3 gap-4 mb-4 pb-4 border-b border-gray-100;
	}

	.case-study-stat {
		@apply flex items-start gap-2;
	}

	.stat-icon {
		@apply text-lg;
	}

	.stat-label {
		@apply text-xs text-gray-500 uppercase tracking-wide;
	}

	.stat-value {
		@apply text-sm font-semibold text-gray-900;
	}

	.case-study-card-tags {
		@apply flex flex-wrap gap-2 mb-4;
	}

	.case-study-tag {
		@apply text-xs px-2 py-1 bg-gray-100 text-gray-600 rounded;
	}

	.case-study-card-cta {
		@apply flex items-center justify-between text-blue-600 font-semibold;
	}

	.case-study-card:hover .case-study-card-cta {
		@apply text-blue-700;
	}

	.case-study-cta-text {
		@apply transition-transform duration-300;
	}

	.case-study-card:hover .case-study-cta-text {
		@apply translate-x-1;
	}

	.case-study-arrow {
		@apply transition-transform duration-300;
	}

	.case-study-card:hover .case-study-arrow {
		@apply translate-x-2;
	}
</style>
