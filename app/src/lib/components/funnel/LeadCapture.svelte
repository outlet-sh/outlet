<script context="module" lang="ts">
	// Declare gtag for TypeScript
	declare const gtag: (command: string, event: string, params: Record<string, unknown>) => void;
</script>

<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { trackEvent, trackConversion } from '$lib/analytics';
	import { trackFunnelStep } from '$lib/funnel';
	import { createContact, type ContactRequest } from '$lib/api';

	interface Props {
		funnel: string;
		redirectTo: string;
		headline?: string;
		subheadline?: string;
		buttonText?: string;
		fields?: ('name' | 'email' | 'company' | 'role')[];
		/** Visual variant: 'light' (white bg), 'dark' (dark inputs) */
		variant?: 'light' | 'dark';
		/** Button style variant */
		buttonVariant?: 'gradient-blue' | 'gradient-orange' | 'solid-cyan';
		/** Capture UTM parameters from URL */
		captureUtm?: boolean;
		/** Track form location for analytics */
		formLocation?: string;
		class?: string;
	}

	let {
		funnel,
		redirectTo,
		headline = '',
		subheadline = '',
		buttonText = 'Submit',
		fields = ['name', 'email'],
		variant = 'light',
		buttonVariant = 'gradient-blue',
		captureUtm = true,
		formLocation = 'lead_capture',
		class: className = ''
	}: Props = $props();

	let isSubmitting = $state(false);
	let errorMessage = $state('');
	let formStarted = false;

	let formData = $state({
		name: '',
		email: '',
		company: '',
		role: ''
	});

	function trackFormStart() {
		if (!formStarted) {
			formStarted = true;
			trackEvent('form_started', {
				form: 'lead_capture',
				funnel
			});
		}
	}

	function canSubmit(): boolean {
		if (fields.includes('name') && !formData.name) return false;
		if (fields.includes('email') && !formData.email) return false;
		return true;
	}

	async function handleSubmit() {
		if (!canSubmit()) return;

		isSubmitting = true;
		errorMessage = '';

		try {
			// Capture UTM params from URL if enabled
			const utmParams = captureUtm ? {
				utm_source: $page.url.searchParams.get('utm_source') || undefined,
				utm_medium: $page.url.searchParams.get('utm_medium') || undefined,
				utm_campaign: $page.url.searchParams.get('utm_campaign') || undefined
			} : {};

			const contactRequest: ContactRequest = {
				email: formData.email,
				name: formData.name || undefined,
				company: formData.company || undefined,
				funnel_slug: funnel,
				...utmParams
			};

			const result = await createContact(contactRequest);

			if (result.id) {
				trackConversion('lead_capture', {
					funnel,
					lead_id: result.id,
					location: formLocation
				});

				trackFunnelStep('optin_complete', funnel);

				trackEvent('funnel_opt_in', {
					funnel,
					location: formLocation
				});

				// Google Ads conversion tracking
				if (browser && typeof gtag !== 'undefined') {
					gtag('event', 'conversion', {
						send_to: 'AW-17764059193/lead_capture'
					});
				}

				// Redirect to next step
				goto(redirectTo);
			} else {
				errorMessage = 'Something went wrong. Please try again.';
			}
		} catch (error: any) {
			console.error('Lead capture error:', error);
			errorMessage = error?.message || 'Something went wrong. Please try again.';
		} finally {
			isSubmitting = false;
		}
	}
</script>

<div class="lead-capture {className}">
	{#if headline}
		<h2 class="{variant === 'dark' ? 'text-white' : 'text-slate-900'} text-2xl md:text-3xl font-bold mb-2 text-center">{headline}</h2>
	{/if}
	{#if subheadline}
		<p class="{variant === 'dark' ? 'text-slate-300' : 'text-slate-600'} mb-6 text-center">{subheadline}</p>
	{/if}

	<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
		{#if fields.includes('name')}
			<div>
				{#if variant === 'light'}
					<label for="lead-name" class="block text-sm font-semibold text-slate-900 mb-2 text-left">Name*</label>
				{/if}
				<input
					type="text"
					id="lead-name"
					name="name"
					autocomplete="name"
					placeholder={variant === 'dark' ? 'Your name' : ''}
					bind:value={formData.name}
					onfocus={trackFormStart}
					class="w-full px-4 py-3 rounded-lg text-slate-900 focus:outline-none focus:ring-2 transition-all {variant === 'dark' ? 'bg-white border-0' : 'bg-white border border-slate-200 focus:ring-blue-500/20 focus:border-blue-500'}"
				/>
			</div>
		{/if}

		{#if fields.includes('email')}
			<div>
				{#if variant === 'light'}
					<label for="lead-email" class="block text-sm font-semibold text-slate-900 mb-2 text-left">Email*</label>
				{/if}
				<input
					type="email"
					id="lead-email"
					name="email"
					autocomplete="email"
					placeholder={variant === 'dark' ? 'Your work email' : ''}
					bind:value={formData.email}
					onfocus={trackFormStart}
					class="w-full px-4 py-3 rounded-lg text-slate-900 focus:outline-none focus:ring-2 transition-all {variant === 'dark' ? 'bg-white border-0' : 'bg-white border border-slate-200 focus:ring-blue-500/20 focus:border-blue-500'}"
				/>
			</div>
		{/if}

		{#if fields.includes('company')}
			<div>
				{#if variant === 'light'}
					<label for="lead-company" class="block text-sm font-semibold text-slate-900 mb-2 text-left">Company</label>
				{/if}
				<input
					type="text"
					id="lead-company"
					name="company"
					autocomplete="organization"
					placeholder={variant === 'dark' ? 'Company' : ''}
					bind:value={formData.company}
					onfocus={trackFormStart}
					class="w-full px-4 py-3 rounded-lg text-slate-900 focus:outline-none focus:ring-2 transition-all {variant === 'dark' ? 'bg-white border-0' : 'bg-white border border-slate-200 focus:ring-blue-500/20 focus:border-blue-500'}"
				/>
			</div>
		{/if}

		{#if fields.includes('role')}
			<div>
				{#if variant === 'light'}
					<label for="lead-role" class="block text-sm font-semibold text-slate-900 mb-2 text-left">Your Role</label>
				{/if}
				<input
					type="text"
					id="lead-role"
					name="role"
					placeholder={variant === 'dark' ? 'Your role' : ''}
					bind:value={formData.role}
					onfocus={trackFormStart}
					class="w-full px-4 py-3 rounded-lg text-slate-900 focus:outline-none focus:ring-2 transition-all {variant === 'dark' ? 'bg-white border-0' : 'bg-white border border-slate-200 focus:ring-blue-500/20 focus:border-blue-500'}"
				/>
			</div>
		{/if}

		{#if errorMessage}
			<div class="p-3 bg-red-50 border border-red-200 rounded-lg">
				<p class="text-sm text-red-600">{errorMessage}</p>
			</div>
		{/if}

		<button
			type="submit"
			disabled={!canSubmit() || isSubmitting}
			class="lead-capture-btn lead-capture-btn--{buttonVariant} w-full py-4 px-6 text-white font-bold rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:translate-y-0 flex items-center justify-center gap-2"
		>
			{#if isSubmitting}
				<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
				</svg>
				Processing...
			{:else}
				{buttonText}
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
				</svg>
			{/if}
		</button>
	</form>
</div>
