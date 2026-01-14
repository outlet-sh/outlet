<script lang="ts">
	import { Copy, Check } from 'lucide-svelte';
	import Button from './Button.svelte';
	import Badge from './Badge.svelte';
	import { cn } from '$lib/utils/cn';

	let {
		code,
		language,
		badgeText,
		badgeColor = 'bg-blue-600 hover:bg-blue-600',
		title,
		footer,
		disabled = false,
		class: className,
		style,
		height
	}: {
		code: string;
		language: string;
		badgeText?: string;
		badgeColor?: string;
		title?: string;
		footer?: any;
		disabled?: boolean;
		class?: string;
		style?: string;
		height?: string | number;
	} = $props();

	let copied = $state(false);

	const handleCopyCode = () => {
		navigator.clipboard.writeText(code);
		copied = true;
		setTimeout(() => (copied = false), 1500);
	};

	// Combine styles including height if provided
	const containerStyle = height
		? `${style || ''}; height: ${typeof height === 'number' ? `${height}px` : height}; display: flex; flex-direction: column;`
		: `${style || ''}; display: flex; flex-direction: column;`;
</script>

<div
	class={cn(
		'rounded-md overflow-hidden border border-border bg-bg text-text',
		className
	)}
	style={containerStyle}
>
	<div
		class={cn(
			'flex items-center px-4 py-2 justify-between border-b border-border'
		)}
	>
		<div class="flex items-center gap-2">
			{#if badgeText}
				<Badge class={`${badgeColor} text-white text-xs px-1.5 py-0 rounded h-5`}>
					{badgeText}
				</Badge>
			{/if}
			{#if title}
				<span class="text-xs font-medium text-text-secondary">{title}</span>
			{/if}
		</div>
		<Button
			type="secondary"
			{disabled}
			onclick={handleCopyCode}
		>
			<div class="h-6 w-6 p-0 flex items-center justify-center text-text-muted hover:text-text">
				{#if copied}
					<Check class="h-3 w-3" />
				{:else}
					<Copy class="h-3 w-3" />
				{/if}
			</div>
		</Button>
	</div>
	<div class="px-4 py-3 flex-1 overflow-auto bg-bg-secondary">
		<pre class="text-xs font-mono m-0 p-0 text-text"><code>{code}</code></pre>
	</div>
	{#if footer}
		<div
			class="px-4 py-2 border-t border-border text-text-secondary"
		>
			{@render footer()}
		</div>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.code-block {
		/* Code block uses utility classes */
	}
</style>
