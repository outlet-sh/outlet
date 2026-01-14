<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { trackEvent } from '$lib/analytics';
	import { trackFunnelStep } from '$lib/funnel';
	import type { Snippet } from 'svelte';

	type OfferType = 'upsell' | 'downsell' | 'order_bump';

	interface OfferProduct {
		name: string;
		description?: string;
		priceCents: number;
		compareAtCents?: number;
		features?: string[];
		image?: string;
	}

	interface Props {
		offerType: OfferType;
		product: OfferProduct;
		acceptUrl: string;
		declineUrl: string;
		headline?: string;
		subheadline?: string;
		acceptButtonText?: string;
		declineButtonText?: string;
		urgencyText?: string;
		guaranteeText?: string;
		testimonial?: { quote: string; author: string; role?: string };
		children?: Snippet;
		class?: string;
	}

	let {
		offerType,
		product,
		acceptUrl,
		declineUrl,
		headline = offerType === 'upsell' ? 'Wait! Special One-Time Offer' : 'Before You Go...',
		subheadline = '',
		acceptButtonText = 'Yes! Add This To My Order',
		declineButtonText = 'No Thanks, I\'ll Pass',
		urgencyText = '',
		guaranteeText = '',
		testimonial,
		children,
		class: className = ''
	}: Props = $props();

	let isProcessing = $state(false);

	onMount(() => {
		trackFunnelStep(`${offerType}_view`);
		trackEvent(`${offerType}_shown`, {
			product_name: product.name,
			price: product.priceCents
		});
	});

	function formatPrice(cents: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',
			minimumFractionDigits: 0
		}).format(cents / 100);
	}

	async function handleAccept() {
		if (isProcessing) return;
		isProcessing = true;

		trackEvent(`${offerType}_accepted`, {
			product_name: product.name,
			price: product.priceCents
		});

		trackFunnelStep(`${offerType}_accepted`);

		// Navigate to accept URL (which should handle the actual purchase)
		if (browser) {
			// Pass session info if available
			const sessionId = $page.url.searchParams.get('session_id');
			const url = new URL(acceptUrl, window.location.origin);
			if (sessionId) {
				url.searchParams.set('session_id', sessionId);
			}
			goto(url.pathname + url.search);
		}
	}

	async function handleDecline() {
		if (isProcessing) return;
		isProcessing = true;

		trackEvent(`${offerType}_declined`, {
			product_name: product.name
		});

		trackFunnelStep(`${offerType}_declined`);

		if (browser) {
			const sessionId = $page.url.searchParams.get('session_id');
			const url = new URL(declineUrl, window.location.origin);
			if (sessionId) {
				url.searchParams.set('session_id', sessionId);
			}
			goto(url.pathname + url.search);
		}
	}

	let savings = $derived(
		product.compareAtCents && product.compareAtCents > product.priceCents
			? product.compareAtCents - product.priceCents
			: 0
	);

	let discountPercent = $derived(
		product.compareAtCents && product.compareAtCents > product.priceCents
			? Math.round((1 - product.priceCents / product.compareAtCents) * 100)
			: 0
	);
</script>

<div class="offer-page min-h-screen bg-gradient-to-b from-slate-50 to-white py-12 px-4 {className}">
	<div class="max-w-2xl mx-auto">
		<!-- Urgency banner -->
		{#if urgencyText}
			<div class="mb-6 bg-amber-500 text-white text-center py-3 px-4 rounded-lg font-semibold animate-pulse">
				{urgencyText}
			</div>
		{/if}

		<!-- Main offer card -->
		<div class="bg-white rounded-2xl shadow-2xl overflow-hidden border border-slate-200">
			<!-- Header -->
			<div class="bg-gradient-to-r from-blue-600 to-violet-600 text-white p-6 text-center">
				{#if offerType === 'upsell'}
					<div class="text-sm font-medium opacity-90 mb-1">EXCLUSIVE ONE-TIME OFFER</div>
				{:else if offerType === 'downsell'}
					<div class="text-sm font-medium opacity-90 mb-1">SPECIAL OFFER JUST FOR YOU</div>
				{/if}
				<h1 class="text-2xl md:text-3xl font-bold">{headline}</h1>
				{#if subheadline}
					<p class="mt-2 opacity-90">{subheadline}</p>
				{/if}
			</div>

			<!-- Product info -->
			<div class="p-6 md:p-8">
				{#if product.image}
					<div class="mb-6">
						<img src={product.image} alt={product.name} class="w-full max-w-md mx-auto rounded-lg shadow" />
					</div>
				{/if}

				<h2 class="text-xl md:text-2xl font-bold text-slate-900 text-center mb-4">{product.name}</h2>

				{#if product.description}
					<p class="text-slate-600 text-center mb-6">{product.description}</p>
				{/if}

				<!-- Features -->
				{#if product.features && product.features.length > 0}
					<div class="mb-6 bg-slate-50 rounded-xl p-6">
						<p class="font-semibold text-slate-900 mb-3">What's Included:</p>
						<ul class="space-y-2">
							{#each product.features as feature}
								<li class="flex items-start gap-3">
									<svg class="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
									</svg>
									<span class="text-slate-700">{feature}</span>
								</li>
							{/each}
						</ul>
					</div>
				{/if}

				<!-- Custom content slot -->
				{#if children}
					<div class="mb-6">
						{@render children()}
					</div>
				{/if}

				<!-- Testimonial -->
				{#if testimonial}
					<div class="mb-6 bg-blue-50 rounded-xl p-6 border border-blue-100">
						<svg class="w-8 h-8 text-blue-400 mb-2" fill="currentColor" viewBox="0 0 24 24">
							<path d="M14.017 21v-7.391c0-5.704 3.731-9.57 8.983-10.609l.995 2.151c-2.432.917-3.995 3.638-3.995 5.849h4v10h-9.983zm-14.017 0v-7.391c0-5.704 3.748-9.57 9-10.609l.996 2.151c-2.433.917-3.996 3.638-3.996 5.849h3.983v10h-9.983z"/>
						</svg>
						<p class="text-slate-700 italic mb-3">"{testimonial.quote}"</p>
						<p class="text-slate-900 font-semibold">{testimonial.author}</p>
						{#if testimonial.role}
							<p class="text-slate-500 text-sm">{testimonial.role}</p>
						{/if}
					</div>
				{/if}

				<!-- Pricing -->
				<div class="text-center mb-6">
					{#if savings > 0}
						<div class="mb-2">
							<span class="text-slate-400 line-through text-lg">{formatPrice(product.compareAtCents || 0)}</span>
							<span class="ml-2 bg-green-100 text-green-800 text-sm font-semibold px-2 py-1 rounded">
								Save {discountPercent}%
							</span>
						</div>
					{/if}
					<div class="text-4xl md:text-5xl font-black text-slate-900">
						{formatPrice(product.priceCents)}
					</div>
					<p class="text-slate-500 mt-1">one-time payment</p>
				</div>

				<!-- Guarantee -->
				{#if guaranteeText}
					<div class="mb-6 text-center">
						<div class="inline-flex items-center gap-2 text-green-700 bg-green-50 px-4 py-2 rounded-lg">
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
							</svg>
							<span class="font-medium">{guaranteeText}</span>
						</div>
					</div>
				{/if}

				<!-- Action buttons -->
				<div class="space-y-3">
					<button
						type="button"
						onclick={handleAccept}
						disabled={isProcessing}
						class="w-full py-4 px-6 bg-gradient-to-r from-green-500 to-emerald-600 text-white font-bold text-lg rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{#if isProcessing}
							Processing...
						{:else}
							{acceptButtonText}
						{/if}
					</button>

					<button
						type="button"
						onclick={handleDecline}
						disabled={isProcessing}
						class="w-full py-3 px-6 text-slate-500 hover:text-slate-700 text-sm underline transition-colors disabled:opacity-50"
					>
						{declineButtonText}
					</button>
				</div>
			</div>
		</div>

		<!-- Trust indicators -->
		<div class="mt-8 flex flex-wrap justify-center gap-6 text-sm text-slate-500">
			<div class="flex items-center gap-2">
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
				</svg>
				<span>Secure Payment</span>
			</div>
			<div class="flex items-center gap-2">
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
				</svg>
				<span>Instant Access</span>
			</div>
			<div class="flex items-center gap-2">
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192l-3.536 3.536M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z" />
				</svg>
				<span>Support Included</span>
			</div>
		</div>
	</div>
</div>
