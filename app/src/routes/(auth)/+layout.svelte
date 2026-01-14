<script lang="ts">
	let { children } = $props();
</script>

<div class="min-h-screen bg-bg flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full">
		{@render children()}
	</div>
</div>
