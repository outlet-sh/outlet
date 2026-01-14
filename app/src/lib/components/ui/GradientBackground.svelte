<!--
  GradientBackground Component
  Consistent gradient background with optional noise texture
-->

<script lang="ts">
	interface Props {
		children: any;
		class?: string;
		withNoise?: boolean;
		noiseFrequency?: number;
		noiseOctaves?: number;
		noiseOpacity?: number;
	}

	let {
		children,
		class: className = '',
		withNoise = true,
		noiseFrequency = 1.2,
		noiseOctaves = 4,
		noiseOpacity = 0.03
	}: Props = $props();

	let key = $state(Date.now());

	// Force re-render when noise parameters change
	$effect(() => {
		key = Date.now();
	});

	const noiseUrl = $derived(
		`data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='${noiseFrequency}' numOctaves='${noiseOctaves}' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E`
	);
</script>

<div class="bg-gradient-to-b from-bg-secondary/40 via-bg-secondary/40 to-bg {className} relative">
	{#if withNoise}
		{#key key}
			<div
				class="absolute top-0 right-0 bottom-0 left-0 pointer-events-none z-0"
				style="opacity: {noiseOpacity}; background-image: url('{noiseUrl}');"
				aria-hidden="true"
			></div>
		{/key}
	{/if}
	<div class="relative z-10">
		{@render children()}
	</div>
</div>

<style>
@reference "$src/app.css";

@layer components.gradient-background {
	/* No custom styles needed - using Tailwind utilities */
}
</style>
