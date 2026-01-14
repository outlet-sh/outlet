<script lang="ts">
	let openIndex = $state<number | null>(null);

	const faqs = [
		{
			question: "What exactly is the AI Pilot-to-Production Diagnostic?",
			answer: "It's a 25-minute conversation designed to identify why your AI initiatives aren't progressing beyond pilots and what the next realistic step is—without committing to a long engagement."
		},
		{
			question: "Do I need to prepare anything?",
			answer: "No. You don't need documents, dashboards, or reports. You talk; the diagnostic framework organizes."
		},
		{
			question: "Is this conversation confidential?",
			answer: "Yes. Completely.\n\nThis diagnostic exists to give you clarity — which only works if the conversation is candid, unfiltered, and protected.\n\nI do not share what we discuss — with your team, your vendors, your board, or anyone else.\n\nYou can speak openly about what's working, what's stuck, what's political, and what you're worried about — without needing to manage perception.\n\nNo recordings. No reporting. Just a trusted, confidential space to see what's actually blocking production."
		},
		{
			question: "Is this safe to discuss given internal politics?",
			answer: "Yes. We frame everything in a way that protects you and positions discoveries as environmental rather than personal gaps."
		},
		{
			question: "What if we've already hired a consulting firm?",
			answer: "This doesn't conflict. Most leaders book the diagnostic because consultants delivered strategy but nothing is shipping."
		},
		{
			question: "Why is this free?",
			answer: "Because the diagnostic is how I determine whether I can help. If there's a real production-readiness path, we talk about it. If not, the call ends — and you keep the insight.\n\nThis protects your time and mine."
		},
		{
			question: "Is this a sales call?",
			answer: "No — it's a diagnostic conversation designed to isolate the blocker that's preventing your AI work from moving into production.\n\nThe first 20–25 minutes are focused entirely on clarity: what's working, what's stuck, and which of the eight structural patterns applies to your environment.\n\nIf the diagnostic reveals that I can materially help you, we'll spend the last 3–5 minutes discussing what that could look like. If not, the call ends there. You keep the clarity either way.\n\nIt's not a sales call. But it is the step that determines whether working together makes sense."
		},
		{
			question: "Will you tell me if the issue is my team?",
			answer: "If the issue is structural, political, or architectural (it usually is), I'll say that. If your environment isn't ready for production, I'll say that too—candidly and privately."
		},
		{
			question: "Will this expose gaps in my leadership?",
			answer: "No. The diagnostic focuses on system blockers—not personal capability. In nearly every case, the blockers are environmental, not leadership failures. You'll leave with clarity about what to address, not shame about what you missed."
		},
		{
			question: "How do I know this isn't just another roadmap exercise?",
			answer: "There is no roadmap in this session. Just signal extraction, pattern recognition, and clarity."
		},
		{
			question: "What happens if we discover we're not ready for production?",
			answer: "You'll understand why, what's missing, and what order to address it in. For many leaders, this is the most valuable outcome."
		},
		{
			question: "What happens after the diagnostic if we want more help?",
			answer: "If appropriate, we can discuss the Production Readiness Assessment, which produces a 6–12 month production plan. Totally optional."
		},
		{
			question: "Do you require access to our systems or data?",
			answer: "No. This is a conversational diagnostic, not a technical review."
		},
		{
			question: "Is this relevant if we already have pilots running?",
			answer: "Especially relevant. Most organizations have pilots running—the issue is they're not moving into production."
		},
		{
			question: "What if everything feels chaotic and I'm not sure what our problem is?",
			answer: "That's exactly when to book. Most leaders come to the diagnostic because the situation is unclear."
		}
	];

	function toggle(index: number) {
		openIndex = openIndex === index ? null : index;
	}
</script>

<section class="py-16 md:py-24 bg-white">
	<div class="max-w-4xl mx-auto px-6 md:px-12">
		<h2 class="text-3xl md:text-4xl lg:text-5xl font-black text-gray-900 mb-12 text-center">
			Frequently Asked Questions
		</h2>

		<div class="space-y-4">
			{#each faqs as faq, index}
				<div class="border border-gray-200 rounded-xl overflow-hidden hover:border-blue-300 transition-colors duration-200">
					<button
						class="w-full text-left px-6 py-5 bg-white hover:bg-slate-50 transition-colors flex items-center justify-between gap-4"
						onclick={() => toggle(index)}
						aria-expanded={openIndex === index}
					>
						<span class="text-lg font-semibold text-gray-900">{faq.question}</span>
						<svg
							class="w-6 h-6 text-blue-600 transition-transform duration-300 flex-shrink-0"
							class:rotate-180={openIndex === index}
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>

					{#if openIndex === index}
						<div class="px-6 py-5 bg-slate-50 border-t border-gray-200 animate-fadeIn">
							<div class="space-y-4">
								{#each faq.answer.split('\n\n') as paragraph}
									<p class="text-lg text-gray-700 leading-relaxed">{paragraph}</p>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</section>

<style>
	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.animate-fadeIn {
		animation: fadeIn 0.3s ease-out;
	}
</style>
