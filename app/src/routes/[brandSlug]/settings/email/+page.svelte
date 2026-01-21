<script lang="ts">
	import { browser } from '$app/environment';
	import { onDestroy } from 'svelte';
	import {
		getOrganization,
		updateOrgEmailSettings,
		listDomainIdentities,
		createDomainIdentity,
		refreshDomainIdentity,
		type OrgInfo,
		type DomainIdentityInfo,
		type DNSRecord
	} from '$lib/api';
	import { getWebSocketClient } from '$lib/websocket/client';
	import { Button, Card, Input, Alert, LoadingSpinner, SaveButton, Badge } from '$lib/components/ui';
	import { Mail, Check, AlertCircle, RefreshCw, Copy, Shield, ExternalLink } from 'lucide-svelte';

	let loading = $state(true);
	let saving = $state(false);
	let saved = $state(false);
	let org = $state<OrgInfo | null>(null);
	let error = $state('');
	let domainIdentities = $state<DomainIdentityInfo[]>([]);
	let refreshingDomain = $state(false);
	let creatingDomain = $state(false);
	let copiedRecord = $state<string | null>(null);

	// Form state
	let fromName = $state('');
	let fromEmail = $state('');
	let replyTo = $state('');

	// WebSocket for real-time updates
	let wsUnsubscribe: (() => void) | null = null;

	$effect(() => {
		loadData();
		setupWebSocket();

		return () => {
			// Cleanup WebSocket subscription on destroy
			if (wsUnsubscribe) {
				wsUnsubscribe();
			}
		};
	});

	function setupWebSocket() {
		if (!browser) return;

		const orgId = localStorage.getItem('currentOrgId');
		if (!orgId) return;

		const ws = getWebSocketClient();

		// Connect if not already connected
		if (!ws.isConnected()) {
			ws.connect();
		}

		// Subscribe to org updates
		ws.send('subscribe', { org_id: orgId });

		// Listen for domain identity updates
		wsUnsubscribe = ws.on('domain_identity_update', (data: any) => {
			console.log('[WebSocket] Domain identity update received:', data);

			// Update the domain identity in our list
			domainIdentities = domainIdentities.map((identity) => {
				if (identity.id === data.id) {
					return {
						...identity,
						verification_status: data.verification_status,
						dkim_status: data.dkim_status,
						mail_from_status: data.mail_from_status,
						last_checked_at: data.last_checked_at
					};
				}
				return identity;
			});
		});
	}

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

			// Load domain identities
			try {
				const identitiesResponse = await listDomainIdentities({}, orgId);
				domainIdentities = identitiesResponse.identities || [];
			} catch (err) {
				console.error('Failed to load domain identities:', err);
				domainIdentities = [];
			}
		} catch (err: any) {
			console.error('Failed to load data:', err);
			error = err.message || 'Failed to load data';
		} finally {
			loading = false;
		}
	}

	async function handleCreateDomainIdentity() {
		if (!org) return;

		creatingDomain = true;
		error = '';

		try {
			const identity = await createDomainIdentity({}, { domain: '' }, org.id);
			domainIdentities = [...domainIdentities, identity];
		} catch (err: any) {
			console.error('Failed to create domain identity:', err);
			error = err.message || 'Failed to create domain identity';
		} finally {
			creatingDomain = false;
		}
	}

	async function handleRefreshDomainIdentity(identityId: string) {
		if (!org) return;

		refreshingDomain = true;
		error = '';

		try {
			const updated = await refreshDomainIdentity({}, org.id, identityId);
			domainIdentities = domainIdentities.map((i) => (i.id === identityId ? updated : i));
		} catch (err: any) {
			console.error('Failed to refresh domain identity:', err);
			error = err.message || 'Failed to refresh domain identity';
		} finally {
			refreshingDomain = false;
		}
	}

	function copyToClipboard(text: string, recordKey: string) {
		navigator.clipboard.writeText(text);
		copiedRecord = recordKey;
		setTimeout(() => {
			copiedRecord = null;
		}, 2000);
	}

	function getStatusBadge(status: string) {
		switch (status) {
			case 'Success':
			case 'success':
				return { variant: 'success' as const, label: 'Verified' };
			case 'Pending':
			case 'pending':
				return { variant: 'warning' as const, label: 'Pending' };
			case 'Failed':
			case 'failed':
			case 'temporary_failure':
				return { variant: 'error' as const, label: 'Failed' };
			default:
				return { variant: 'default' as const, label: status || 'Unknown' };
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

		<!-- Domain Verification Card -->
		<Card>
			<div class="card-header">
				<h2 class="card-title">Domain Verification</h2>
				<p class="card-subtitle">Verify your domain to improve email deliverability and enable DKIM signing</p>
			</div>

			{#if domainIdentities.length === 0}
				<div class="border border-border rounded-lg p-6 text-center">
					<div class="w-12 h-12 mx-auto bg-base-200 rounded-lg flex items-center justify-center mb-4">
						<Shield class="h-6 w-6 text-base-content/50" />
					</div>
					<h3 class="text-base font-medium mb-2">No domains configured</h3>
					<p class="text-sm text-base-content/60 mb-4">
						{#if !fromEmail}
							Set your From Email above first, then verify your domain.
						{:else}
							Verify your domain to authenticate emails and improve deliverability.
						{/if}
					</p>
					{#if fromEmail}
						<Button
							variant="primary"
							onclick={handleCreateDomainIdentity}
							disabled={creatingDomain}
						>
							{#if creatingDomain}
								<LoadingSpinner size="normal" />
								<span>Setting up...</span>
							{:else}
								<Shield class="h-4 w-4" />
								<span>Verify Domain</span>
							{/if}
						</Button>
					{/if}
				</div>
			{:else}
				{#each domainIdentities as identity (identity.id)}
					{@const verificationBadge = getStatusBadge(identity.verification_status)}
					{@const dkimBadge = getStatusBadge(identity.dkim_status)}
					<div class="border border-border rounded-lg overflow-hidden mb-4 last:mb-0">
						<!-- Domain Header -->
						<div class="p-4 flex items-center justify-between bg-base-200/50">
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 bg-primary/10 rounded-lg flex items-center justify-center">
									<Shield class="h-5 w-5 text-primary" />
								</div>
								<div>
									<h3 class="font-medium">{identity.domain}</h3>
									<div class="flex items-center gap-2 mt-1">
										<span class="text-xs text-base-content/60">Verification:</span>
										<Badge variant={verificationBadge.variant} size="sm">{verificationBadge.label}</Badge>
										<span class="text-xs text-base-content/60 ml-2">DKIM:</span>
										<Badge variant={dkimBadge.variant} size="sm">{dkimBadge.label}</Badge>
									</div>
								</div>
							</div>
							<Button
								variant="ghost"
								size="sm"
								onclick={() => handleRefreshDomainIdentity(identity.id)}
								disabled={refreshingDomain}
							>
								<RefreshCw class="h-4 w-4 {refreshingDomain ? 'animate-spin' : ''}" />
								<span>Refresh</span>
							</Button>
						</div>

						<!-- DNS Records -->
						{#if identity.dns_records && identity.dns_records.length > 0}
							<div class="p-4">
								<div class="flex items-center gap-2 mb-3">
									<AlertCircle class="h-4 w-4 text-warning" />
									<span class="text-sm font-medium">Add these DNS records to your domain</span>
								</div>

								<div class="overflow-x-auto">
									<table class="table table-sm w-full">
										<thead>
											<tr>
												<th class="text-xs font-medium text-base-content/60 uppercase">Type</th>
												<th class="text-xs font-medium text-base-content/60 uppercase">Name/Host</th>
												<th class="text-xs font-medium text-base-content/60 uppercase">Value</th>
												<th class="text-xs font-medium text-base-content/60 uppercase">Purpose</th>
												<th class="w-10"></th>
											</tr>
										</thead>
										<tbody>
											{#each identity.dns_records as record, i}
												{@const recordKey = `${identity.id}-${i}`}
												<tr>
													<td>
														<Badge variant="default" size="sm">{record.type}</Badge>
													</td>
													<td>
														<code class="text-xs bg-base-200 px-2 py-1 rounded break-all">{record.name}</code>
													</td>
													<td>
														<code class="text-xs bg-base-200 px-2 py-1 rounded break-all max-w-xs block overflow-hidden text-ellipsis">{record.value}</code>
													</td>
													<td>
														<span class="text-xs text-base-content/60 capitalize">{record.purpose.replace(/_/g, ' ')}</span>
													</td>
													<td>
														<button
															class="btn btn-ghost btn-xs"
															onclick={() => copyToClipboard(record.value, recordKey)}
															title="Copy value"
														>
															{#if copiedRecord === recordKey}
																<Check class="h-3 w-3 text-success" />
															{:else}
																<Copy class="h-3 w-3" />
															{/if}
														</button>
													</td>
												</tr>
											{/each}
										</tbody>
									</table>
								</div>

								<div class="mt-4 p-3 bg-info/10 rounded-lg">
									<p class="text-xs text-info">
										<strong>Note:</strong> DNS changes can take up to 72 hours to propagate. Click "Refresh" to check verification status.
									</p>
								</div>
							</div>
						{/if}
					</div>
				{/each}
			{/if}
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
						These settings are used as defaults for all emails sent from this brand, including campaigns, sequences, and transactional emails.
						You can override these settings on a per-campaign or per-sequence basis.
					</p>
				</div>
			</div>
		</Card>
	</div>
{/if}
