<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import * as api from '$lib/api';
	import type { OrgInfo } from '$lib/api';
	import { Card, Input, Button, Alert, LoadingSpinner } from '$lib/components/ui';
	import { ChevronRight, ChevronLeft, Mail, Check, Loader2, Building2, HelpCircle, Lightbulb, Shield } from 'lucide-svelte';

	type Step = 'welcome' | 'brand' | 'complete';
	let currentStep = $state<Step>('welcome');
	let loading = $state(false);
	let pageLoading = $state(true);
	let error = $state('');
	let org = $state<OrgInfo | null>(null);

	// Brand form fields
	let fromName = $state('');
	let fromEmail = $state('');
	let replyTo = $state('');

	// Get org slug from URL
	let orgSlug = $derived($page.params.orgSlug);

	$effect(() => {
		if (orgSlug) {
			loadOrg();
		}
	});

	async function loadOrg() {
		pageLoading = true;
		try {
			org = await api.getOrganizationBySlug(orgSlug);
			// Pre-fill form with existing values
			if (org) {
				fromName = org.from_name || org.name;
				fromEmail = org.from_email || '';
				replyTo = org.reply_to || '';
			}
		} catch (err) {
			console.error('Failed to load organization:', err);
			error = 'Organization not found';
		} finally {
			pageLoading = false;
		}
	}

	function nextStep() {
		if (currentStep === 'welcome') {
			currentStep = 'brand';
		} else if (currentStep === 'brand') {
			saveAndComplete();
		}
	}

	function prevStep() {
		if (currentStep === 'brand') {
			currentStep = 'welcome';
		}
	}

	async function saveAndComplete() {
		if (!org) return;

		loading = true;
		error = '';

		try {
			await api.updateOrgEmailSettings(
				{},
				{
					from_name: fromName,
					from_email: fromEmail,
					reply_to: replyTo || fromEmail
				},
				org.id
			);

			currentStep = 'complete';
		} catch (err: any) {
			console.error('Failed to save org email settings:', err);
			error = err.message || 'Failed to save email settings';
		} finally {
			loading = false;
		}
	}

	function goToDashboard() {
		goto(`/${orgSlug}`);
	}

	function skipForNow() {
		localStorage.setItem(`orgSetupSkipped_${orgSlug}`, 'true');
		goto(`/${orgSlug}`);
	}
</script>

<svelte:head>
	<title>Welcome - {org?.name || 'Setup'}</title>
</svelte:head>

{#if pageLoading}
	<Card>
		<div class="flex items-center justify-center py-12">
			<LoadingSpinner size="large" />
		</div>
	</Card>
{:else if !org}
	<Card>
		<Alert type="error" title="Organization not found">
			<p>The organization you're looking for doesn't exist.</p>
		</Alert>
		<div class="mt-4">
			<Button type="secondary" onclick={() => goto('/')}>Go Back</Button>
		</div>
	</Card>
{:else}
	<div class="grid grid-cols-1 lg:grid-cols-5 gap-8">
		<!-- Main Content -->
		<div class="lg:col-span-3">
			<Card>
				<!-- Progress indicator -->
				<div class="mb-8">
					<div class="flex items-center justify-between text-sm">
						<span class="text-text-muted">Step {currentStep === 'welcome' ? 1 : currentStep === 'brand' ? 2 : 3} of 3</span>
					</div>
					<div class="mt-2 h-2 bg-border rounded-full overflow-hidden">
						<div
							class="h-full bg-primary transition-all duration-300"
							style="width: {currentStep === 'welcome' ? '33%' : currentStep === 'brand' ? '66%' : '100%'}"
						></div>
					</div>
				</div>

				{#if currentStep === 'welcome'}
					<!-- Welcome Step -->
					<div class="text-center">
						<div class="mx-auto w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mb-6">
							<Building2 class="w-8 h-8 text-primary" />
						</div>
						<h1 class="text-2xl font-bold text-text mb-2">Welcome to {org.name}</h1>
						<p class="text-text-muted mb-8">
							Let's set up your email sender identity so your subscribers know who's sending them emails.
						</p>

						<div class="space-y-4 text-left mb-8">
							<div class="flex items-start gap-3 p-4 bg-bg-muted rounded-lg">
								<Mail class="w-5 h-5 text-primary mt-0.5" />
								<div>
									<p class="font-medium text-text">Brand Setup</p>
									<p class="text-sm text-text-muted">Configure your from name and email address</p>
								</div>
							</div>
						</div>

						<div class="flex gap-3">
							<Button type="primary" size="lg" onclick={nextStep} class="flex-1 justify-center">
								Get Started
								<ChevronRight class="w-4 h-4 ml-2" />
							</Button>
						</div>

						<button
							type="button"
							onclick={skipForNow}
							class="mt-4 text-sm text-text-muted hover:text-text transition-colors"
						>
							Skip for now
						</button>
					</div>
				{:else if currentStep === 'brand'}
					<!-- Brand Configuration Step -->
					<div>
						<div class="mb-6">
							<h2 class="text-xl font-bold text-text">Email Sender Settings</h2>
							<p class="text-text-muted mt-1">How do you want your emails to appear?</p>
						</div>

						{#if error}
							<Alert type="error" title="Error" class="mb-6">
								<p>{error}</p>
							</Alert>
						{/if}

						<form class="space-y-5" onsubmit={(e) => { e.preventDefault(); nextStep(); }}>
							<div>
								<label for="from-name" class="form-label">From Name</label>
								<Input
									id="from-name"
									type="text"
									placeholder="Your Company Name"
									bind:value={fromName}
									required
								/>
								<p class="mt-1 text-xs text-text-muted">The name that appears as the sender</p>
							</div>

							<div>
								<label for="from-email" class="form-label">From Email</label>
								<Input
									id="from-email"
									type="email"
									placeholder="newsletter@example.com"
									bind:value={fromEmail}
									required
								/>
								<p class="mt-1 text-xs text-text-muted">Must be a verified email address</p>
							</div>

							<div>
								<label for="reply-to" class="form-label">Reply-To Email (optional)</label>
								<Input
									id="reply-to"
									type="email"
									placeholder="support@example.com"
									bind:value={replyTo}
								/>
								<p class="mt-1 text-xs text-text-muted">Where replies will be sent</p>
							</div>

							<div class="flex gap-3 pt-4">
								<Button type="secondary" onclick={prevStep} disabled={loading}>
									<ChevronLeft class="w-4 h-4 mr-2" />
									Back
								</Button>
								<Button htmlType="submit" type="primary" disabled={loading || !fromName || !fromEmail} class="flex-1 justify-center">
									{#if loading}
										<Loader2 class="w-4 h-4 mr-2 animate-spin" />
										Saving...
									{:else}
										Save & Complete
										<ChevronRight class="w-4 h-4 ml-2" />
									{/if}
								</Button>
							</div>
						</form>

						<button
							type="button"
							onclick={skipForNow}
							class="mt-4 text-sm text-text-muted hover:text-text transition-colors w-full text-center"
						>
							Skip for now
						</button>
					</div>
				{:else if currentStep === 'complete'}
					<!-- Complete Step -->
					<div class="text-center">
						<div class="mx-auto w-16 h-16 rounded-full bg-green-500/10 flex items-center justify-center mb-6">
							<Check class="w-8 h-8 text-green-500" />
						</div>
						<h2 class="text-2xl font-bold text-text mb-2">You're all set!</h2>
						<p class="text-text-muted mb-8">
							Your sender settings have been saved. You're ready to start sending emails from {org.name}.
						</p>

						<Button type="primary" size="lg" onclick={goToDashboard} class="w-full justify-center">
							Go to Dashboard
							<ChevronRight class="w-4 h-4 ml-2" />
						</Button>
					</div>
				{/if}
			</Card>
		</div>

		<!-- Help Sidebar -->
		<div class="lg:col-span-2">
			<div class="sticky top-8 space-y-6">
				{#if currentStep === 'welcome'}
					<!-- Welcome Help -->
					<Card>
						<div class="flex items-center gap-2 mb-4">
							<HelpCircle class="w-5 h-5 text-primary" />
							<h3 class="font-semibold text-text">Why set this up?</h3>
						</div>
						<ul class="space-y-3 text-sm text-text-muted">
							<li class="flex items-start gap-2">
								<span class="text-primary mt-0.5">1.</span>
								<span>Subscribers see a familiar name in their inbox</span>
							</li>
							<li class="flex items-start gap-2">
								<span class="text-primary mt-0.5">2.</span>
								<span>Improves email deliverability and trust</span>
							</li>
							<li class="flex items-start gap-2">
								<span class="text-primary mt-0.5">3.</span>
								<span>Replies go to the right place</span>
							</li>
						</ul>
					</Card>

					<Card>
						<div class="flex items-center gap-2 mb-4">
							<Lightbulb class="w-5 h-5 text-amber-500" />
							<h3 class="font-semibold text-text">Pro tip</h3>
						</div>
						<p class="text-sm text-text-muted">
							Use a recognizable sender name like your company name or newsletter title. Avoid generic names like "Newsletter" or "Updates".
						</p>
					</Card>
				{:else if currentStep === 'brand'}
					<!-- Brand Help -->
					<Card>
						<div class="flex items-center gap-2 mb-4">
							<HelpCircle class="w-5 h-5 text-primary" />
							<h3 class="font-semibold text-text">Sender Settings Help</h3>
						</div>
						<div class="space-y-4 text-sm">
							<div>
								<p class="font-medium text-text">From Name</p>
								<p class="text-text-muted">This appears as the sender in email clients. Use your brand name or a person's name for a personal touch.</p>
							</div>
							<div>
								<p class="font-medium text-text">From Email</p>
								<p class="text-text-muted">The email address that sends your campaigns. This should be verified with your email provider.</p>
							</div>
							<div>
								<p class="font-medium text-text">Reply-To</p>
								<p class="text-text-muted">Where subscriber replies are sent. Use a monitored inbox like support@ or hello@.</p>
							</div>
						</div>
					</Card>

					<Card>
						<div class="flex items-center gap-2 mb-4">
							<Lightbulb class="w-5 h-5 text-amber-500" />
							<h3 class="font-semibold text-text">Best practices</h3>
						</div>
						<ul class="space-y-2 text-sm text-text-muted">
							<li class="flex items-start gap-2">
								<Check class="w-4 h-4 text-green-500 mt-0.5 flex-shrink-0" />
								<span>Use a consistent from name across all emails</span>
							</li>
							<li class="flex items-start gap-2">
								<Check class="w-4 h-4 text-green-500 mt-0.5 flex-shrink-0" />
								<span>Use a domain you own (not gmail.com)</span>
							</li>
							<li class="flex items-start gap-2">
								<Check class="w-4 h-4 text-green-500 mt-0.5 flex-shrink-0" />
								<span>Monitor your reply-to inbox regularly</span>
							</li>
						</ul>
					</Card>

					<Card>
						<div class="flex items-center gap-2 mb-4">
							<Shield class="w-5 h-5 text-green-500" />
							<h3 class="font-semibold text-text">Email verification</h3>
						</div>
						<p class="text-sm text-text-muted">
							Make sure your from email is verified with your email provider (like Amazon SES) before sending campaigns.
						</p>
					</Card>
				{:else if currentStep === 'complete'}
					<!-- Complete Help -->
					<Card>
						<div class="flex items-center gap-2 mb-4">
							<Lightbulb class="w-5 h-5 text-amber-500" />
							<h3 class="font-semibold text-text">What's next?</h3>
						</div>
						<ul class="space-y-3 text-sm text-text-muted">
							<li class="flex items-start gap-2">
								<Check class="w-4 h-4 text-green-500 mt-0.5" />
								<span>Create your first email list</span>
							</li>
							<li class="flex items-start gap-2">
								<Check class="w-4 h-4 text-green-500 mt-0.5" />
								<span>Add subscribers manually or import a CSV</span>
							</li>
							<li class="flex items-start gap-2">
								<Check class="w-4 h-4 text-green-500 mt-0.5" />
								<span>Create your first campaign</span>
							</li>
							<li class="flex items-start gap-2">
								<Check class="w-4 h-4 text-green-500 mt-0.5" />
								<span>Set up welcome email automation</span>
							</li>
						</ul>
					</Card>
				{/if}
			</div>
		</div>
	</div>
{/if}
