<!--
  ScrollArea Component
  Custom scrollable area with styled scrollbars
-->

<script lang="ts">
	interface Props {
		children: any;
		class?: string;
		orientation?: 'vertical' | 'horizontal';
	}

	let {
		children,
		class: className = '',
		orientation = 'vertical'
	}: Props = $props();

	const orientationClasses = {
		vertical: 'overflow-y-auto',
		horizontal: 'overflow-x-auto'
	};

	const scrollbarClasses = {
		vertical: 'scrollbar-thin scrollbar-track-transparent',
		horizontal: 'scrollbar-thin scrollbar-track-transparent'
	};
</script>

<div class="relative overflow-hidden {className}">
	<div class="h-full w-full rounded-[inherit] {orientationClasses[orientation]} {scrollbarClasses[orientation]}">
		{@render children()}
	</div>
</div>

<style>
	@reference "$src/app.css";

@layer components.scroll-area {
	/* Custom scrollbar styles for browsers that don't support scrollbar-* utilities */
	:global(.scrollbar-thin) {
		scrollbar-width: thin;
	}

	:global(.scrollbar-thin::-webkit-scrollbar) {
		width: 10px;
		height: 10px;
	}

	:global(.scrollbar-thin::-webkit-scrollbar-track) {
		background: transparent;
	}

	:global(.scrollbar-thin::-webkit-scrollbar-thumb) {
		background-color: var(--color-border);
		border-radius: 9999px;
		border: 2px solid transparent;
		background-clip: padding-box;
	}

	:global(.scrollbar-thin::-webkit-scrollbar-thumb:hover) {
		background-color: var(--color-text-muted);
	}
}
</style>
