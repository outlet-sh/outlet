<script lang="ts">
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';

	const mounted = writable(false);
	let {
		header,
		subheader,
		children
	}: {
		header: string;
		subheader: string;
		children?: any;
	} = $props();

	onMount(() => {
		mounted.set(true);
	});
</script>

<!-- Header -->
<div
	class="mb-16 opacity-0 transition-all duration-700 ease-out md:mb-24"
	class:opacity-100={mounted}
	class:translate-y-0={mounted}
	class:translate-y-5={!mounted}
>
	<h1 class="mb-6 text-6xl font-semibold tracking-tight text-black md:mb-8 md:text-7xl">{header}</h1>
	<p class="text-2xl leading-relaxed text-gray-600 md:text-3xl">
		{subheader}
	</p>
	{#if children}
		{children}
	{/if}
</div>
