<!--
  Textarea Component
  Multi-line text input
-->

<script lang="ts">
	interface Props {
		value?: string;
		placeholder?: string;
		disabled?: boolean;
		rows?: number;
		class?: string;
		id?: string;
		name?: string;
		required?: boolean;
		label?: string;
		description?: string;
	}

	let {
		value = $bindable(''),
		placeholder = '',
		disabled = false,
		rows = 3,
		class: className = '',
		id,
		name,
		required = false,
		label,
		description
	}: Props = $props();
</script>

{#if label || description}
	<div class="space-y-1.5">
		{#if label}
			<label for={id} class="text-sm font-medium text-text">{label}</label>
		{/if}
		<textarea
			bind:value
			{placeholder}
			{disabled}
			{rows}
			{id}
			{name}
			{required}
			class="textarea {className}"
		></textarea>
		{#if description}
			<p class="text-xs text-text-muted">{description}</p>
		{/if}
	</div>
{:else}
	<textarea
		bind:value
		{placeholder}
		{disabled}
		{rows}
		{id}
		{name}
		{required}
		class="textarea {className}"
	></textarea>
{/if}

<style>
	@reference "$src/app.css";
	@layer components.textarea {
		.textarea {
			@apply flex min-h-[80px] w-full rounded-md border px-3 py-2 text-sm;
			@apply bg-bg text-text;
			@apply focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2;
			@apply disabled:cursor-not-allowed disabled:opacity-50;
			border-color: var(--color-border);
		}

		.textarea::placeholder {
			@apply text-text-muted;
		}

		.textarea:focus-visible {
			@apply ring-primary ring-offset-bg;
		}
	}
</style>
