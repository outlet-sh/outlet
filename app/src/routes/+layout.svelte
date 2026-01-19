<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { getPublicSetupStatus } from '$lib/api';
	import '$src/app.css';

	let { children } = $props();

	onMount(async () => {
		// Skip setup check if already on setup page
		if ($page.url.pathname.startsWith('/setup')) {
			return;
		}

		try {
			const status = await getPublicSetupStatus();
			if (status.setup_required) {
				goto('/setup');
			}
		} catch (err) {
			// If API fails, might be first run - try setup page
			console.error('Failed to check setup status:', err);
		}
	});
</script>

<svelte:head>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</svelte:head>

{@render children()}
