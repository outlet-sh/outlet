<script lang="ts">
	import { trackEvent } from '$lib/analytics';

	function handleClick() {
		trackEvent('cta_click', {
			location: 'bottom_cta',
			button_text: 'Apply for Your AI Pilot-to-Production Diagnostic'
		});
	}
</script>

<section class="section diagonal-accent bg-slate-950 text-slate-100">
	<div class="absolute inset-0 opacity-30 pointer-events-none">
		<div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] rounded-full blur-3xl bg-radial-blue-glow"></div>
	</div>

	<div class="container-content relative z-10">
		<div class="max-w-4xl mx-auto text-center">
			<h2 class="mb-6 leading-tight text-slate-50">
				Apply for Your AI Pilot-to-Production Diagnostic
			</h2>
			<p class="text-xl mb-4 text-slate-300 leading-relaxed">
				You don't need more input. You need to isolate the blocker. That's what this does.
			</p>

			<!-- Strategic framing -->
			<p class="text-lg mb-10 text-slate-400 leading-relaxed italic max-w-2xl mx-auto">
				Most leaders wait for clarity. The best extract it.
			</p>

			<a href="/apply" onclick={handleClick} class="btn-primary text-lg px-12 py-5 mb-4 bg-white text-slate-950">
				Apply for Your AI Pilot-to-Production Diagnostic
				<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
					<path d="M7.5 5L12.5 10L7.5 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
				</svg>
			</a>

			<p class="text-base text-slate-300 leading-relaxed max-w-2xl mx-auto">
				A 25-minute structured diagnostic that isolates your blocker — and reveals whether working together makes sense. You keep the clarity either way.
			</p>
		</div>
	</div>
</section>
