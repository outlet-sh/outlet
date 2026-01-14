<script lang="ts">
	import { goto } from '$app/navigation';
	import * as api from '$lib/api';
	import type { OrgInfo } from '$lib/api';
	import { Button, Card, Input, Badge, Modal, Alert, LoadingSpinner } from '$lib/components/ui';
	import { Plus, Building2, ChevronRight } from 'lucide-svelte';

	let loading = $state(true);
	let organizations = $state<OrgInfo[]>([]);
	let error = $state('');

	async function selectOrg(org: OrgInfo) {
		localStorage.setItem('currentOrgId', org.id);
		localStorage.setItem('currentOrgName', org.name);
		localStorage.setItem('currentOrgSlug', org.slug);

		// Check if org has email configured, unless user has skipped setup
		const skipped = localStorage.getItem(`orgSetupSkipped_${org.slug}`);
		if (!skipped && !org.from_email) {
			goto(`/${org.slug}/welcome`);
			return;
		}

		goto(`/${org.slug}`);
	}

	// Create modal
	let showCreateModal = $state(false);
	let creating = $state(false);
	let newOrgName = $state('');
	let newOrgSlug = $state('');
	let newAppUrl = $state('');

	$effect(() => {
		checkPlatformSetup();
	});

	async function checkPlatformSetup() {
		// Check if platform is configured (unless user has skipped)
		const skipped = localStorage.getItem('setupSkipped');
		if (!skipped) {
			try {
				const status = await api.getSetupStatus();
				if (!status.platform_configured) {
					goto('/setup');
					return;
				}
			} catch (err) {
				console.error('Failed to check setup status:', err);
				// Continue to load orgs even if check fails
			}
		}

		loadOrganizations();
	}

	async function loadOrganizations() {
		loading = true;
		error = '';

		try {
			const response = await api.listOrganizations();
			organizations = response.organizations || [];
		} catch (err) {
			console.error('Failed to load organizations:', err);
			error = 'Failed to load organizations';
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		newOrgName = '';
		newOrgSlug = '';
		newAppUrl = '';
		showCreateModal = true;
	}

	function closeCreateModal() {
		showCreateModal = false;
	}

	function generateSlug(name: string): string {
		return name
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, '-')
			.replace(/^-|-$/g, '');
	}

	function handleNameInput() {
		newOrgSlug = generateSlug(newOrgName);
	}

	async function submitCreate() {
		if (!newOrgName || !newOrgSlug) return;

		creating = true;
		try {
			await api.createOrganization({
				name: newOrgName,
				slug: newOrgSlug,
				app_url: newAppUrl
			});
			closeCreateModal();
			await loadOrganizations();
		} catch (err: any) {
			console.error('Failed to create organization:', err);
			error = err.message || 'Failed to create organization';
		} finally {
			creating = false;
		}
	}
</script>

<svelte:head>
	<title>Sites - Outlet</title>
</svelte:head>

<div class="py-6">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="mb-6 flex justify-between items-center">
			<div>
				<h1 class="text-2xl font-semibold text-text">Sites</h1>
				<p class="mt-1 text-sm text-text-muted">Manage your client sites</p>
			</div>
			<Button type="primary" onclick={openCreateModal}>
				<Plus class="mr-2 h-4 w-4" />
				New Site
			</Button>
		</div>

		{#if error}
			<div class="mb-6">
				<Alert type="error" title="Error">
					<p>{error}</p>
				</Alert>
			</div>
		{/if}

		{#if loading}
			<LoadingSpinner size="large" />
		{:else if organizations.length === 0}
			<Card>
				<div class="text-center py-12">
					<Building2 class="mx-auto h-12 w-12 text-text-muted" />
					<h3 class="mt-2 text-sm font-medium text-text">No sites yet</h3>
					<p class="mt-1 text-sm text-text-muted">Get started by adding your first client site.</p>
					<div class="mt-6">
						<Button type="primary" onclick={openCreateModal}>
							<Plus class="mr-2 h-4 w-4" />
							New Site
						</Button>
					</div>
				</div>
			</Card>
		{:else}
			<div class="space-y-3">
				{#each organizations as org}
					<Button
						type="ghost"
						onclick={() => selectOrg(org)}
						class="w-full text-left bg-bg rounded-lg border border-border hover:border-primary hover:shadow-md transition-all p-4 group"
					>
						<div class="flex items-center justify-between w-full">
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
									<Building2 class="w-5 h-5 text-primary" />
								</div>
								<div>
									<div class="flex items-center gap-2">
										<span class="font-semibold text-text group-hover:text-primary transition-colors">{org.name}</span>
										<Badge variant="default" size="sm">{org.slug}</Badge>
									</div>
								</div>
							</div>
							<ChevronRight class="w-5 h-5 text-text-muted group-hover:text-primary transition-colors" />
						</div>
					</Button>
				{/each}
			</div>
		{/if}
	</div>
</div>

<!-- Create Modal -->
<Modal bind:show={showCreateModal} title="Add Site">
	<div class="space-y-4">
		<div>
			<label for="org-name" class="form-label">Site Name</label>
			<Input
				id="org-name"
				type="text"
				bind:value={newOrgName}
				oninput={handleNameInput}
				placeholder="Acme Corp"
			/>
		</div>
		<div>
			<label for="org-slug" class="form-label">Slug</label>
			<Input
				id="org-slug"
				type="text"
				bind:value={newOrgSlug}
				placeholder="acme-corp"
			/>
			<p class="mt-1 text-xs text-text-muted">Used as a unique identifier for API calls</p>
		</div>
		<div>
			<label for="app-url" class="form-label">App URL</label>
			<Input
				id="app-url"
				type="text"
				bind:value={newAppUrl}
				placeholder="https://yourapp.com"
			/>
			<p class="mt-1 text-xs text-text-muted">Base URL of your application for confirmation links</p>
		</div>
	</div>

	{#snippet footer()}
		<div class="flex justify-end gap-3">
			<Button type="secondary" onclick={closeCreateModal} disabled={creating}>
				Cancel
			</Button>
			<Button type="primary" onclick={submitCreate} disabled={!newOrgName || !newOrgSlug || creating}>
				{creating ? 'Creating...' : 'Create'}
			</Button>
		</div>
	{/snippet}
</Modal>
