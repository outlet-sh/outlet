<script>
	import CTAButton from './CTAButton.svelte';
</script>

<section class="relative flex items-center overflow-hidden bg-gradient-dark-hero hero-padding">
	<!-- Diagonal Light Beam -->
	<div class="absolute top-0 right-0 w-1/2 h-full opacity-10 pointer-events-none">
		<div class="absolute top-0 right-0 w-full h-full diagonal-beam"></div>
	</div>

	<!-- Subtle Grid -->
	<div class="absolute inset-0 grid-background opacity-20"></div>

	<div class="container-wide relative z-10">
		<div class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16 items-center">
			<!-- LEFT: Headline + CTA -->
			<div class="lg:col-span-7">
				<!-- Badge -->
				<div class="mb-6">
					<span class="inline-flex items-center gap-2 px-4 py-2 rounded-full border font-mono text-xs font-semibold uppercase tracking-wider badge-blue text-blue-light">
						<svg width="10" height="10" viewBox="0 0 10 10" fill="none" xmlns="http://www.w3.org/2000/svg">
							<circle cx="5" cy="5" r="5" fill="currentColor"/>
						</svg>
						For Enterprise Leaders
					</span>
				</div>

				<h1 class="text-4xl md:text-4xl lg:text-5xl mb-6 text-balance text-slate-50 leading-tight font-bold">
					Is it unreasonable to want <span class="gradient-text">a clear path</span> before committing <span class="gradient-text">to your next AI initiative</span>?
				</h1>

				<!-- NFY Recognition (before stats) -->
				<p class="text-xl mb-6 max-w-2xl text-slate-300 leading-relaxed">
					You've seen pilots that work in demos but stall before production. You've watched promising initiatives quietly fade. The gap between what's possible and what's actually happening keeps widening.
				</p>

				<p class="text-xl md:text-2xl mb-10 max-w-2xl text-slate-300 leading-normal">
					Most organizations reach this point. Here's how the successful 12% moved forward.
				</p>

				<div class="mb-4">
					<CTAButton />
				</div>

				<!-- Micro-clarifier for cold traffic -->
				<p class="text-sm text-slate-500 text-center max-w-md mx-auto">
					No prep. No obligation. Just clarity.
				</p>
			</div>

			<!-- RIGHT: Stats Card (High-Converting Data Visualization) -->
			<div class="lg:col-span-5">
				<div class="relative">
					<!-- Glow Effect -->
					<div class="absolute inset-0 rounded-2xl blur-2xl opacity-20 glow-blue-small"></div>

					<!-- Card -->
					<div class="relative rounded-2xl p-8 border card-glass">
						<p class="text-sm font-mono font-semibold uppercase tracking-wider mb-2 text-slate-400">
							The Reality of AI Initiatives
						</p>

						<!-- Interpretive frame for cold traffic -->
						<p class="text-sm text-slate-300 mb-6 italic">
							The numbers aren't the problem. It's what they reveal.
						</p>

						<!-- Stats Grid -->
						<div class="grid grid-cols-2 gap-6 mb-6">
							<div>
								<div class="stat-number mb-2 stat-number-large">88%</div>
								<p class="text-sm font-mono uppercase tracking-wide leading-tight text-slate-400">Never Reach Production</p>
							</div>
							<div>
								<div class="stat-number mb-2 stat-number-large">95%</div>
								<p class="text-sm font-mono uppercase tracking-wide leading-tight text-slate-400">Show Zero ROI</p>
							</div>
							<div>
								<div class="stat-number mb-2 stat-number-large">74%</div>
								<p class="text-sm font-mono uppercase tracking-wide leading-tight text-slate-400">CEOs at Career Risk</p>
							</div>
							<div>
								<div class="stat-number mb-2 stat-number-large">54%</div>
								<p class="text-sm font-mono uppercase tracking-wide leading-tight text-slate-400">Feel They've Failed</p>
							</div>
						</div>

						<!-- Divider -->
						<div class="h-px mb-6 divider-gradient"></div>

						<!-- Quote (pattern interrupt spacing) -->
						<div class="relative pl-4 border-l-2 border-blue">
							<p class="text-base italic leading-relaxed text-slate-300 mb-1">
								"Feels like AI's everywhere except where it actually matters.
							</p>
							<p class="text-base italic leading-relaxed text-slate-300 mb-4">
								Everything works as a demo. Nothing works in production."
							</p>
							<p class="text-xs font-mono font-semibold uppercase tracking-wider text-blue-light">
								— Enterprise CTO, Nov 2025
							</p>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>
