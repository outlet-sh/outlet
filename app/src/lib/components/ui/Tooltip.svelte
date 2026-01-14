<!--
  Tooltip Component
  Hover or click tooltip with positioning
-->

<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	interface Props {
		content: string;
		position?: 'top' | 'bottom' | 'left' | 'right';
		trigger?: 'hover' | 'click';
		children: any;
	}

	let {
		content,
		position = 'bottom',
		trigger = 'click',
		children
	}: Props = $props();

	let showTooltip = $state(false);
	let containerEl: HTMLSpanElement;

	function handleClick(e: MouseEvent) {
		if (trigger === 'click') {
			showTooltip = !showTooltip;
		}
	}

	function handleClickOutside(e: MouseEvent) {
		if (trigger === 'click' && showTooltip) {
			// Close if click is outside this tooltip's container
			if (containerEl && !containerEl.contains(e.target as Node)) {
				showTooltip = false;
			}
		}
	}

	function handleMouseEnter() {
		if (trigger === 'hover') {
			showTooltip = true;
		}
	}

	function handleMouseLeave() {
		if (trigger === 'hover') {
			showTooltip = false;
		}
	}

	onMount(() => {
		// Use capture phase to catch clicks before they reach other handlers
		document.addEventListener('click', handleClickOutside, true);
	});

	onDestroy(() => {
		document.removeEventListener('click', handleClickOutside, true);
	});
</script>

<span
	class="tooltip-container"
	bind:this={containerEl}
	onmouseenter={handleMouseEnter}
	onmouseleave={handleMouseLeave}
	onclick={handleClick}
	role="button"
	tabindex="0"
>
	{@render children()}

	{#if showTooltip}
		<div
			class="tooltip-wrapper tooltip-{position}"
			role="tooltip"
		>
			<div class="tooltip-content">
				{content}
			</div>
			<div class="tooltip-arrow tooltip-arrow-{position}"></div>
		</div>
	{/if}
</span>

<style>
	@reference "$src/app.css";
	@layer components.tooltip {
		.tooltip-container {
			@apply relative cursor-pointer;
			display: inline;
		}

		.tooltip-wrapper {
			@apply absolute z-50;
			animation: fade-in 200ms ease-out;
		}

		.tooltip-top {
			@apply bottom-full left-1/2 -translate-x-1/2 mb-2;
		}

		.tooltip-bottom {
			@apply top-full left-1/2 -translate-x-1/2 mt-2;
		}

		.tooltip-left {
			@apply right-full top-1/2 -translate-y-1/2 mr-2;
		}

		.tooltip-right {
			@apply left-full top-1/2 -translate-y-1/2 ml-2;
		}

		.tooltip-content {
			@apply rounded-lg px-3 py-2 text-sm font-medium shadow-xl;
			background: var(--color-base-800);
			color: var(--color-base-50);
			max-width: 320px;
			min-width: 200px;
		}

		.tooltip-arrow {
			@apply absolute border-4;
			border-color: var(--color-base-800);
		}

		.tooltip-arrow-top {
			@apply top-full left-1/2 -translate-x-1/2;
			border-left-color: transparent;
			border-right-color: transparent;
			border-bottom-color: transparent;
		}

		.tooltip-arrow-bottom {
			@apply bottom-full left-1/2 -translate-x-1/2;
			border-left-color: transparent;
			border-right-color: transparent;
			border-top-color: transparent;
		}

		.tooltip-arrow-left {
			@apply left-full top-1/2 -translate-y-1/2;
			border-top-color: transparent;
			border-bottom-color: transparent;
			border-right-color: transparent;
		}

		.tooltip-arrow-right {
			@apply right-full top-1/2 -translate-y-1/2;
			border-top-color: transparent;
			border-bottom-color: transparent;
			border-left-color: transparent;
		}
	}
</style>
