<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { trackEvent } from '$lib/analytics';
	import { trackFunnelStep } from '$lib/funnel';
	import { sDKGetProduct, createCheckoutSession, type BillingProductInfo, type BillingPriceInfo, type CheckoutRequest } from '$lib/api';
	import type { Snippet } from 'svelte';

	// Helper to get the default price from a product
	function getDefaultPrice(prod: BillingProductInfo | null): BillingPriceInfo | undefined {
		return prod?.prices?.[0];
	}

	interface Props {
		productSlug: string;
		successRedirect?: string;
		cancelRedirect?: string;
		children?: Snippet;
		class?: string;
	}

	let {
		productSlug,
		successRedirect = '/thank-you',
		cancelRedirect,
		children,
		class: className = ''
	}: Props = $props();

	let loading = $state(true);
	let product = $state<BillingProductInfo | null>(null);
	let loadError = $state('');

	let isCheckingOut = $state(false);
	let checkoutError = $state('');

	// Form for checkout
	let formData = $state({
		name: '',
		email: ''
	});

	// Derived default price for template use
	let defaultPrice = $derived(getDefaultPrice(product));

	onMount(async () => {
		await loadProduct();
		trackFunnelStep('product_page_view', `/product/${productSlug}`);
	});

	async function loadProduct() {
		loading = true;
		loadError = '';

		try {
			product = await sDKGetProduct({}, productSlug);
			if (!product) {
				loadError = 'Product not found';
			}
		} catch (err: any) {
			console.error('Failed to load product:', err);
			loadError = err.message || 'Failed to load product details';
		} finally {
			loading = false;
		}
	}

	function formatPrice(cents: number): string {
		const price = getDefaultPrice(product);
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: price?.currency || 'USD',
			minimumFractionDigits: 0
		}).format(cents / 100);
	}

	function isInStock(): boolean {
		if (!product) return false;
		if (product.inventory_type === 'unlimited') return true;
		if (product.inventory_type === 'limited' && product.inventory_count && product.inventory_count > 0) return true;
		return false;
	}

	function canCheckout(): boolean {
		return !!formData.name && !!formData.email && isInStock();
	}

	async function handleCheckout() {
		const price = getDefaultPrice(product);
		if (!product || !price || !canCheckout()) return;

		isCheckingOut = true;
		checkoutError = '';

		try {
			trackEvent('product_checkout_started', {
				product_slug: productSlug,
				price: price.unit_amount_cents
			});

			if (!price.stripe_price_id) {
				throw new Error('Product is not available for purchase');
			}

			const checkoutRequest: CheckoutRequest = {
				customer_email: formData.email,
				line_items: [{
					price_id: price.stripe_price_id,
					quantity: 1
				}],
				mode: 'payment',
				success_url: `${window.location.origin}${successRedirect}`,
				cancel_url: cancelRedirect ? `${window.location.origin}${cancelRedirect}` : window.location.href
			};

			const response = await createCheckoutSession(checkoutRequest);

			if (!response.checkout_url) {
				throw new Error('Failed to create checkout session');
			}

			// Redirect to Stripe Checkout
			if (browser) {
				window.location.href = response.checkout_url;
			}
		} catch (err: any) {
			checkoutError = err.message || 'Something went wrong. Please try again.';
			trackEvent('product_checkout_error', { error: checkoutError });
		} finally {
			isCheckingOut = false;
		}
	}
</script>

<div class="product-page {className}">
	{#if loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-600"></div>
		</div>
	{:else if loadError}
		<div class="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
			<p class="text-red-800">{loadError}</p>
		</div>
	{:else if product}
		<!-- Slot for custom sales content -->
		{#if children}
			{@render children()}
		{/if}

		<!-- Checkout section -->
		<div class="mt-12 max-w-xl mx-auto">
			<div class="bg-white rounded-2xl border-2 border-slate-200 shadow-lg p-8">
				<!-- Price display -->
				<div class="text-center mb-6">
					{#if defaultPrice?.compare_at_cents && defaultPrice.compare_at_cents > defaultPrice.unit_amount_cents}
						<p class="text-slate-400 line-through text-lg">{formatPrice(defaultPrice.compare_at_cents)}</p>
					{/if}
					<p class="text-4xl font-black text-slate-900">{formatPrice(defaultPrice?.unit_amount_cents ?? 0)}</p>
					<p class="text-slate-500">one-time payment</p>
				</div>

				<!-- What's included (if provided) -->
				{#if product.whats_included && product.whats_included.length > 0}
					<div class="mb-6 space-y-2">
						{#each product.whats_included as item}
							<div class="flex items-center gap-2 text-slate-700">
								<svg class="w-5 h-5 text-emerald-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
								</svg>
								<span>{item}</span>
							</div>
						{/each}
					</div>
				{/if}

				<!-- Features (if provided) -->
				{#if product.features && product.features.length > 0}
					<div class="mb-6 p-4 bg-slate-50 rounded-lg">
						<p class="text-sm font-semibold text-slate-700 mb-2">Features:</p>
						<ul class="space-y-1 text-sm text-slate-600">
							{#each product.features as feature}
								<li class="flex items-start gap-2">
									<span class="text-blue-500">•</span>
									<span>{feature}</span>
								</li>
							{/each}
						</ul>
					</div>
				{/if}

				<!-- Out of stock message -->
				{#if !isInStock()}
					<div class="mb-6 p-4 bg-amber-50 border border-amber-200 rounded-lg text-center">
						<p class="text-amber-800 font-semibold">Currently Sold Out</p>
						<p class="text-amber-700 text-sm">Join the waitlist to be notified when available.</p>
					</div>
				{/if}

				<!-- Checkout form -->
				<form onsubmit={(e) => { e.preventDefault(); handleCheckout(); }} class="space-y-4">
					<div>
						<label for="product-name" class="block text-sm font-semibold text-slate-700 mb-1">Name *</label>
						<input
							type="text"
							id="product-name"
							bind:value={formData.name}
							required
							class="w-full px-4 py-3 rounded-lg border border-slate-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-colors"
							placeholder="Your name"
						/>
					</div>

					<div>
						<label for="product-email" class="block text-sm font-semibold text-slate-700 mb-1">Email *</label>
						<input
							type="email"
							id="product-email"
							bind:value={formData.email}
							required
							class="w-full px-4 py-3 rounded-lg border border-slate-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-colors"
							placeholder="you@example.com"
						/>
					</div>

					{#if checkoutError}
						<div class="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
							{checkoutError}
						</div>
					{/if}

					<button
						type="submit"
						disabled={!canCheckout() || isCheckingOut}
						class="w-full py-4 px-6 bg-blue-600 text-white font-bold rounded-lg hover:bg-blue-700 disabled:bg-slate-400 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
					>
						{#if isCheckingOut}
							<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							Processing...
						{:else if !isInStock()}
							Join Waitlist
						{:else}
							Get Access - {formatPrice(defaultPrice?.unit_amount_cents ?? 0)}
						{/if}
					</button>
				</form>

				<!-- Trust indicators -->
				<div class="mt-6 pt-6 border-t border-slate-200 space-y-2 text-sm text-slate-500 text-center">
					<div class="flex items-center justify-center gap-2">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
						</svg>
						<span>Secure payment via Stripe</span>
					</div>
					<p>Instant access after purchase</p>
				</div>
			</div>
		</div>
	{/if}
</div>
