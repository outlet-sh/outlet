<script lang="ts">
	import PageHeader from '$lib/components/PageHeader.svelte';
	import { onMount } from 'svelte';

	let {
		header,
		subheader,
		children
	}: {
		header: string;
		subheader: string;
		children?: any;
	} = $props();

	let mounted = $state(false);

	onMount(() => {
		mounted = true;
	});
</script>

<PageHeader {header} {subheader} />
<div class="prose-lg space-y-12 md:space-y-16">
	{@render children()}
</div>
