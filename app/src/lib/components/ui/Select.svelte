<!--
  Reusable Select Component
  Supports both children (slot-based) and options prop patterns
-->

<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Option {
		value: string | number;
		label: string;
	}

	let {
		value = $bindable(),
		size = 'md',
		disabled = false,
		options,
		children,
		onchange,
		class: extraClass = '',
		...restProps
	}: {
		value?: string | number;
		size?: 'sm' | 'md' | 'lg';
		disabled?: boolean;
		options?: Option[];
		children?: Snippet;
		onchange?: () => void;
		class?: string;
		[key: string]: any;
	} = $props();

	let safeValue = $derived(value ?? '');

	const sizeClasses = {
		sm: 'select-sm',
		md: 'select-md',
		lg: 'select-lg'
	};

	const className = `select ${sizeClasses[size]} ${extraClass}`.trim();

	function handleChange(e: Event) {
		const target = e.currentTarget as HTMLSelectElement;
		// Try to preserve number type if original value was a number
		if (typeof value === 'number') {
			value = Number(target.value);
		} else {
			value = target.value;
		}
		onchange?.();
	}
</script>

<select class={className} value={safeValue} onchange={handleChange} {disabled} {...restProps}>
	{#if options}
		{#each options as opt}
			<option value={opt.value}>{opt.label}</option>
		{/each}
	{:else if children}
		{@render children()}
	{/if}
</select>

<style>
	@reference "$src/app.css";
	@layer components.select {
		.select {
			@apply w-full px-4 py-3 rounded-xl transition-all duration-200 bg-bg text-text;
			@apply disabled:opacity-50 disabled:cursor-not-allowed;
			border: 1px solid var(--color-border);
			font-size: 1rem; /* 16px minimum prevents iOS zoom on focus */
		}

		.select:focus {
			@apply outline-none;
			border-color: var(--color-primary);
			box-shadow: 0 0 0 3px color-mix(in srgb, var(--color-primary) 20%, transparent);
		}

		.select-sm {
			@apply h-8 px-3 py-1 text-sm;
		}

		.select-md {
			@apply h-11;
		}

		.select-lg {
			@apply h-12 px-5 py-4 text-base;
		}
	}
</style>
