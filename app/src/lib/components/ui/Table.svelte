<!--
  Reusable Table Component
-->

<script lang="ts">
	let { children }: { children: any } = $props();
</script>

<div class="table-container">
	<table class="table">
		{@render children()}
	</table>
</div>

<style>
	@reference "$src/app.css";
	@layer components.table {
		.table-container {
			@apply overflow-x-auto rounded-lg bg-bg;
			border: 1px solid var(--color-border);
		}

		.table {
			@apply w-full text-sm;
		}

		.table :global(thead) {
			@apply border-b;
			border-color: var(--color-border);
		}

		.table :global(th) {
			@apply px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-text-muted;
		}

		.table :global(td) {
			@apply px-4 py-3 text-text;
		}

		.table :global(tr) {
			@apply border-b;
			border-color: color-mix(in srgb, var(--color-border) 50%, transparent);
		}

		.table :global(tbody tr:hover) {
			@apply bg-bg-secondary;
		}
	}
</style>
