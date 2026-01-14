<script>
	import CTAButton from './CTAButton.svelte';
</script>

<!-- Example "About" section - customize for your landing pages -->
<section class="section bg-white">
	<div class="container-wide">

		<!-- Headline Section -->
		<div class="mb-12 max-w-4xl mx-auto text-center">
			<h2 class="mb-6 text-slate-950">
				Why <span class="gradient-text">This Platform</span>
			</h2>
			<p class="text-xl text-slate-600 leading-relaxed">
				Built by developers, for developers. Everything you need to launch and scale your SaaS.
			</p>
		</div>

		<!-- 2-UP GRID: TEXT LEFT, IMAGE RIGHT -->
		<div class="grid md:grid-cols-2 gap-16 items-start max-w-6xl">

			<!-- LEFT COLUMN -->
			<div class="space-y-6">

				<p class="text-lg leading-relaxed text-slate-700">
					We've been building SaaS products for years. We know what it takes.
				</p>

				<p class="text-lg leading-relaxed text-slate-700">
					This platform includes everything we wish we had when we started:
				</p>

				<ul class="space-y-2 text-lg list-none pl-0 mb-6">
					<li class="flex items-start gap-3">
						<span class="text-slate-400 mt-1">•</span>
						<span class="text-slate-700 leading-relaxed">Multi-tenant architecture from day one</span>
					</li>
					<li class="flex items-start gap-3">
						<span class="text-slate-400 mt-1">•</span>
						<span class="text-slate-700 leading-relaxed">Email marketing with sequences and campaigns</span>
					</li>
					<li class="flex items-start gap-3">
						<span class="text-slate-400 mt-1">•</span>
						<span class="text-slate-700 leading-relaxed">Billing with Stripe integration (BYOK)</span>
					</li>
					<li class="flex items-start gap-3">
						<span class="text-slate-400 mt-1">•</span>
						<span class="text-slate-700 leading-relaxed">AI/LLM gateway with cost tracking</span>
					</li>
					<li class="flex items-start gap-3">
						<span class="text-slate-400 mt-1">•</span>
						<span class="text-slate-700 leading-relaxed">Headless CMS for content management</span>
					</li>
				</ul>

				<p class="text-lg leading-relaxed text-slate-700">
					Stop rebuilding the same infrastructure for every project. Start with everything you need.
				</p>

				<p class="text-lg leading-relaxed font-semibold text-slate-900">
					Deploy in minutes. Scale with confidence.
				</p>

			</div>

			<!-- RIGHT COLUMN (STATS + FEATURES) -->
			<div class="flex flex-col items-end">
				<!-- Stats Card -->
				<div class="bg-gradient-to-br from-slate-50 to-white rounded-xl border border-gray-200 shadow-sm p-6 mb-4 w-full max-w-sm">
					<div class="grid grid-cols-2 gap-4 text-sm">
						<div class="text-center">
							<p class="text-2xl font-bold text-emerald-600">10X</p>
							<p class="text-xs text-gray-600 leading-tight">Faster<br />Development</p>
						</div>
						<div class="text-center">
							<p class="text-2xl font-bold text-emerald-600">100%</p>
							<p class="text-xs text-gray-600 leading-tight">Source Code<br />Ownership</p>
						</div>
						<div class="text-center col-span-2">
							<p class="text-xs font-semibold text-gray-700 mb-1">Built for</p>
							<p class="text-sm font-bold text-gray-900">Solo Founders & Small Teams</p>
						</div>
					</div>
				</div>

				<p class="text-base text-slate-700 font-medium text-right leading-relaxed max-w-sm">
					Self-hosted. No vendor lock-in. No recurring platform fees.
				</p>
			</div>

		</div>

	</div>
</section>
