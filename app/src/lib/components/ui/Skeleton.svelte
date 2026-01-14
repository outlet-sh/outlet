<!--
  Skeleton Component
  Loading placeholder with pulse animation
-->

<script lang="ts">
	interface Props {
		class?: string;
	}

	let {
		class: className = ''
	}: Props = $props();
</script>

<div class="animate-pulse rounded-md bg-bg-secondary {className}"></div>

<style>
@reference "$src/app.css";

@layer components.skeleton {
	/* No custom styles needed - using Tailwind utilities */
}
</style>
