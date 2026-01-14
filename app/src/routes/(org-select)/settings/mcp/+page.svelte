<script lang="ts">
	import { Card, Button, Input, Alert, LoadingSpinner, Badge, Modal } from '$lib/components/ui';
	import { Copy, Plus, Trash2, Key, ChevronDown, ChevronRight, CheckCircle } from 'lucide-svelte';
	import { listMCPAPIKeys, createMCPAPIKey, revokeMCPAPIKey, type MCPAPIKeyInfo } from '$lib/api';
	import { onMount } from 'svelte';

	let mcpEndpoint = 'https://outlet.sh/mcp';

	let loading = $state(true);
	let creating = $state(false);
	let error = $state('');
	let success = $state('');
	let copiedEndpoint = $state(false);
	let copiedKey = $state('');

	let apiKeys = $state<MCPAPIKeyInfo[]>([]);
	let showCreateModal = $state(false);
	let newKeyName = $state('');
	let newKeyValue = $state('');

	// Setup guide state
	let activeTab = $state<'claude-desktop' | 'chatgpt' | 'claude-code' | 'cursor'>('claude-desktop');

	// MCP Tools grouped by category
	const mcpTools = {
		Organization: [
			{ name: 'org_list', desc: 'List organizations you have access to' },
			{ name: 'org_select', desc: 'Select an organization to work with' },
			{ name: 'org_get', desc: 'Get organization settings and config' },
			{ name: 'org_update', desc: 'Update organization settings' },
			{ name: 'org_payment_setup', desc: 'Configure Stripe credentials' }
		],
		Products: [
			{ name: 'product_create', desc: 'Create a new product' },
			{ name: 'product_list', desc: 'List all products' },
			{ name: 'product_update', desc: 'Update product details' }
		],
		Billing: [
			{ name: 'price_create', desc: 'Create a price for a product' },
			{ name: 'price_list', desc: 'List all prices' },
			{ name: 'price_update', desc: 'Update price properties' },
			{ name: 'price_delete', desc: 'Delete a price' },
			{ name: 'payment_sync', desc: 'Sync to Stripe' }
		],
		Emails: [
			{ name: 'email_list_create', desc: 'Create an email list' },
			{ name: 'email_list_list', desc: 'List all email lists' },
			{ name: 'sequence_create', desc: 'Create an email sequence' },
			{ name: 'sequence_list', desc: 'List all sequences' },
			{ name: 'sequence_email_add', desc: 'Add email to sequence' },
			{ name: 'sequence_email_update', desc: 'Update sequence email' },
			{ name: 'sequence_get', desc: 'Get sequence with emails' }
		],
		'LLM Gateway': [
			{ name: 'llm_setup', desc: 'Configure LLM provider credentials' },
			{ name: 'llm_config', desc: 'View/update LLM configuration' },
			{ name: 'llm_pricing_get', desc: 'View model credit multipliers' },
			{ name: 'llm_pricing_update', desc: 'Set credit multiplier overrides' },
			{ name: 'llm_test', desc: 'Test LLM connection' },
			{ name: 'llm_stats', desc: 'Get LLM usage statistics' }
		]
	};

	let expandedCategories = $state<Record<string, boolean>>({
		Organization: true,
		Products: false,
		Billing: false,
		Emails: false,
		'LLM Gateway': false
	});

	onMount(async () => {
		await loadKeys();
	});

	async function loadKeys() {
		loading = true;
		error = '';
		try {
			const response = await listMCPAPIKeys();
			apiKeys = response.keys || [];
		} catch (err: any) {
			// No keys yet is ok
			console.log('No MCP keys yet');
		} finally {
			loading = false;
		}
	}

	async function createKey() {
		if (!newKeyName.trim()) {
			error = 'Please enter a name for the API key';
			return;
		}

		creating = true;
		error = '';
		try {
			const response = await createMCPAPIKey({ name: newKeyName.trim() });
			newKeyValue = response.key || '';
			await loadKeys();
			success = "API key created successfully. Copy it now - you won't be able to see it again!";
		} catch (err: any) {
			error = err.message || 'Failed to create API key';
		} finally {
			creating = false;
		}
	}

	async function revokeKey(id: string) {
		if (!confirm('Are you sure you want to revoke this API key? This cannot be undone.')) {
			return;
		}

		try {
			await revokeMCPAPIKey({}, id);
			await loadKeys();
			success = 'API key revoked';
		} catch (err: any) {
			error = err.message || 'Failed to revoke API key';
		}
	}

	function copyToClipboard(text: string, type: 'endpoint' | 'key') {
		navigator.clipboard.writeText(text);
		if (type === 'endpoint') {
			copiedEndpoint = true;
			setTimeout(() => (copiedEndpoint = false), 2000);
		} else {
			copiedKey = text;
			setTimeout(() => (copiedKey = ''), 2000);
		}
	}

	function closeCreateModal() {
		showCreateModal = false;
		newKeyName = '';
		newKeyValue = '';
	}

	function formatDate(dateStr: string | null | undefined): string {
		if (!dateStr) return 'Never';
		return new Date(dateStr).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	function toggleCategory(category: string) {
		expandedCategories[category] = !expandedCategories[category];
	}

	// Calculate total tools
	let totalTools = $derived(Object.values(mcpTools).flat().length);
</script>

<svelte:head>
	<title>MCP Configuration - Outlet</title>
</svelte:head>

<div class="space-y-6">
	{#if error}
		<Alert type="error" title="Error" onclose={() => (error = '')}>
			<p>{error}</p>
		</Alert>
	{/if}

	{#if success}
		<Alert type="success" title="Success" onclose={() => (success = '')}>
			<p>{success}</p>
		</Alert>
	{/if}

	<!-- MCP Endpoint -->
	<Card>
		<div class="flex items-start justify-between gap-4">
			<div class="flex-1">
				<h2 class="text-lg font-medium text-text mb-1">MCP Endpoint</h2>
				<p class="text-sm text-text-muted">
					Connect AI assistants to help you configure and manage your platform
				</p>
			</div>
			<Badge type="success">Active</Badge>
		</div>

		<div class="mt-4">
			<label class="form-label">Your MCP Server URL</label>
			<div class="flex gap-2">
				<Input type="text" value={mcpEndpoint} readonly class="font-mono text-sm" />
				<Button type="secondary" onclick={() => copyToClipboard(mcpEndpoint, 'endpoint')}>
					{#if copiedEndpoint}
						<CheckCircle class="h-4 w-4 text-green-500" />
					{:else}
						<Copy class="h-4 w-4" />
					{/if}
				</Button>
			</div>
		</div>
	</Card>

	<!-- API Keys -->
	<Card>
		<div class="flex items-center justify-between">
			<div>
				<h2 class="text-lg font-medium text-text mb-1">API Keys</h2>
				<p class="text-sm text-text-muted">Manage API keys for MCP authentication</p>
			</div>
			<Button type="primary" size="sm" onclick={() => (showCreateModal = true)}>
				<Plus class="mr-1.5 h-4 w-4" />
				Create Key
			</Button>
		</div>

		{#if loading}
			<div class="flex justify-center py-8">
				<LoadingSpinner />
			</div>
		{:else if apiKeys.length === 0}
			<div class="mt-6 text-center py-8 text-text-muted">
				<Key class="mx-auto h-12 w-12 opacity-50" />
				<p class="mt-2">No API keys yet</p>
				<p class="text-sm">Create an API key to connect AI assistants</p>
			</div>
		{:else}
			<div class="mt-6 overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="border-b border-border">
							<th class="text-left py-2 font-medium text-text-muted">Name</th>
							<th class="text-left py-2 font-medium text-text-muted">Key</th>
							<th class="text-left py-2 font-medium text-text-muted">Last Used</th>
							<th class="text-left py-2 font-medium text-text-muted">Created</th>
							<th class="text-right py-2 font-medium text-text-muted">Actions</th>
						</tr>
					</thead>
					<tbody>
						{#each apiKeys as key (key.id)}
							<tr class="border-b border-border/50">
								<td class="py-3 font-medium">{key.name}</td>
								<td class="py-3 font-mono text-text-muted">{key.key_prefix}...</td>
								<td class="py-3 text-text-muted">{formatDate(key.last_used)}</td>
								<td class="py-3 text-text-muted">{formatDate(key.created_at)}</td>
								<td class="py-3 text-right">
									<Button type="danger" size="icon" onclick={() => revokeKey(key.id)}>
										<Trash2 class="h-4 w-4" />
									</Button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</Card>

	<!-- Setup Instructions -->
	<Card>
		<h2 class="text-lg font-medium text-text mb-1">Setup Instructions</h2>
		<p class="text-sm text-text-muted">Connect your favorite AI assistant to Outlet</p>

		<div class="mt-4 flex flex-wrap gap-2">
			<Button
				type={activeTab === 'claude-desktop' ? 'primary' : 'secondary'}
				size="sm"
				onclick={() => (activeTab = 'claude-desktop')}
			>
				Claude Desktop
			</Button>
			<Button
				type={activeTab === 'chatgpt' ? 'primary' : 'secondary'}
				size="sm"
				onclick={() => (activeTab = 'chatgpt')}
			>
				ChatGPT
			</Button>
			<Button
				type={activeTab === 'claude-code' ? 'primary' : 'secondary'}
				size="sm"
				onclick={() => (activeTab = 'claude-code')}
			>
				Claude Code
			</Button>
			<Button
				type={activeTab === 'cursor' ? 'primary' : 'secondary'}
				size="sm"
				onclick={() => (activeTab = 'cursor')}
			>
				Cursor / IDE
			</Button>
		</div>

		<div class="mt-4 rounded-lg bg-surface-secondary p-4">
			{#if activeTab === 'claude-desktop'}
				<div class="flex items-center gap-2 mb-3">
					<h4 class="font-medium text-text">Claude Desktop Setup</h4>
					<Badge type="success" size="sm">OAuth</Badge>
				</div>
				<p class="text-sm text-text-muted mb-3">
					Claude Desktop uses OAuth for secure authentication. Add the following to your <code
						class="bg-surface-tertiary px-1 rounded">claude_desktop_config.json</code
					>:
				</p>
				<pre class="bg-surface-tertiary p-3 rounded text-sm overflow-x-auto font-mono">{`{
  "mcpServers": {
    "outlet": {
      "url": "${mcpEndpoint}",
      "oauth": {
        "authorize_url": "https://outlet.sh/oauth/authorize",
        "token_url": "https://outlet.sh/oauth/token",
        "client_id": "claude-desktop",
        "scopes": ["mcp:full"]
      }
    }
  }
}`}</pre>
				<p class="mt-3 text-sm text-text-muted">
					When you connect, Claude Desktop will open a browser window for you to log in and
					authorize access.
				</p>
			{:else if activeTab === 'chatgpt'}
				<div class="flex items-center gap-2 mb-3">
					<h4 class="font-medium text-text">ChatGPT Setup</h4>
					<Badge type="success" size="sm">OAuth</Badge>
				</div>
				<p class="text-sm text-text-muted mb-3">
					ChatGPT uses OAuth for secure authentication when connecting to MCP servers:
				</p>
				<ol class="text-sm text-text-muted space-y-2 list-decimal list-inside">
					<li>Go to <strong>Configure</strong> &rarr; <strong>Actions</strong></li>
					<li>Click <strong>Create new action</strong></li>
					<li>
						Set the server URL to: <code class="bg-surface-tertiary px-1 rounded"
							>{mcpEndpoint}</code
						>
					</li>
					<li>Select <strong>OAuth</strong> for authentication</li>
					<li>
						Use authorize URL: <code class="bg-surface-tertiary px-1 rounded"
							>https://outlet.sh/oauth/authorize</code
						>
					</li>
					<li>
						Use token URL: <code class="bg-surface-tertiary px-1 rounded"
							>https://outlet.sh/oauth/token</code
						>
					</li>
				</ol>
				<p class="mt-3 text-sm text-text-muted">
					ChatGPT will prompt you to authorize access when first connecting.
				</p>
			{:else if activeTab === 'claude-code'}
				<div class="flex items-center gap-2 mb-3">
					<h4 class="font-medium text-text">Claude Code Setup</h4>
					<Badge type="neutral" size="sm">API Key</Badge>
				</div>
				<p class="text-sm text-text-muted mb-3">
					Claude Code (the CLI) uses API keys for authentication. Create an API key above, then add
					to your <code class="bg-surface-tertiary px-1 rounded">~/.claude/settings.json</code>:
				</p>
				<pre class="bg-surface-tertiary p-3 rounded text-sm overflow-x-auto font-mono">{`{
  "mcpServers": {
    "outlet": {
      "url": "${mcpEndpoint}",
      "headers": {
        "Authorization": "Bearer YOUR_API_KEY"
      }
    }
  }
}`}</pre>
				<p class="mt-3 text-sm text-text-muted">
					Replace <code class="bg-surface-tertiary px-1 rounded">YOUR_API_KEY</code> with an API key created
					above.
				</p>
			{:else if activeTab === 'cursor'}
				<div class="flex items-center gap-2 mb-3">
					<h4 class="font-medium text-text">Cursor / Windsurf / VS Code Setup</h4>
					<Badge type="neutral" size="sm">API Key</Badge>
				</div>
				<p class="text-sm text-text-muted mb-3">
					IDE-based tools use API keys for authentication. Create an API key above, then add to your
					IDE config:
				</p>
				<pre class="bg-surface-tertiary p-3 rounded text-sm overflow-x-auto font-mono">{`{
  "mcp.servers": [
    {
      "name": "Outlet",
      "url": "${mcpEndpoint}",
      "auth": {
        "type": "bearer",
        "token": "YOUR_API_KEY"
      }
    }
  ]
}`}</pre>
				<p class="mt-3 text-sm text-text-muted">
					Replace <code class="bg-surface-tertiary px-1 rounded">YOUR_API_KEY</code> with an API key created
					above. Configuration file location varies by IDE.
				</p>
			{/if}
		</div>
	</Card>

	<!-- Available Tools -->
	<Card>
		<h2 class="text-lg font-medium text-text mb-1">Available Tools ({totalTools} tools)</h2>
		<p class="text-sm text-text-muted">
			MCP tools that AI assistants can use to help you manage your platform
		</p>

		<div class="mt-4 space-y-2">
			{#each Object.entries(mcpTools) as [category, tools]}
				<div class="border border-border rounded-lg">
					<button
						type="button"
						class="w-full flex items-center justify-between p-3 hover:bg-surface-secondary transition-colors"
						onclick={() => toggleCategory(category)}
					>
						<span class="font-medium text-text">{category}</span>
						<div class="flex items-center gap-2">
							<Badge type="neutral">{tools.length}</Badge>
							{#if expandedCategories[category]}
								<ChevronDown class="h-4 w-4 text-text-muted" />
							{:else}
								<ChevronRight class="h-4 w-4 text-text-muted" />
							{/if}
						</div>
					</button>
					{#if expandedCategories[category]}
						<div class="border-t border-border p-3 space-y-2">
							{#each tools as tool}
								<div class="flex items-start gap-3 text-sm">
									<code
										class="bg-surface-tertiary px-1.5 py-0.5 rounded text-xs font-mono text-primary"
										>{tool.name}</code
									>
									<span class="text-text-muted">{tool.desc}</span>
								</div>
							{/each}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	</Card>
</div>

<!-- Create Key Modal -->
<Modal bind:show={showCreateModal} title="Create API Key" onclose={closeCreateModal}>
	{#if newKeyValue}
		<div class="space-y-4">
			<Alert type="warning" title="Save your API key">
				<p>This is the only time you'll see this key. Copy it now and store it securely.</p>
			</Alert>
			<div>
				<label class="form-label">API Key</label>
				<div class="flex gap-2">
					<Input type="text" value={newKeyValue} readonly class="font-mono text-sm" />
					<Button type="secondary" onclick={() => copyToClipboard(newKeyValue, 'key')}>
						{#if copiedKey === newKeyValue}
							<CheckCircle class="h-4 w-4 text-green-500" />
						{:else}
							<Copy class="h-4 w-4" />
						{/if}
					</Button>
				</div>
			</div>
		</div>

		{#snippet footer()}
			<div class="flex justify-end">
				<Button type="primary" onclick={closeCreateModal}>Done</Button>
			</div>
		{/snippet}
	{:else}
		<div class="space-y-4">
			<div>
				<label for="key-name" class="form-label">Key Name</label>
				<Input
					type="text"
					id="key-name"
					bind:value={newKeyName}
					placeholder="e.g., Claude Code, My ChatGPT"
				/>
				<p class="mt-1 text-xs text-text-muted">A descriptive name to identify this key</p>
			</div>
		</div>

		{#snippet footer()}
			<div class="flex justify-end gap-3">
				<Button type="secondary" onclick={closeCreateModal}>Cancel</Button>
				<Button type="primary" onclick={createKey} disabled={creating || !newKeyName.trim()}>
					{creating ? 'Creating...' : 'Create Key'}
				</Button>
			</div>
		{/snippet}
	{/if}
</Modal>
