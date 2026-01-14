<script lang="ts">
	import type { Snippet } from 'svelte';
	import { slide } from 'svelte/transition';
	import { X } from 'lucide-svelte';
	import { Button } from '$lib/components/ui';

	interface Props {
		show: boolean;
		title?: string;
		children: Snippet;
		footer?: Snippet;
		onclose?: () => void;
	}

	let { show = $bindable(), title, children, footer, onclose }: Props = $props();

	function handleClose() {
		show = false;
		onclose?.();
	}
</script>

{#if show}
	<div class="slide-form" transition:slide={{ duration: 200 }}>
		<div class="slide-form-header">
			{#if title}
				<h3 class="slide-form-title">{title}</h3>
			{/if}
			<button onclick={handleClose} class="slide-form-close" aria-label="Close">
				<X class="h-4 w-4" />
			</button>
		</div>
		<div class="slide-form-content">
			{@render children()}
		</div>
		{#if footer}
			<div class="slide-form-footer">
				{@render footer()}
			</div>
		{/if}
	</div>
{/if}

<style>
	@reference "$src/app.css";
	@layer components.slide-form {
		.slide-form {
			@apply rounded-lg border mb-4 overflow-hidden;
			background-color: var(--color-bg);
			border-color: var(--color-border);
		}

		.slide-form-header {
			@apply flex items-center justify-between px-4 py-3 border-b;
			background-color: var(--color-bg-secondary);
			border-color: var(--color-border);
		}

		.slide-form-title {
			@apply text-sm font-semibold;
			color: var(--color-text);
		}

		.slide-form-close {
			@apply p-1 rounded transition-colors;
			color: var(--color-text-muted);
		}

		.slide-form-close:hover {
			background-color: var(--color-bg-tertiary);
			color: var(--color-text);
		}

		.slide-form-content {
			@apply p-4;
		}

		.slide-form-footer {
			@apply px-4 py-3 border-t flex justify-end gap-3;
			background-color: var(--color-bg-secondary);
			border-color: var(--color-border);
		}
	}
</style>
