<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import {
		getSequence,
		getList,
		listLists,
		listSequences,
		updateSequence,
		deleteSequence,
		createTemplate,
		updateTemplate,
		deleteTemplate,
		type SequenceInfo,
		type TemplateInfo,
		type ListInfo
	} from '$lib/api';
	import {
		Button,
		Card,
		Input,
		Alert,
		LoadingSpinner,
		Badge,
		EmptyState,
		SlideForm,
		Select,
		Toggle,
		AlertDialog,
		Textarea
	} from '$lib/components/ui';
	import {
		ArrowLeft,
		Plus,
		Mail,
		Trash2,
		Edit,
		Play,
		Pause,
		GripVertical,
		Clock,
		ArrowRight
	} from 'lucide-svelte';

	// Route params
	let listId = $derived($page.params.id);
	let seqId = $derived($page.params.seqId);
	let basePath = $derived(`/${$page.params.orgSlug}`);

	// State
	let loading = $state(true);
	let sequence = $state<SequenceInfo | null>(null);
	let templates = $state<TemplateInfo[]>([]);
	let list = $state<ListInfo | null>(null);
	let allLists = $state<ListInfo[]>([]);
	let allSequences = $state<SequenceInfo[]>([]);
	let error = $state('');

	// Edit sequence state
	let editName = $state('');
	let editTrigger = $state('signup');
	let editActive = $state(true);
	let editListId = $state('');
	let editOnCompletionSequenceId = $state<string>('');
	let saving = $state(false);
	let saved = $state(false);

	// Create email state
	let showCreateEmail = $state(false);
	let creatingEmail = $state(false);
	let newSubject = $state('');
	let newDelay = $state(0);
	let newHtmlBody = $state('');

	// Edit email state
	let editingEmail = $state<TemplateInfo | null>(null);
	let editSubject = $state('');
	let editDelay = $state(0);
	let editHtmlBody = $state('');
	let savingEmail = $state(false);

	// Delete states
	let showDeleteSequence = $state(false);
	let deletingSequence = $state(false);
	let deletingEmailId = $state<string | null>(null);

	$effect(() => {
		loadData();
	});

	async function loadData() {
		loading = true;
		error = '';
		try {
			const [seqResponse, listResponse, allSeqResponse, allListsResponse] = await Promise.all([
				getSequence({}, seqId),
				getList({}, listId),
				listSequences(),
				listLists()
			]);
			sequence = seqResponse.sequence;
			templates = (seqResponse.templates || []).sort((a, b) => a.position - b.position);
			list = listResponse;
			allLists = allListsResponse.lists || [];
			// Filter out the current sequence from completion options
			allSequences = (allSeqResponse.sequences || []).filter(s => s.id !== seqId);

			// Set edit state
			editName = sequence.name;
			editTrigger = sequence.trigger_event;
			editActive = sequence.is_active;
			editListId = sequence.list_id || listId;
			editOnCompletionSequenceId = sequence.on_completion_sequence_id || '';
		} catch (err) {
			console.error('Failed to load sequence:', err);
			error = 'Failed to load autoresponder';
		} finally {
			loading = false;
		}
	}

	async function saveSequence() {
		if (!sequence || !editName.trim()) return;
		saving = true;
		saved = false;
		error = '';
		try {
			const listChanged = editListId !== listId;
			const updatedSeq = await updateSequence({}, {
				name: editName.trim(),
				trigger_event: editTrigger,
				is_active: editActive,
				list_id: editListId,
				on_completion_sequence_id: editOnCompletionSequenceId || null
			}, seqId);
			sequence.name = editName.trim();
			sequence.trigger_event = editTrigger;
			sequence.is_active = editActive;
			sequence.list_id = editListId;
			sequence.on_completion_sequence_id = editOnCompletionSequenceId || '';
			sequence.on_completion_sequence_name = updatedSeq.on_completion_sequence_name || '';
			saved = true;

			// If list changed, redirect to the new list's sequence page
			if (listChanged) {
				goto(`${basePath}/lists/${editListId}/autoresponders/${seqId}`);
			} else {
				setTimeout(() => { saved = false; }, 2000);
			}
		} catch (err: any) {
			error = err.message || 'Failed to save';
		} finally {
			saving = false;
		}
	}

	async function toggleActive() {
		if (!sequence) return;
		try {
			await updateSequence({}, { is_active: !sequence.is_active }, seqId);
			sequence.is_active = !sequence.is_active;
			editActive = sequence.is_active;
		} catch (err: any) {
			error = err.message || 'Failed to toggle status';
		}
	}

	function openCreateEmail() {
		newSubject = '';
		newDelay = templates.length === 0 ? 0 : 24; // First email immediately, others 24h
		newHtmlBody = '';
		showCreateEmail = true;
	}

	async function submitCreateEmail() {
		if (!newSubject.trim()) return;
		creatingEmail = true;
		error = '';
		try {
			await createTemplate({
				sequence_id: seqId,
				subject: newSubject.trim(),
				delay_hours: newDelay,
				html_body: newHtmlBody || '<p>Email content here</p>',
				is_active: true,
				position: templates.length + 1
			});
			showCreateEmail = false;
			await loadData();
		} catch (err: any) {
			error = err.message || 'Failed to create email';
		} finally {
			creatingEmail = false;
		}
	}

	function startEditEmail(template: TemplateInfo) {
		editingEmail = template;
		editSubject = template.subject;
		editDelay = template.delay_hours;
		editHtmlBody = template.html_body;
	}

	async function saveEmail() {
		if (!editingEmail || !editSubject.trim()) return;
		savingEmail = true;
		error = '';
		try {
			await updateTemplate({}, {
				subject: editSubject.trim(),
				delay_hours: editDelay,
				html_body: editHtmlBody
			}, editingEmail.id);
			editingEmail = null;
			await loadData();
		} catch (err: any) {
			error = err.message || 'Failed to save email';
		} finally {
			savingEmail = false;
		}
	}

	async function handleDeleteEmail(templateId: string) {
		deletingEmailId = templateId;
		try {
			await deleteTemplate({}, templateId);
			await loadData();
		} catch (err: any) {
			error = err.message || 'Failed to delete email';
		} finally {
			deletingEmailId = null;
		}
	}

	async function executeDeleteSequence() {
		deletingSequence = true;
		try {
			await deleteSequence({}, seqId);
			goto(`${basePath}/lists/${listId}?tab=autoresponders`);
		} catch (err: any) {
			error = err.message || 'Failed to delete autoresponder';
		} finally {
			deletingSequence = false;
		}
	}

	function formatDelay(hours: number): string {
		if (hours === 0) return 'Immediately';
		if (hours < 24) return `${hours} hour${hours === 1 ? '' : 's'} after`;
		const days = Math.floor(hours / 24);
		const remainingHours = hours % 24;
		if (remainingHours === 0) return `${days} day${days === 1 ? '' : 's'} after`;
		return `${days}d ${remainingHours}h after`;
	}

	function getTriggerLabel(trigger: string): string {
		switch (trigger) {
			case 'signup': return 'On signup';
			case 'manual': return 'Manual enrollment';
			default: return trigger;
		}
	}
</script>

<svelte:head>
	<title>{sequence?.name || 'Autoresponder'} | Outlet</title>
</svelte:head>

<div class="p-6 max-w-4xl mx-auto">
	{#if loading}
		<div class="flex justify-center py-12">
			<LoadingSpinner size="large" />
		</div>
	{:else if !sequence}
		<Alert type="error" title="Not found">
			<p>The requested autoresponder could not be found.</p>
			<Button type="primary" class="mt-3" onclick={() => goto(`${basePath}/lists/${listId}`)}>
				Back to List
			</Button>
		</Alert>
	{:else}
		<!-- Header -->
		<div class="flex items-center gap-4 mb-6">
			<a href="{basePath}/lists/{listId}" class="p-2 rounded-lg hover:bg-bg-secondary transition-colors text-text-muted hover:text-text">
				<ArrowLeft class="h-5 w-5" />
			</a>
			<div class="flex-1">
				<div class="flex items-center gap-2 text-sm text-text-muted mb-1">
					<span>{list?.name}</span>
					<span>/</span>
					<span>Autoresponder</span>
				</div>
				<div class="flex items-center gap-3">
					<h1 class="text-2xl font-semibold text-text">{sequence.name}</h1>
					{#if sequence.is_active}
						<Badge variant="success" size="sm">Active</Badge>
					{:else}
						<Badge variant="default" size="sm">Paused</Badge>
					{/if}
				</div>
				{#if sequence.on_completion_sequence_name}
					<div class="flex items-center gap-1 text-sm text-text-muted mt-1">
						<span>On completion:</span>
						<ArrowRight class="h-3 w-3" />
						<span class="text-primary">{sequence.on_completion_sequence_name}</span>
					</div>
				{/if}
			</div>
			<Button
				type={sequence.is_active ? 'secondary' : 'primary'}
				onclick={toggleActive}
			>
				{#if sequence.is_active}
					<Pause class="mr-2 h-4 w-4" />
					Pause
				{:else}
					<Play class="mr-2 h-4 w-4" />
					Activate
				{/if}
			</Button>
		</div>

		{#if error}
			<Alert type="error" title="Error" class="mb-4">
				<p>{error}</p>
			</Alert>
		{/if}

		<!-- Sequence Settings -->
		<Card class="mb-6">
			<h2 class="text-lg font-medium text-text mb-4">Settings</h2>
			<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
				<div>
					<label for="seq-name" class="form-label">Name</label>
					<Input id="seq-name" type="text" bind:value={editName} />
				</div>
				<div>
					<label for="seq-trigger" class="form-label">Trigger</label>
					<Select
						id="seq-trigger"
						bind:value={editTrigger}
						options={[
							{ value: 'signup', label: 'On signup' },
							{ value: 'manual', label: 'Manual enrollment' }
						]}
					/>
				</div>
			</div>
			<div class="mt-4">
				<label for="seq-list" class="form-label">List</label>
				<p class="text-sm text-text-muted mb-2">Move this autoresponder to a different list.</p>
				<Select
					id="seq-list"
					bind:value={editListId}
					options={allLists.map(l => ({ value: String(l.id), label: l.name }))}
				/>
			</div>
			<div class="mt-4">
				<label for="seq-on-completion" class="form-label">On Completion</label>
				<p class="text-sm text-text-muted mb-2">Automatically enroll contacts in another sequence when they complete this one.</p>
				<Select
					id="seq-on-completion"
					bind:value={editOnCompletionSequenceId}
					options={[
						{ value: '', label: 'None (do nothing)' },
						...allSequences.map(s => ({ value: s.id, label: s.name }))
					]}
				/>
			</div>
			<div class="mt-4 flex items-center gap-4">
				<Button
					type="primary"
					onclick={saveSequence}
					disabled={saving || (editName === sequence.name && editTrigger === sequence.trigger_event && editListId === (sequence.list_id || listId) && editOnCompletionSequenceId === (sequence.on_completion_sequence_id || ''))}
				>
					{#if saving}
						Saving...
					{:else if saved}
						Saved!
					{:else}
						Save Changes
					{/if}
				</Button>
				<Button type="danger" onclick={() => showDeleteSequence = true}>
					<Trash2 class="mr-2 h-4 w-4" />
					Delete
				</Button>
			</div>
		</Card>

		<!-- Emails Section -->
		<div class="flex items-center justify-between mb-4">
			<h2 class="text-lg font-medium text-text">Emails</h2>
			<Button type="primary" onclick={openCreateEmail} disabled={showCreateEmail}>
				<Plus class="mr-2 h-4 w-4" />
				Add Email
			</Button>
		</div>

		<!-- Create Email Form -->
		<SlideForm bind:show={showCreateEmail} title="Add Email to Sequence">
			<div class="space-y-4">
				<div>
					<label for="email-subject" class="form-label">Subject</label>
					<Input id="email-subject" type="text" bind:value={newSubject} placeholder="Welcome to our newsletter!" />
				</div>
				<div>
					<label for="email-delay" class="form-label">Send</label>
					<div class="flex items-center gap-2">
						<Input id="email-delay" type="number" bind:value={newDelay} min={0} class="w-24" />
						<span class="text-text-muted">hours after {templates.length === 0 ? 'signup' : 'previous email'}</span>
					</div>
				</div>
				<div>
					<label for="email-body" class="form-label">Content (HTML)</label>
					<Textarea id="email-body" bind:value={newHtmlBody} rows={6} placeholder="<p>Your email content...</p>" />
				</div>
			</div>

			{#snippet footer()}
				<Button type="secondary" onclick={() => showCreateEmail = false} disabled={creatingEmail}>
					Cancel
				</Button>
				<Button type="primary" onclick={submitCreateEmail} disabled={!newSubject.trim() || creatingEmail}>
					{creatingEmail ? 'Adding...' : 'Add Email'}
				</Button>
			{/snippet}
		</SlideForm>

		<!-- Edit Email Form -->
		{#if editingEmail}
			<Card class="mb-4 border-primary">
				<h3 class="font-medium text-text mb-4">Edit Email</h3>
				<div class="space-y-4">
					<div>
						<label for="edit-subject" class="form-label">Subject</label>
						<Input id="edit-subject" type="text" bind:value={editSubject} />
					</div>
					<div>
						<label for="edit-delay" class="form-label">Delay (hours)</label>
						<Input id="edit-delay" type="number" bind:value={editDelay} min={0} class="w-32" />
					</div>
					<div>
						<label for="edit-body" class="form-label">Content (HTML)</label>
						<Textarea id="edit-body" bind:value={editHtmlBody} rows={8} />
					</div>
					<div class="flex gap-2">
						<Button type="primary" onclick={saveEmail} disabled={savingEmail}>
							{savingEmail ? 'Saving...' : 'Save'}
						</Button>
						<Button type="secondary" onclick={() => editingEmail = null}>
							Cancel
						</Button>
					</div>
				</div>
			</Card>
		{/if}

		<!-- Email List -->
		{#if templates.length === 0}
			<EmptyState
				icon={Mail}
				title="No emails yet"
				description="Add your first email to this autoresponder sequence."
			>
				<Button type="primary" onclick={openCreateEmail}>
					<Plus class="mr-2 h-4 w-4" />
					Add Email
				</Button>
			</EmptyState>
		{:else}
			<div class="space-y-3">
				{#each templates as template, i}
					<Card class="flex items-center gap-4 {!template.is_active ? 'opacity-60' : ''}">
						<div class="text-text-muted">
							<GripVertical class="h-5 w-5" />
						</div>
						<div class="flex items-center justify-center h-8 w-8 rounded-full bg-primary text-white text-sm font-medium">
							{i + 1}
						</div>
						<div class="flex-1 min-w-0">
							<div class="font-medium text-text truncate">{template.subject}</div>
							<div class="text-sm text-text-muted flex items-center gap-2">
								<Clock class="h-3 w-3" />
								{formatDelay(template.delay_hours)}
							</div>
						</div>
						<div class="flex items-center gap-2">
							<Button type="ghost" size="sm" onclick={() => startEditEmail(template)}>
								<Edit class="h-4 w-4" />
							</Button>
							<Button
								type="ghost"
								size="sm"
								onclick={() => handleDeleteEmail(template.id)}
								disabled={deletingEmailId === template.id}
								class="text-error hover:bg-error/10"
							>
								<Trash2 class="h-4 w-4" />
							</Button>
						</div>
					</Card>
				{/each}
			</div>
		{/if}
	{/if}
</div>

<AlertDialog
	bind:open={showDeleteSequence}
	title="Delete Autoresponder"
	description={`Are you sure you want to delete "${sequence?.name}"? This will delete all emails in this sequence. This action cannot be undone.`}
	actionLabel={deletingSequence ? 'Deleting...' : 'Delete'}
	actionType="danger"
	onAction={executeDeleteSequence}
	onCancel={() => showDeleteSequence = false}
/>
