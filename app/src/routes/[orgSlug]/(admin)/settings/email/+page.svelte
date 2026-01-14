<script lang="ts">
	import { browser } from '$app/environment';
	import { getOrganization, updateOrgEmailSettings, type OrgInfo } from '$lib/api';
	import { Button, Card, Input, Alert, LoadingSpinner, SaveButton } from '$lib/components/ui';
	import { Mail, Check, AlertCircle } from 'lucide-svelte';

	let loading = $state(true);
	let saving = $state(false);
	let saved = $state(false);
	let org = $state<OrgInfo | null>(null);
	let error = $state('');

	// Form state
	let fromName = $state('');
	let fromEmail = $state('');
	let replyTo = $state('');

	$effect(() => {
		loadData();
	});

	async function loadData() {
		loading = true;
		error = '';

		try {
			const orgId = browser ? localStorage.getItem('currentOrgId') : null;
			if (!orgId) {
				error = 'No organization selected';
				loading = false;
				return;
			}

			org = await getOrganization({}, orgId);
			// Populate form with existing values
			fromName = org.from_name || '';
			fromEmail = org.from_email || '';
			replyTo = org.reply_to || '';
		} catch (err: any) {
			console.error('Failed to load data:', err);
			error = err.message || 'Failed to load data';
		} finally {
			loading = false;
		}
	}

	async function saveEmailSettings() {
		if (!org) return;

		saving = true;
		saved = false;
		error = '';

		try {
			await updateOrgEmailSettings({}, {
				from_name: fromName || undefined,
				from_email: fromEmail || undefined,
				reply_to: replyTo || undefined
			}, org.id);
			saved = true;
			setTimeout(() => { saved = false; }, 2000);
			await loadData();
		} catch (err: any) {
			console.error('Failed to save email settings:', err);
			error = err.message || 'Failed to save email settings';
		} finally {
			saving = false;
		}
	}

	let hasChanges = $derived(
		org && (
			fromName !== (org.from_name || '') ||
			fromEmail !== (org.from_email || '') ||
			replyTo !== (org.reply_to || '')
		)
	);

	let isConfigured = $derived(org?.from_name || org?.from_email);
</script>

<svelte:head>
	<title>Email Settings - {org?.name || 'Loading'}</title>
</svelte:head>

{#if error}
	<Alert type="error" title="Error">
		<p>{error}</p>
	</Alert>
{/if}

{#if loading}
	<LoadingSpinner size="large" />
{:else}
	<div class="space-y-6">
		<Card>
			<div class="card-header">
				<h2 class="card-title">Sender Identity</h2>
				<p class="card-subtitle">Configure how your emails appear to recipients</p>
			</div>

			<div class="border border-border rounded-lg">
				<div class="p-4 flex items-center justify-between">
					<div class="flex items-center gap-4">
						<div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
							<Mail class="h-6 w-6 text-blue-600" />
						</div>
						<div>
							<h3 class="card-section-title">Email Sender</h3>
							<p class="card-section-subtitle">
								{#if isConfigured}
									{org?.from_name || 'No name'} &lt;{org?.from_email || 'No email'}&gt;
								{:else}
									Not configured - using platform defaults
								{/if}
							</p>
						</div>
					</div>
					<div class="flex items-center gap-2">
						{#if isConfigured}
							<div class="flex items-center gap-2 text-green-600">
								<Check class="h-5 w-5" />
								<span class="text-sm font-medium">Configured</span>
							</div>
						{:else}
							<div class="flex items-center gap-2 text-text-muted">
								<AlertCircle class="h-5 w-5" />
								<span class="text-sm">Using defaults</span>
							</div>
						{/if}
					</div>
				</div>

				<div class="border-t border-border p-4 bg-bg-secondary">
					<div class="space-y-4 max-w-lg">
						<div>
							<label for="from-name" class="form-label">From Name</label>
							<Input
								id="from-name"
								type="text"
								bind:value={fromName}
								placeholder="Acme Inc"
							/>
							<p class="mt-1 text-xs text-text-muted">The name shown in the recipient's inbox</p>
						</div>

						<div>
							<label for="from-email" class="form-label">From Email</label>
							<Input
								id="from-email"
								type="email"
								bind:value={fromEmail}
								placeholder="hello@yourdomain.com"
							/>
							<p class="mt-1 text-xs text-text-muted">Must be a verified email address with your email provider</p>
						</div>

						<div>
							<label for="reply-to" class="form-label">Reply-To Email</label>
							<Input
								id="reply-to"
								type="email"
								bind:value={replyTo}
								placeholder="support@yourdomain.com"
							/>
							<p class="mt-1 text-xs text-text-muted">Where replies will be sent (optional, defaults to From Email)</p>
						</div>

						<div class="flex justify-end gap-3 pt-2">
							<SaveButton
								label="Save Settings"
								{saving}
								{saved}
								onclick={saveEmailSettings}
								disabled={!hasChanges}
							/>
						</div>
					</div>
				</div>
			</div>
		</Card>

		<!-- Info Card -->
		<Card>
			<div class="flex items-start gap-4">
				<div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center flex-shrink-0">
					<svg class="h-5 w-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</div>
				<div>
					<h3 class="card-title">How it works</h3>
					<p class="card-subtitle">
						These settings are used as defaults for all emails sent from your site, including campaigns, sequences, and transactional emails.
						You can override these settings on a per-campaign or per-sequence basis.
					</p>
				</div>
			</div>
		</Card>
	</div>
{/if}
