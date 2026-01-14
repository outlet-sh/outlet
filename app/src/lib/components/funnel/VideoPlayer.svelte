<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { trackEvent } from '$lib/analytics';
	import type { Snippet } from 'svelte';

	interface Props {
		/** YouTube embed URL (e.g., https://www.youtube.com/embed/VIDEO_ID) */
		src: string;
		/** Seconds to wait before revealing CTA slot content */
		revealCtaAfter?: number;
		/** Video title for accessibility */
		title?: string;
		/** Aspect ratio class */
		aspectRatio?: 'video' | '16/9' | '4/3';
		/** Show countdown timer (set false for B2B) */
		showCountdown?: boolean;
		/** Slot for CTA content that appears after timer */
		cta?: Snippet;
		class?: string;
	}

	let {
		src,
		revealCtaAfter = 0,
		title = 'Video',
		aspectRatio = 'video',
		showCountdown = false,
		cta,
		class: className = ''
	}: Props = $props();

	let showCta = $state(revealCtaAfter === 0);
	let secondsWatched = $state(0);
	let timerInterval: ReturnType<typeof setInterval> | null = null;

	onMount(() => {
		if (revealCtaAfter > 0) {
			timerInterval = setInterval(() => {
				secondsWatched++;
				if (secondsWatched >= revealCtaAfter && !showCta) {
					showCta = true;
					trackEvent('video_cta_revealed', {
						src,
						seconds_watched: secondsWatched
					});
				}
			}, 1000);
		}

		trackEvent('video_player_loaded', { src });
	});

	onDestroy(() => {
		if (timerInterval) {
			clearInterval(timerInterval);
		}
	});

	function getAspectRatioClass(): string {
		switch (aspectRatio) {
			case '16/9':
				return 'aspect-[16/9]';
			case '4/3':
				return 'aspect-[4/3]';
			default:
				return 'aspect-video';
		}
	}

	function formatTime(seconds: number): string {
		const mins = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${mins}:${secs.toString().padStart(2, '0')}`;
	}
</script>

<div class="video-player {className}">
	<!-- Video embed -->
	<div class="relative {getAspectRatioClass()} w-full rounded-xl overflow-hidden bg-slate-900 shadow-lg">
		<iframe
			{src}
			{title}
			class="absolute inset-0 w-full h-full"
			frameborder="0"
			allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
			allowfullscreen
		></iframe>
	</div>

	<!-- CTA section -->
	{#if cta}
		<div class="mt-6">
			{#if !showCta && revealCtaAfter > 0 && showCountdown}
				<!-- Countdown until CTA reveals (consumer funnel mode) -->
				<div class="text-center py-4">
					<p class="text-slate-500 text-sm">
						Special offer unlocks in {formatTime(revealCtaAfter - secondsWatched)}
					</p>
					<div class="mt-2 w-full max-w-xs mx-auto bg-slate-200 rounded-full h-2">
						<div
							class="bg-cyan-600 h-2 rounded-full transition-all duration-1000"
							style="width: {Math.min(100, (secondsWatched / revealCtaAfter) * 100)}%"
						></div>
					</div>
				</div>
			{:else if showCta}
				<!-- CTA content revealed -->
				<div class="animate-fade-in">
					{@render cta()}
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	@keyframes fade-in {
		from {
			opacity: 0;
			transform: translateY(10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.animate-fade-in {
		animation: fade-in 0.5s ease-out;
	}
</style>
