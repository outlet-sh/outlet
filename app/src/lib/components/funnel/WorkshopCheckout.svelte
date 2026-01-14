<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { trackEvent } from '$lib/analytics';
	import { trackFunnelStep } from '$lib/funnel';
	import { sDKGetWorkshop, sDKGetProduct, createCheckoutSession, type WorkshopEventInfo, type BillingProductInfo, type CheckoutRequest } from '$lib/api';

	interface BumpOption {
		productSlug: string;
		label: string;
		description?: string;
		priceCents: number;
	}

	interface Props {
		eventSlug: string;
		bumps?: BumpOption[];
		successRedirect?: string;
		class?: string;
	}

	let {
		eventSlug,
		bumps = [],
		successRedirect = '/workshop/success',
		class: className = ''
	}: Props = $props();

	// Step state: 1 = registration, 2 = payment
	let step = $state(1);

	let loading = $state(true);
	let workshop = $state<WorkshopEventInfo | null>(null);
	let product = $state<BillingProductInfo | null>(null);
	let loadError = $state('');

	let isSubmitting = $state(false);
	let checkoutError = $state('');

	// Registration form
	let formData = $state({
		name: '',
		email: '',
		company: '',
		role: ''
	});

	// Selected bumps
	let selectedBumps = $state<string[]>([]);

	onMount(async () => {
		await loadWorkshop();
	});

	async function loadWorkshop() {
		loading = true;
		loadError = '';

		try {
			workshop = await sDKGetWorkshop({}, eventSlug);

			// Also load the product to get stripe_price_id
			if (workshop?.product_slug) {
				product = await sDKGetProduct({}, workshop.product_slug);
			}

			trackFunnelStep('workshop_checkout_view', eventSlug);
		} catch (err: any) {
			console.error('Failed to load workshop:', err);
			loadError = err.message || 'Failed to load workshop details';
		} finally {
			loading = false;
		}
	}

	function formatPrice(cents: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',
			minimumFractionDigits: 0
		}).format(cents / 100);
	}

	function formatDateRange(startDate: string, endDate: string): string {
		const start = new Date(startDate);
		const end = new Date(endDate);

		if (start.getMonth() === end.getMonth() && start.getFullYear() === end.getFullYear()) {
			return `${start.toLocaleDateString('en-US', { month: 'long' })} ${start.getDate()}-${end.getDate()}, ${start.getFullYear()}`;
		}

		return `${start.toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })} - ${end.toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })}`;
	}

	function canProceedToPayment(): boolean {
		return !!formData.name && !!formData.email;
	}

	function handleProceedToPayment() {
		if (!canProceedToPayment()) return;

		trackEvent('workshop_registration_step_complete', {
			event_slug: eventSlug,
			has_bumps_selected: selectedBumps.length > 0
		});

		step = 2;
	}

	function toggleBump(productSlug: string) {
		if (selectedBumps.includes(productSlug)) {
			selectedBumps = selectedBumps.filter(s => s !== productSlug);
		} else {
			selectedBumps = [...selectedBumps, productSlug];
		}
	}

	function calculateTotal(): number {
		if (!workshop) return 0;
		let total = workshop.price_cents;
		for (const slug of selectedBumps) {
			const bump = bumps.find(b => b.productSlug === slug);
			if (bump) total += bump.priceCents;
		}
		return total;
	}

	// Helper to get the default price from a product
	function getDefaultPrice(prod: BillingProductInfo | null) {
		return prod?.prices?.[0];
	}

	async function handleCheckout() {
		const price = getDefaultPrice(product);
		if (!workshop || !price?.stripe_price_id) return;

		isSubmitting = true;
		checkoutError = '';

		try {
			trackEvent('workshop_checkout_started', {
				event_slug: eventSlug,
				price: workshop.price_cents,
				bumps_selected: selectedBumps
			});

			const checkoutRequest: CheckoutRequest = {
				customer_email: formData.email,
				line_items: [{
					price_id: price.stripe_price_id!,
					quantity: 1
				}],
				mode: 'payment',
				success_url: `${window.location.origin}${successRedirect}`,
				cancel_url: window.location.href
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
			trackEvent('workshop_checkout_error', { error: checkoutError });
		} finally {
			isSubmitting = false;
		}
	}
</script>

<div class="workshop-checkout {className}">
	{#if loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-10 w-10 border-b-2 border-cyan-600"></div>
		</div>
	{:else if loadError}
		<div class="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
			<p class="text-red-800">{loadError}</p>
		</div>
	{:else if workshop}
		<!-- Progress indicator -->
		<div class="flex items-center justify-center mb-8">
			<div class="flex items-center">
				<div class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold {step >= 1 ? 'bg-cyan-600 text-white' : 'bg-slate-200 text-slate-500'}">
					1
				</div>
				<div class="w-16 h-1 {step >= 2 ? 'bg-cyan-600' : 'bg-slate-200'}"></div>
				<div class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold {step >= 2 ? 'bg-cyan-600 text-white' : 'bg-slate-200 text-slate-500'}">
					2
				</div>
			</div>
		</div>

		{#if step === 1}
			<!-- Step 1: Registration -->
			<div class="bg-white rounded-xl border border-slate-200 p-6">
				<h3 class="text-lg font-bold text-slate-900 mb-4">Your Information</h3>

				<div class="space-y-4">
					<div>
						<label for="ws-name" class="block text-sm font-semibold text-slate-700 mb-1">Full Name *</label>
						<input
							type="text"
							id="ws-name"
							bind:value={formData.name}
							required
							class="w-full px-4 py-3 rounded-lg border border-slate-300 focus:border-cyan-500 focus:ring-2 focus:ring-cyan-200 transition-colors"
							placeholder="John Smith"
						/>
					</div>

					<div>
						<label for="ws-email" class="block text-sm font-semibold text-slate-700 mb-1">Work Email *</label>
						<input
							type="email"
							id="ws-email"
							bind:value={formData.email}
							required
							class="w-full px-4 py-3 rounded-lg border border-slate-300 focus:border-cyan-500 focus:ring-2 focus:ring-cyan-200 transition-colors"
							placeholder="john@company.com"
						/>
					</div>

					<div>
						<label for="ws-company" class="block text-sm font-semibold text-slate-700 mb-1">Company</label>
						<input
							type="text"
							id="ws-company"
							bind:value={formData.company}
							class="w-full px-4 py-3 rounded-lg border border-slate-300 focus:border-cyan-500 focus:ring-2 focus:ring-cyan-200 transition-colors"
							placeholder="Acme Corp"
						/>
					</div>

					<div>
						<label for="ws-role" class="block text-sm font-semibold text-slate-700 mb-1">Your Role</label>
						<input
							type="text"
							id="ws-role"
							bind:value={formData.role}
							class="w-full px-4 py-3 rounded-lg border border-slate-300 focus:border-cyan-500 focus:ring-2 focus:ring-cyan-200 transition-colors"
							placeholder="VP of Engineering"
						/>
					</div>

					<button
						type="button"
						onclick={handleProceedToPayment}
						disabled={!canProceedToPayment()}
						class="w-full py-4 px-6 bg-cyan-600 text-white font-bold rounded-lg hover:bg-cyan-700 disabled:bg-slate-400 disabled:cursor-not-allowed transition-colors"
					>
						Continue to Payment
					</button>
				</div>
			</div>
		{:else}
			<!-- Step 2: Payment -->
			<div class="space-y-6">
				<!-- Order summary -->
				<div class="bg-slate-50 rounded-xl border border-slate-200 p-6">
					<h3 class="text-lg font-bold text-slate-900 mb-4">Order Summary</h3>

					<div class="space-y-3">
						<div class="flex justify-between items-start">
							<div>
								<p class="font-semibold text-slate-900">{workshop.title}</p>
								<p class="text-sm text-slate-500">{formatDateRange(workshop.start_date, workshop.end_date)}</p>
							</div>
							<p class="font-semibold text-slate-900">{formatPrice(workshop.price_cents)}</p>
						</div>

						{#each selectedBumps as slug}
							{@const bump = bumps.find(b => b.productSlug === slug)}
							{#if bump}
								<div class="flex justify-between text-sm">
									<span class="text-slate-700">{bump.label}</span>
									<span class="text-slate-900">{formatPrice(bump.priceCents)}</span>
								</div>
							{/if}
						{/each}

						<div class="border-t border-slate-200 pt-3 flex justify-between">
							<span class="font-bold text-slate-900">Total</span>
							<span class="font-bold text-slate-900">{formatPrice(calculateTotal())}</span>
						</div>
					</div>
				</div>

				<!-- Bumps (if any) -->
				{#if bumps.length > 0}
					<div class="bg-amber-50 rounded-xl border-2 border-amber-200 p-6">
						<h4 class="font-bold text-amber-900 mb-4">Special Offers</h4>
						<div class="space-y-3">
							{#each bumps as bump}
								<label class="flex items-start gap-3 cursor-pointer p-3 rounded-lg hover:bg-amber-100/50 transition-colors">
									<input
										type="checkbox"
										checked={selectedBumps.includes(bump.productSlug)}
										onchange={() => toggleBump(bump.productSlug)}
										class="mt-1 w-5 h-5 rounded border-amber-400 text-amber-600 focus:ring-amber-500"
									/>
									<div class="flex-1">
										<p class="font-semibold text-amber-900">{bump.label}</p>
										{#if bump.description}
											<p class="text-sm text-amber-700">{bump.description}</p>
										{/if}
									</div>
									<span class="font-semibold text-amber-900">+{formatPrice(bump.priceCents)}</span>
								</label>
							{/each}
						</div>
					</div>
				{/if}

				<!-- Seats remaining warning -->
				{#if workshop.seats_remaining <= 4 && workshop.seats_remaining > 0}
					<div class="bg-amber-50 border border-amber-200 rounded-lg p-3 text-center">
						<span class="text-amber-800 font-semibold">Only {workshop.seats_remaining} seats remaining</span>
					</div>
				{/if}

				{#if checkoutError}
					<div class="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
						{checkoutError}
					</div>
				{/if}

				<div class="flex gap-3">
					<button
						type="button"
						onclick={() => step = 1}
						class="px-6 py-4 border border-slate-300 text-slate-700 font-semibold rounded-lg hover:bg-slate-50 transition-colors"
					>
						Back
					</button>
					<button
						type="button"
						onclick={handleCheckout}
						disabled={isSubmitting || !getDefaultPrice(product)?.stripe_price_id}
						class="flex-1 py-4 px-6 bg-cyan-600 text-white font-bold rounded-lg hover:bg-cyan-700 disabled:bg-slate-400 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
					>
						{#if isSubmitting}
							<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							Processing...
						{:else}
							Complete Purchase - {formatPrice(calculateTotal())}
						{/if}
					</button>
				</div>

				<!-- Trust indicators -->
				<div class="text-center space-y-2 text-sm text-slate-500">
					<div class="flex items-center justify-center gap-2">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
						</svg>
						<span>Secure payment via Stripe</span>
					</div>
					<p>Calendar invite sent immediately after purchase</p>
				</div>
			</div>
		{/if}
	{/if}
</div>
