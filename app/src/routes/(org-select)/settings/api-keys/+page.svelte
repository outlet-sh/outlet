<script lang="ts">
	import { listMCPAPIKeys, createMCPAPIKey, revokeMCPAPIKey, type MCPAPIKeyInfo } from '$lib/api';
	import {
		Button,
		Card,
		Input,
		Badge,
		Modal,
		AlertDialog,
		Alert,
		LoadingSpinner
	} from '$lib/components/ui';
	import { Plus, Key, Copy, Check } from 'lucide-svelte';

	let loading = $state(true);
	let keys = $state<MCPAPIKeyInfo[]>([]);
	let error = $state('');

	// Add key modal
	let showAddModal = $state(false);
	let adding = $state(false);
	let newName = $state('');
	let newExpiresAt = $state('');

	// New key display modal
	let showNewKeyModal = $state(false);
	let newKeyValue = $state('');
	let copied = $state(false);

	// Revoke confirm
	let showRevokeConfirm = $state(false);
	let revokingKey = $state<MCPAPIKeyInfo | null>(null);
	let revoking = $state(false);

	$effect(() => {
		loadKeys();
	});

	async function loadKeys() {
		loading = true;
		error = '';

		try {
			const response = await listMCPAPIKeys();
			keys = response.keys || [];
		} catch (err) {
			console.error('Failed to load API keys:', err);
			error = 'Failed to load API keys';
		} finally {
			loading = false;
		}
	}

	function openAddModal() {
		newName = '';
		newExpiresAt = '';
		showAddModal = true;
	}

	function closeAddModal() {
		showAddModal = false;
	}

	async function submitAdd() {
		if (!newName) return;

		adding = true;
		error = '';

		try {
			const response = await createMCPAPIKey({
				name: newName,
				expires_at: newExpiresAt || undefined
			});
			closeAddModal();
			// Show the new key
			newKeyValue = response.key || '';
			showNewKeyModal = true;
			await loadKeys();
		} catch (err: any) {
			console.error('Failed to create API key:', err);
			error = err.message || 'Failed to create API key';
		} finally {
			adding = false;
		}
	}

	function closeNewKeyModal() {
		showNewKeyModal = false;
		newKeyValue = '';
		copied = false;
	}

	async function copyKey() {
		if (newKeyValue) {
			await navigator.clipboard.writeText(newKeyValue);
			copied = true;
			setTimeout(() => {
				copied = false;
			}, 2000);
		}
	}

	function confirmRevoke(key: MCPAPIKeyInfo) {
		revokingKey = key;
		showRevokeConfirm = true;
	}

	async function executeRevoke() {
		if (!revokingKey) return;

		revoking = true;
		try {
			await revokeMCPAPIKey({}, revokingKey.id);
			showRevokeConfirm = false;
			revokingKey = null;
			await loadKeys();
		} catch (err: any) {
			console.error('Failed to revoke API key:', err);
			error = err.message || 'Failed to revoke API key';
		} finally {
			revoking = false;
		}
	}

	function formatDate(dateString: string): string {
		if (!dateString) return 'N/A';
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	function isExpired(expiresAt?: string): boolean {
		if (!expiresAt) return false;
		return new Date(expiresAt) < new Date();
	}
</script>

<svelte:head>
	<title>API Keys - Outlet</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<div>
			<p class="text-sm text-text-secondary">
				Manage API keys for MCP (Claude, ChatGPT) integrations
			</p>
			<p class="text-xs text-text-muted mt-1">
				API keys allow AI assistants to manage your SaaS through the Model Context Protocol
			</p>
		</div>
		<Button type="primary" onclick={openAddModal}>
			<Plus class="mr-2 h-4 w-4" />
			Create API Key
		</Button>
	</div>

	{#if error}
		<Alert type="error" title="Error">
			<p>{error}</p>
		</Alert>
	{/if}

	{#if loading}
		<LoadingSpinner size="large" />
	{:else if keys.length === 0}
		<Card>
			<div class="text-center py-12">
				<Key class="mx-auto h-12 w-12 text-text-muted" />
				<h3 class="mt-2 text-sm font-medium text-text">No API keys</h3>
				<p class="mt-1 text-sm text-text-muted">
					Create an API key to connect AI assistants via MCP.
				</p>
				<div class="mt-6">
					<Button type="primary" onclick={openAddModal}>
						<Plus class="mr-2 h-4 w-4" />
						Create API Key
					</Button>
				</div>
			</div>
		</Card>
	{:else}
		<Card>
			<div class="divide-y divide-border">
				{#each keys as key}
					<div class="flex items-center justify-between py-4 first:pt-0 last:pb-0">
						<div class="flex items-center gap-3">
							<div
								class="h-10 w-10 rounded-full bg-surface-elevated flex items-center justify-center"
							>
								<Key class="h-5 w-5 text-text-muted" />
							</div>
							<div>
								<div class="flex items-center gap-2">
									<span class="font-medium text-text">{key.name}</span>
									<code class="text-xs bg-surface-elevated px-2 py-0.5 rounded text-text-muted"
										>{key.key_prefix}...</code
									>
									{#if key.expires_at && isExpired(key.expires_at)}
										<Badge variant="error" size="sm">Expired</Badge>
									{:else if key.expires_at}
										<Badge variant="warning" size="sm">Expires {formatDate(key.expires_at)}</Badge>
									{:else}
										<Badge variant="success" size="sm">Active</Badge>
									{/if}
								</div>
								<div class="flex items-center gap-4 text-sm text-text-muted">
									<span>Created {formatDate(key.created_at)}</span>
									{#if key.last_used}
										<span>Last used {formatDate(key.last_used)}</span>
									{:else}
										<span>Never used</span>
									{/if}
								</div>
							</div>
						</div>
						<div class="flex items-center gap-2">
							<Button type="danger" size="sm" onclick={() => confirmRevoke(key)}>Revoke</Button>
						</div>
					</div>
				{/each}
			</div>
		</Card>
	{/if}

	<Card>
		<h3 class="text-sm font-medium text-text mb-3">Usage Instructions</h3>
		<div class="text-sm text-text-muted space-y-2">
			<p>
				Use your API key to connect AI assistants to Outlet via the Model Context Protocol (MCP).
			</p>
			<p class="font-medium text-text mt-4">For Claude Desktop:</p>
			<pre class="bg-surface-elevated p-3 rounded text-xs overflow-x-auto mt-2"><code
					>{`{
  "mcpServers": {
    "outlet": {
      "transport": {
        "type": "sse",
        "url": "https://outlet.sh/mcp/sse",
        "headers": {
          "Authorization": "Bearer lv_your_api_key"
        }
      }
    }
  }
}`}</code
				></pre>
		</div>
	</Card>
</div>

<!-- Create API Key Modal -->
<Modal bind:show={showAddModal} title="Create API Key">
	<div class="space-y-4">
		<div>
			<label for="new-name" class="form-label">Name</label>
			<Input id="new-name" type="text" bind:value={newName} placeholder="My API Key" />
			<p class="text-xs text-text-muted mt-1">A name to identify this key</p>
		</div>
		<div>
			<label for="new-expires" class="form-label">Expiration (optional)</label>
			<Input id="new-expires" type="date" bind:value={newExpiresAt} />
			<p class="text-xs text-text-muted mt-1">Leave empty for no expiration</p>
		</div>
	</div>

	{#snippet footer()}
		<div class="flex justify-end gap-3">
			<Button type="secondary" onclick={closeAddModal} disabled={adding}>Cancel</Button>
			<Button type="primary" onclick={submitAdd} disabled={!newName || adding}>
				{adding ? 'Creating...' : 'Create Key'}
			</Button>
		</div>
	{/snippet}
</Modal>

<!-- New Key Display Modal -->
<Modal bind:show={showNewKeyModal} title="API Key Created">
	<div class="space-y-4">
		<Alert type="warning" title="Save this key now">
			<p>
				This is the only time you'll see the full API key. Copy it now - you won't be able to see it
				again.
			</p>
		</Alert>

		<div>
			<label class="form-label">Your API Key</label>
			<div class="flex gap-2">
				<code class="flex-1 bg-surface-elevated p-3 rounded text-sm font-mono break-all">
					{newKeyValue}
				</code>
				<Button type="secondary" onclick={copyKey}>
					{#if copied}
						<Check class="h-4 w-4 text-success" />
					{:else}
						<Copy class="h-4 w-4" />
					{/if}
				</Button>
			</div>
		</div>
	</div>

	{#snippet footer()}
		<div class="flex justify-end">
			<Button type="primary" onclick={closeNewKeyModal}>Done</Button>
		</div>
	{/snippet}
</Modal>

<!-- Revoke Confirmation -->
<AlertDialog
	bind:open={showRevokeConfirm}
	title="Revoke API Key"
	description={revokingKey
		? `Are you sure you want to revoke "${revokingKey.name}"? Any applications using this key will lose access immediately.`
		: ''}
	actionLabel={revoking ? 'Revoking...' : 'Revoke Key'}
	actionType="danger"
	onAction={executeRevoke}
	onCancel={() => {
		showRevokeConfirm = false;
		revokingKey = null;
	}}
/>
