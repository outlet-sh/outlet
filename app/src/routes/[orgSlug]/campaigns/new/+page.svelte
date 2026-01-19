<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import {
		listLists,
		listEmailDesigns,
		createCampaign,
		type ListInfo,
		type EmailDesignInfo,
		type CreateCampaignRequest
	} from '$lib/api';
	import {
		Button,
		Input,
		Checkbox,
		Alert,
		LoadingSpinner,
		HtmlEditor,
		Textarea,
		Card,
		Select
	} from '$lib/components/ui';
	import { ArrowLeft, Send, Save, Eye, FileText } from 'lucide-svelte';

	// Form state
	let name = $state('');
	let subject = $state('');
	let previewText = $state('');
	let fromName = $state('');
	let fromEmail = $state('');
	let replyTo = $state('');
	let htmlBody = $state('');
	let plainText = $state('');
	let trackOpens = $state(true);
	let trackClicks = $state(true);
	let selectedListIds = $state<string[]>([]);
	let excludeListIds = $state<string[]>([]);
	let selectedDesignId = $state('');

	// Data
	let lists = $state<ListInfo[]>([]);
	let designs = $state<EmailDesignInfo[]>([]);
	let loading = $state(true);
	let saving = $state(false);
	let error = $state('');

	// Derived
	let basePath = $derived(`/${$page.params.orgSlug}`);
	let totalSubscribers = $derived(() => {
		return lists
			.filter((l) => selectedListIds.includes(l.id))
			.reduce((sum, l) => sum + l.subscriber_count, 0);
	});

	// Form validation
	let isValid = $derived(() => {
		return (
			name.trim() !== '' &&
			subject.trim() !== '' &&
			htmlBody.trim() !== '' &&
			selectedListIds.length > 0
		);
	});

	$effect(() => {
		loadData();
	});

	async function loadData() {
		loading = true;
		error = '';

		try {
			const [listsRes, designsRes] = await Promise.all([
				listLists(),
				listEmailDesigns({})
			]);
			lists = listsRes.lists || [];
			designs = (designsRes.designs || []).filter((d) => d.is_active);
		} catch (err) {
			console.error('Failed to load data:', err);
			error = 'Failed to load lists and templates';
		} finally {
			loading = false;
		}
	}

	function handleDesignSelect(designId: string) {
		selectedDesignId = designId;
		if (designId) {
			const design = designs.find((d) => d.id === designId);
			if (design) {
				htmlBody = design.html_body;
				if (design.plain_text) {
					plainText = design.plain_text;
				}
			}
		}
	}

	function toggleList(listId: string) {
		if (selectedListIds.includes(listId)) {
			selectedListIds = selectedListIds.filter((id) => id !== listId);
		} else {
			selectedListIds = [...selectedListIds, listId];
		}
	}

	function toggleExcludeList(listId: string) {
		if (excludeListIds.includes(listId)) {
			excludeListIds = excludeListIds.filter((id) => id !== listId);
		} else {
			excludeListIds = [...excludeListIds, listId];
		}
	}

	async function handleSave(sendNow = false) {
		if (!isValid()) {
			error = 'Please fill in all required fields';
			return;
		}

		saving = true;
		error = '';

		try {
			const request: CreateCampaignRequest = {
				name: name.trim(),
				subject: subject.trim(),
				preview_text: previewText.trim() || undefined,
				from_name: fromName.trim() || undefined,
				from_email: fromEmail.trim() || undefined,
				reply_to: replyTo.trim() || undefined,
				html_body: htmlBody,
				plain_text: plainText.trim() || undefined,
				list_ids: selectedListIds,
				exclude_list_ids: excludeListIds.length > 0 ? excludeListIds : undefined,
				track_opens: trackOpens,
				track_clicks: trackClicks,
				design_id: selectedDesignId || undefined
			};

			const campaign = await createCampaign(request);
			goto(`${basePath}/campaigns/${campaign.id}`);
		} catch (err: any) {
			console.error('Failed to save campaign:', err);
			error = err.message || 'Failed to save campaign';
		} finally {
			saving = false;
		}
	}

	function generatePlainText() {
		// Simple HTML to plain text conversion
		const div = document.createElement('div');
		div.innerHTML = htmlBody;
		plainText = div.textContent || div.innerText || '';
	}
</script>

<svelte:head>
	<title>New Campaign | Outlet</title>
</svelte:head>

<div class="p-6 max-w-5xl mx-auto">
	<!-- Header -->
	<div class="flex items-center gap-4 mb-6">
		<Button type="ghost" size="sm" onclick={() => goto(`${basePath}/campaigns`)}>
			<ArrowLeft class="h-4 w-4" />
		</Button>
		<div class="flex-1">
			<h1 class="text-2xl font-semibold text-text">New Campaign</h1>
			<p class="mt-1 text-sm text-text-muted">Create and send an email campaign to your subscribers</p>
		</div>
	</div>

	{#if error}
		<Alert type="error" title="Error" class="mb-6" onclose={() => (error = '')}>
			<p>{error}</p>
		</Alert>
	{/if}

	{#if loading}
		<div class="flex justify-center py-12">
			<LoadingSpinner size="lg" />
		</div>
	{:else}
		<div class="space-y-6">
			<!-- Campaign Details -->
			<Card>
				<h2 class="text-lg font-medium mb-4">Campaign Details</h2>
				<div class="space-y-4">
					<Input
						label="Campaign Name"
						placeholder="e.g., January Newsletter"
						bind:value={name}
						required
					/>
					<Input
						label="Subject Line"
						placeholder="e.g., Check out what's new this month!"
						bind:value={subject}
						required
					/>
					<Input
						label="Preview Text"
						placeholder="Text shown after the subject in inbox previews"
						bind:value={previewText}
					/>
				</div>
			</Card>

			<!-- Sender Details -->
			<Card>
				<h2 class="text-lg font-medium mb-4">Sender Details</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<Input
						label="From Name"
						placeholder="e.g., John from Acme"
						bind:value={fromName}
					/>
					<Input
						label="From Email"
						type="email"
						placeholder="e.g., john@acme.com"
						bind:value={fromEmail}
					/>
					<Input
						label="Reply-To Email"
						type="email"
						placeholder="e.g., support@acme.com"
						bind:value={replyTo}
					/>
				</div>
				<p class="text-xs text-text-muted mt-2">
					Leave blank to use your organization's default sender settings.
				</p>
			</Card>

			<!-- Recipients -->
			<Card>
				<h2 class="text-lg font-medium mb-4">Recipients</h2>
				{#if lists.length === 0}
					<p class="text-text-muted">No lists available. Create a list first to send campaigns.</p>
				{:else}
					<div class="space-y-4">
						<!-- Send to Lists -->
						<div>
							<label class="label">
								<span class="label-text font-medium">Send to these lists</span>
							</label>
							<div class="grid grid-cols-1 md:grid-cols-2 gap-2 mt-2">
								{#each lists as list}
									<label class="flex items-center gap-3 p-3 rounded-lg border border-base-300 hover:bg-base-200 cursor-pointer transition-colors">
										<input
											type="checkbox"
											checked={selectedListIds.includes(list.id)}
											onchange={() => toggleList(list.id)}
											class="checkbox checkbox-primary checkbox-sm"
										/>
										<div class="flex-1 min-w-0">
											<span class="font-medium text-text block truncate">{list.name}</span>
											<span class="text-xs text-text-muted">{list.subscriber_count.toLocaleString()} subscribers</span>
										</div>
									</label>
								{/each}
							</div>
							{#if selectedListIds.length > 0}
								<p class="text-sm text-text-muted mt-2">
									Total recipients: <strong class="text-text">{totalSubscribers().toLocaleString()}</strong>
								</p>
							{/if}
						</div>

						<!-- Exclude Lists -->
						{#if selectedListIds.length > 0 && lists.length > 1}
							<div>
								<label class="label">
									<span class="label-text font-medium">Exclude subscribers from these lists (optional)</span>
								</label>
								<div class="grid grid-cols-1 md:grid-cols-2 gap-2 mt-2">
									{#each lists.filter((l) => !selectedListIds.includes(l.id)) as list}
										<label class="flex items-center gap-3 p-3 rounded-lg border border-base-300 hover:bg-base-200 cursor-pointer transition-colors">
											<input
												type="checkbox"
												checked={excludeListIds.includes(list.id)}
												onchange={() => toggleExcludeList(list.id)}
												class="checkbox checkbox-warning checkbox-sm"
											/>
											<div class="flex-1 min-w-0">
												<span class="font-medium text-text block truncate">{list.name}</span>
												<span class="text-xs text-text-muted">{list.subscriber_count.toLocaleString()} subscribers</span>
											</div>
										</label>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				{/if}
			</Card>

			<!-- Template Selection -->
			{#if designs.length > 0}
				<Card>
					<h2 class="text-lg font-medium mb-4">Start from Template</h2>
					<div class="grid grid-cols-1 md:grid-cols-3 gap-3">
						<button
							type="button"
							onclick={() => handleDesignSelect('')}
							class="p-4 rounded-lg border-2 text-left transition-colors {selectedDesignId === '' ? 'border-primary bg-primary/5' : 'border-base-300 hover:border-base-content/20'}"
						>
							<FileText class="h-6 w-6 mb-2 text-text-muted" />
							<span class="font-medium text-text block">Blank</span>
							<span class="text-xs text-text-muted">Start from scratch</span>
						</button>
						{#each designs as design}
							<button
								type="button"
								onclick={() => handleDesignSelect(design.id)}
								class="p-4 rounded-lg border-2 text-left transition-colors {selectedDesignId === design.id ? 'border-primary bg-primary/5' : 'border-base-300 hover:border-base-content/20'}"
							>
								<FileText class="h-6 w-6 mb-2 text-primary" />
								<span class="font-medium text-text block truncate">{design.name}</span>
								<span class="text-xs text-text-muted">{design.category}</span>
							</button>
						{/each}
					</div>
				</Card>
			{/if}

			<!-- Email Content -->
			<Card>
				<h2 class="text-lg font-medium mb-4">Email Content</h2>
				<div class="space-y-4">
					<div>
						<label class="label">
							<span class="label-text font-medium">HTML Content</span>
						</label>
						<HtmlEditor bind:value={htmlBody} minHeight="400px" />
					</div>

					<div>
						<div class="flex items-center justify-between mb-2">
							<label class="label">
								<span class="label-text font-medium">Plain Text Version</span>
							</label>
							<Button type="ghost" size="xs" onclick={generatePlainText}>
								Generate from HTML
							</Button>
						</div>
						<Textarea
							bind:value={plainText}
							placeholder="Plain text version for email clients that don't support HTML..."
							rows={6}
						/>
						<p class="text-xs text-text-muted mt-1">
							Optional. A plain text version will be auto-generated if not provided.
						</p>
					</div>
				</div>
			</Card>

			<!-- Tracking Options -->
			<Card>
				<h2 class="text-lg font-medium mb-4">Tracking</h2>
				<div class="space-y-3">
					<Checkbox
						bind:checked={trackOpens}
						label="Track opens"
					/>
					<p class="text-xs text-text-muted ml-6 -mt-1">
						Adds a tracking pixel to measure how many recipients open the email.
					</p>

					<Checkbox
						bind:checked={trackClicks}
						label="Track clicks"
					/>
					<p class="text-xs text-text-muted ml-6 -mt-1">
						Wraps links to track which links are clicked.
					</p>
				</div>
			</Card>

			<!-- Template Tags Help -->
			<Card>
				<h2 class="text-lg font-medium mb-4">Available Template Tags</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
					<div>
						<code class="bg-base-200 px-2 py-1 rounded text-xs">{'{{email}}'}</code>
						<span class="text-text-muted ml-2">Subscriber's email</span>
					</div>
					<div>
						<code class="bg-base-200 px-2 py-1 rounded text-xs">{'{{name}}'}</code>
						<span class="text-text-muted ml-2">Subscriber's name</span>
					</div>
					<div>
						<code class="bg-base-200 px-2 py-1 rounded text-xs">{'{{unsubscribe_url}}'}</code>
						<span class="text-text-muted ml-2">Unsubscribe link</span>
					</div>
					<div>
						<code class="bg-base-200 px-2 py-1 rounded text-xs">{'{{web_version_url}}'}</code>
						<span class="text-text-muted ml-2">Web version link</span>
					</div>
				</div>
			</Card>

			<!-- Actions -->
			<div class="flex flex-col sm:flex-row gap-3 pt-4 border-t border-base-300">
				<Button type="secondary" onclick={() => goto(`${basePath}/campaigns`)}>
					Cancel
				</Button>
				<div class="flex-1"></div>
				<Button
					type="secondary"
					onclick={() => handleSave(false)}
					disabled={saving || !isValid()}
				>
					<Save class="h-4 w-4 mr-2" />
					{saving ? 'Saving...' : 'Save as Draft'}
				</Button>
			</div>
		</div>
	{/if}
</div>
