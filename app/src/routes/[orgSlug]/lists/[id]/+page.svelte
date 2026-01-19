<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	// Redirect to subscribers tab
	$effect(() => {
		const basePath = `/${$page.params.orgSlug}/lists/${$page.params.id}`;
		goto(`${basePath}/subscribers`, { replaceState: true });
	});
</script>

<div class="flex justify-center py-12">
	<p class="text-text-muted">Redirecting...</p>
</div>
