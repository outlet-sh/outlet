<script lang="ts">
	let { children } = $props();
</script>

<div class="min-h-screen bg-bg flex flex-col items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
	<div class="w-full max-w-5xl">
		{@render children()}
	</div>
</div>
