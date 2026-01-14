<!--
  Reusable Card Component
  Ensures 100% consistency for all card-based layouts
-->

<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		title,
		subtitle,
		children,
		action,
		class: extraClass = '',
		hover = false,
		onclick
	}: {
		title?: string;
		subtitle?: string;
		children: Snippet;
		action?: Snippet;
		class?: string;
		hover?: boolean;
		onclick?: () => void;
	} = $props();

	const isClickable = !!onclick;
	const className = `card ${hover ? 'group hover:border-indigo-500/50 transition-all' : ''} ${isClickable ? 'cursor-pointer' : ''} ${extraClass}`.trim();
</script>

<div class={className} onclick={onclick} onkeydown={(e) => e.key === 'Enter' && onclick?.()} role={isClickable ? 'button' : undefined} tabindex={isClickable ? 0 : undefined}>
	{#if title}
		<div class="flex items-center justify-between mb-4">
			<div>
				<h2 class="text-lg font-semibold text-text">{title}</h2>
				{#if subtitle}
					<p class="mt-1 text-sm text-text-muted">{subtitle}</p>
				{/if}
			</div>
			{#if action}
				{@render action()}
			{/if}
		</div>
	{/if}

	{@render children()}
</div>

<style>
	@reference "$src/app.css";
	@layer components.card {
		/* Card styles are defined in app.css */
	}
</style>
