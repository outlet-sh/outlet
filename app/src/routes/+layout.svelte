<script lang="ts">
	import { onMount } from 'svelte';
	import '../app.css';
	import '../styles/toast.css';

	const { children }: { children: any } = $props();
</script>

{@render children()}
