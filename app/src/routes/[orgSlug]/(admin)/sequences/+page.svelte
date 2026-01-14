<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import {
		listSequences,
		listLists,
		createSequence,
		deleteSequence,
		updateSequence,
		type SequenceInfo,
		type ListInfo
	} from '$lib/api';
	import {
		Button,
		Alert,
		LoadingSpinner,
		Badge,
		EmptyState,
		SearchInput,
		DropdownMenu,
		SlideForm,
		Input,
		Select,
		Toggle
	} from '$lib/components/ui';
	import { Plus, Workflow, MoreVertical, Trash2, Edit, Play, Pause, Copy } from 'lucide-svelte';

	let loading = $state(true);
	let sequences = $state<SequenceInfo[]>([]);
	let lists = $state<ListInfo[]>([]);
	let error = $state('');
	let searchQuery = $state('');

	// Create form
	let showCreateForm = $state(false);
	let creating = $state(false);
	let newName = $state('');
	let newSlug = $state('');
	let newListId = $state('');
	let newTriggerEvent = $state('signup');
	let newIsActive = $state(true);

	// Build base path with orgSlug
	let basePath = $derived(`/${$page.params.orgSlug}`);

	// Filtered sequences based on search
	let filteredSequences = $derived(
		sequences.filter(seq =>
			seq.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			seq.slug.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	$effect(() => {
		loadData();
	});

	async function loadData() {
		loading = true;
		error = '';

		try {
			const [seqResponse, listResponse] = await Promise.all([
				listSequences(),
				listLists()
			]);
			sequences = seqResponse.sequences || [];
			lists = listResponse.lists || [];
		} catch (err) {
			console.error('Failed to load sequences:', err);
			error = 'Failed to load sequences';
		} finally {
			loading = false;
		}
	}

	function openCreateForm() {
		newName = '';
		newSlug = '';
		newListId = lists[0]?.id || '';
		newTriggerEvent = 'signup';
		newIsActive = true;
		showCreateForm = true;
	}

	function closeCreateForm() {
		showCreateForm = false;
	}

	function generateSlug(name: string): string {
		return name
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, '-')
			.replace(/^-|-$/g, '');
	}

	function handleNameInput() {
		newSlug = generateSlug(newName);
	}

	async function submitCreate() {
		if (!newName || !newSlug || !newListId) return;

		creating = true;
		error = '';
		try {
			await createSequence({
				name: newName,
				slug: newSlug,
				list_id: newListId,
				trigger_event: newTriggerEvent,
				is_active: newIsActive,
				sequence_type: 'lifecycle'
			});
			closeCreateForm();
			await loadData();
		} catch (err: any) {
			console.error('Failed to create sequence:', err);
			error = err.message || 'Failed to create sequence';
		} finally {
			creating = false;
		}
	}

	async function handleDelete(sequence: SequenceInfo) {
		if (!confirm(`Delete "${sequence.name}"? This cannot be undone.`)) return;

		try {
			await deleteSequence({}, sequence.id);
			await loadData();
		} catch (err: any) {
			console.error('Failed to delete sequence:', err);
			error = err.message || 'Failed to delete sequence';
		}
	}

	async function toggleSequenceActive(sequence: SequenceInfo) {
		try {
			await updateSequence({}, { is_active: !sequence.is_active }, sequence.id);
			await loadData();
		} catch (err: any) {
			console.error('Failed to toggle sequence:', err);
			error = err.message || 'Failed to update sequence';
		}
	}

	function getTriggerLabel(trigger: string): string {
		switch (trigger) {
			case 'signup': return 'On signup';
			case 'purchase': return 'On purchase';
			case 'tag_added': return 'Tag added';
			case 'manual': return 'Manual';
			default: return trigger;
		}
	}
</script>

<svelte:head>
	<title>Sequences | Outlet</title>
</svelte:head>

<div class="p-6 max-w-5xl mx-auto">
	<!-- Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
		<div>
			<h1 class="text-2xl font-semibold text-text">Sequences</h1>
			<p class="mt-1 text-sm text-text-muted">Automated email sequences triggered by events</p>
		</div>
		<Button type="primary" onclick={openCreateForm} disabled={showCreateForm || lists.length === 0}>
			<Plus class="mr-2 h-4 w-4" />
			New Sequence
		</Button>
	</div>

	{#if error}
		<Alert type="error" title="Error" class="mb-4">
			<p>{error}</p>
		</Alert>
	{/if}

	<!-- Inline Create Form -->
	<SlideForm bind:show={showCreateForm} title="Create New Sequence">
		<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
			<div>
				<label for="seq-name" class="form-label">Name</label>
				<Input
					id="seq-name"
					type="text"
					bind:value={newName}
					oninput={handleNameInput}
					placeholder="Welcome Series"
				/>
			</div>
			<div>
				<label for="seq-slug" class="form-label">Slug</label>
				<Input
					id="seq-slug"
					type="text"
					bind:value={newSlug}
					placeholder="welcome-series"
				/>
			</div>
			<div>
				<label for="seq-list" class="form-label">Email List</label>
				<Select
					id="seq-list"
					bind:value={newListId}
					options={lists.map(l => ({ value: l.id, label: l.name }))}
				/>
			</div>
			<div>
				<label for="seq-trigger" class="form-label">Trigger</label>
				<Select
					id="seq-trigger"
					bind:value={newTriggerEvent}
					options={[
						{ value: 'signup', label: 'On signup' },
						{ value: 'purchase', label: 'On purchase' },
						{ value: 'tag_added', label: 'Tag added' },
						{ value: 'manual', label: 'Manual enrollment' }
					]}
				/>
			</div>
			<div class="sm:col-span-2 flex items-center gap-3">
				<Toggle bind:checked={newIsActive} />
				<span class="text-sm text-text">Active (starts enrolling subscribers immediately)</span>
			</div>
		</div>

		{#snippet footer()}
			<Button type="secondary" onclick={closeCreateForm} disabled={creating}>
				Cancel
			</Button>
			<Button type="primary" onclick={submitCreate} disabled={!newName || !newSlug || !newListId || creating}>
				{creating ? 'Creating...' : 'Create Sequence'}
			</Button>
		{/snippet}
	</SlideForm>

	{#if loading}
		<div class="flex justify-center py-12">
			<LoadingSpinner size="large" />
		</div>
	{:else if lists.length === 0}
		<Alert type="info" title="Create a list first">
			<p>You need at least one email list before creating sequences.</p>
			<Button type="primary" class="mt-3" onclick={() => goto(`${basePath}/lists`)}>
				Go to Lists
			</Button>
		</Alert>
	{:else if sequences.length === 0}
		<EmptyState
			icon={Workflow}
			title="No sequences yet"
			description="Create automated email sequences to nurture your subscribers."
		>
			<Button type="primary" onclick={openCreateForm}>
				<Plus class="mr-2 h-4 w-4" />
				Create Sequence
			</Button>
		</EmptyState>
	{:else}
		<!-- Search -->
		{#if sequences.length > 3}
			<div class="mb-4">
				<SearchInput
					bind:value={searchQuery}
					placeholder="Search sequences..."
				/>
			</div>
		{/if}

		<!-- Sequences Table -->
		<div class="sequences-table">
			<table class="w-full">
				<thead>
					<tr>
						<th class="text-left">Sequence</th>
						<th class="text-left">List</th>
						<th class="text-left">Trigger</th>
						<th class="text-center">Status</th>
						<th class="text-right w-10"></th>
					</tr>
				</thead>
				<tbody>
					{#each filteredSequences as sequence}
						<tr>
							<td>
								<a href="{basePath}/sequences/{sequence.id}" class="block group">
									<div class="flex items-center gap-3">
										<div class="sequences-icon" class:active={sequence.is_active}>
											<Workflow class="h-4 w-4" />
										</div>
										<div class="min-w-0 flex-1">
											<span class="font-medium text-text group-hover:text-primary transition-colors block truncate">
												{sequence.name}
											</span>
											<span class="text-sm text-text-muted">{sequence.slug}</span>
										</div>
									</div>
								</a>
							</td>
							<td>
								<span class="text-text">{sequence.list_name || '-'}</span>
							</td>
							<td>
								<Badge variant="default" size="sm">
									{getTriggerLabel(sequence.trigger_event)}
								</Badge>
							</td>
							<td class="text-center">
								{#if sequence.is_active}
									<Badge variant="success" size="sm">Active</Badge>
								{:else}
									<Badge variant="default" size="sm">Paused</Badge>
								{/if}
							</td>
							<td class="text-right">
								<DropdownMenu
									trigger={{
										icon: MoreVertical,
										class: 'p-1.5 rounded-lg hover:bg-bg-secondary transition-colors text-text-muted hover:text-text'
									}}
									items={[
										{
											label: 'Edit',
											icon: Edit,
											onclick: () => goto(`${basePath}/sequences/${sequence.id}`)
										},
										{
											label: sequence.is_active ? 'Pause' : 'Activate',
											icon: sequence.is_active ? Pause : Play,
											onclick: () => toggleSequenceActive(sequence)
										},
										{
											label: 'Duplicate',
											icon: Copy,
											onclick: () => goto(`${basePath}/sequences/new?duplicate=${sequence.id}`)
										},
										{ divider: true },
										{
											label: 'Delete',
											icon: Trash2,
											variant: 'danger',
											onclick: () => handleDelete(sequence)
										}
									]}
								/>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		{#if searchQuery && filteredSequences.length === 0}
			<div class="text-center py-8 text-text-muted">
				No sequences match "{searchQuery}"
			</div>
		{/if}
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.sequences {
		.sequences-table {
			@apply rounded-lg border overflow-hidden;
			background-color: var(--color-bg);
			border-color: var(--color-border);
		}

		.sequences-table thead {
			background-color: var(--color-bg-secondary);
		}

		.sequences-table th {
			@apply px-4 py-3 text-xs font-medium uppercase tracking-wide;
			color: var(--color-text-muted);
		}

		.sequences-table td {
			@apply px-4 py-4 border-t;
			border-color: var(--color-border);
		}

		.sequences-table tbody tr:hover {
			background-color: var(--color-bg-secondary);
		}

		.sequences-icon {
			@apply flex items-center justify-center h-10 w-10 rounded-lg;
			background-color: var(--color-bg-secondary);
			color: var(--color-text-muted);
		}

		.sequences-icon.active {
			background-color: color-mix(in srgb, var(--color-success) 15%, transparent);
			color: var(--color-success);
		}
	}
</style>
