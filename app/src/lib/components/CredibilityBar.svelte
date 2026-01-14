<script>
	// Video placeholder hidden for now - uncomment when video is ready
	// import HeroVideoPlaceholder from './HeroVideoPlaceholder.svelte';
</script>

<!--
  CREDIBILITY BAR
  Purpose: Show proof of competence visually
  Position: Right after emotional intro, before heavy text begins
-->
<section class="section-compact bg-slate-50">
	<div class="container-wide">
		<!-- Trust Statement -->
		<div class="text-center mb-12">
			<p class="text-sm font-mono uppercase tracking-wider mb-4 text-slate-500">Engaged and Consulted With Teams At</p>
			<p class="text-xl md:text-2xl max-w-4xl mx-auto font-semibold text-slate-900 leading-relaxed">
				Google, Walmart, Nike, HotelPlanner, North American Bancard, Northwestern Mutual, HarvestRight, SRT Group, and others.
			</p>
		</div>

		<!-- Logo Grid - Optimized for Trust -->
		<div class="mb-16">
			<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4 max-w-5xl mx-auto">
				<!-- Replace these divs with actual <img> tags when logos are available -->
				<div class="group relative rounded-xl border flex items-center justify-center h-24 transition-all duration-300 hover:scale-105 bg-white border-slate-200 shadow-card">
					<span class="text-sm font-bold font-mono transition-colors group-hover:text-blue-600 text-slate-400">Google</span>
				</div>
				<div class="group relative rounded-xl border flex items-center justify-center h-24 transition-all duration-300 hover:scale-105 bg-white border-slate-200 shadow-card">
					<span class="text-sm font-bold font-mono transition-colors group-hover:text-blue-600 text-slate-400">Walmart</span>
				</div>
				<div class="group relative rounded-xl border flex items-center justify-center h-24 transition-all duration-300 hover:scale-105 bg-white border-slate-200 shadow-card">
					<span class="text-sm font-bold font-mono transition-colors group-hover:text-blue-600 text-slate-400">Nike</span>
				</div>
				<div class="group relative rounded-xl border flex items-center justify-center h-24 transition-all duration-300 hover:scale-105 bg-white border-slate-200 shadow-card">
					<span class="text-sm font-bold font-mono transition-colors group-hover:text-blue-600 text-slate-400">Northwestern Mutual</span>
				</div>
				<div class="group relative rounded-xl border flex items-center justify-center h-24 transition-all duration-300 hover:scale-105 bg-white border-slate-200 shadow-card">
					<span class="text-sm font-bold font-mono transition-colors group-hover:text-blue-600 text-slate-400">HotelPlanner</span>
				</div>
			</div>
		</div>

		<!-- Hero VSL Video - Hidden until video is ready -->
		<!-- <HeroVideoPlaceholder /> -->
	</div>
</section>
