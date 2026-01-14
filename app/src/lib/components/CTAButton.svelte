<script lang="ts">
	import { trackEvent } from '$lib/analytics';

	function handleClick() {
		trackEvent('cta_click', {
			location: 'hero',
			button_text: 'Apply for Your AI Production Diagnostic'
		});
	}
</script>

<a href="/apply" onclick={handleClick} class="btn-primary group text-lg px-10 py-5 glow-emerald">
	<!-- Shimmer Effect -->
	<div class="absolute inset-0 -translate-x-full group-hover:translate-x-full transition-transform duration-1000 ease-in-out bg-gradient-to-r from-transparent via-white/20 to-transparent"></div>

	<span class="relative">Apply for Your AI Production Diagnostic</span>

	<svg class="relative w-5 h-5 group-hover:translate-x-1 transition-transform duration-300" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
		<path d="M7.5 5L12.5 10L7.5 15" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
	</svg>
</a>
