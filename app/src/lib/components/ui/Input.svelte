<!--
  Reusable Input Component
-->

<script lang="ts">
	let {
		type = 'text',
		size = 'md',
		placeholder = '',
		value = $bindable(),
		element = $bindable(),
		disabled = false,
		id,
		class: extraClass = '',
		step,
		min,
		max,
		label,
		description,
		name,
		required = false,
		readonly = false,
		autocomplete = 'off',
		onfocus,
		oninput,
		onblur,
		onkeydown
	}: {
		type?: string;
		size?: 'sm' | 'md' | 'lg';
		placeholder?: string;
		value?: string | number;
		element?: HTMLInputElement;
		disabled?: boolean;
		id?: string;
		class?: string;
		step?: string | number;
		min?: string | number;
		max?: string | number;
		label?: string;
		description?: string;
		name?: string;
		required?: boolean;
		readonly?: boolean;
		autocomplete?: AutoFill;
		onfocus?: (e: Event) => void;
		oninput?: (e: Event) => void;
		onblur?: (e: Event) => void;
		onkeydown?: (e: KeyboardEvent) => void;
	} = $props();

	// Ensure value is never undefined
	let safeValue = $derived(value ?? '');

	const sizeClasses = {
		sm: 'input-sm',
		md: 'input-md',
		lg: 'input-lg'
	};

	const className = `input ${sizeClasses[size]} ${extraClass}`.trim();

	function handleInput(e: Event) {
		const input = e.currentTarget as HTMLInputElement;
		// Preserve number type for number inputs
		if (type === 'number') {
			value = input.valueAsNumber;
		} else {
			value = input.value;
		}
		oninput?.(e);
	}
</script>

{#if label}
	<label for={id} class="block text-sm font-medium text-text mb-1.5">{label}</label>
{/if}
{#if description}
	<p class="text-xs text-text-muted mb-1.5">{description}</p>
{/if}
<input
	bind:this={element}
	class={className}
	{type}
	{placeholder}
	value={safeValue}
	oninput={handleInput}
	{onfocus}
	{onblur}
	{onkeydown}
	{disabled}
	{id}
	{name}
	{step}
	{min}
	{max}
	{required}
	{readonly}
	{autocomplete}
/>

<style>
	@reference "$src/app.css";
	@layer components.input {
		/* ===== Input ===== */
		.input {
			@apply w-full px-4 py-3 rounded-xl transition-all duration-200;
			@apply disabled:opacity-50 disabled:cursor-not-allowed;
			@apply bg-bg text-text;
			border: 1px solid var(--color-border);
			font-size: 1rem; /* 16px minimum prevents iOS zoom on focus */
		}

		.input::placeholder {
			@apply text-text-muted;
		}

		.input:focus {
			@apply outline-none;
			border-color: var(--color-primary);
			box-shadow: 0 0 0 3px color-mix(in srgb, var(--color-primary) 20%, transparent);
		}

		/* Size variants */
		.input-sm {
			@apply h-8 px-3 py-1 text-sm;
		}

		.input-md {
			@apply h-11;
		}

		.input-lg {
			@apply h-12 px-5 py-4 text-base;
		}
	}
</style>
