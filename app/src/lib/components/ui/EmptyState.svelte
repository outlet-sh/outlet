<!--
  Reusable Empty State Component
-->

<script lang="ts">
	import type { ComponentType, Snippet } from 'svelte';

	let {
		icon = 'inbox',
		title,
		message,
		description, // alias for message
		children
	}: {
		icon?: string | ComponentType;
		title: string;
		message?: string;
		description?: string;
		children?: Snippet;
	} = $props();

	const displayMessage = $derived(message ?? description);
</script>

<div class="p-12 text-center">
	<div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-bg-secondary">
		{#if typeof icon === 'string'}
			<i class="fas fa-{icon} text-2xl text-text-muted"></i>
		{:else}
			<svelte:component this={icon} class="h-8 w-8 text-text-muted" />
		{/if}
	</div>
	<p class="text-sm font-medium text-text">{title}</p>
	{#if displayMessage}
		<p class="mt-1 text-xs text-text-secondary">{displayMessage}</p>
	{/if}
	{#if children}
		<div class="mt-4">
			{@render children()}
		</div>
	{/if}
</div>

<style>
@reference "$src/app.css";

@layer components.empty-state {
	/* No custom styles needed - using Tailwind utilities */
}
</style>
