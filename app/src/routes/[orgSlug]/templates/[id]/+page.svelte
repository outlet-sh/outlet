<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import {
		getEmailDesign,
		updateEmailDesign,
		deleteEmailDesign,
		type EmailDesignInfo
	} from '$lib/api';
	import {
		Button,
		Card,
		Input,
		Alert,
		LoadingSpinner,
		Badge,
		Select,
		Toggle,
		AlertDialog,
		HtmlEditor
	} from '$lib/components/ui';
	import {
		ArrowLeft,
		FileText,
		Save,
		Trash2,
		Eye
	} from 'lucide-svelte';

	// Get template ID from URL
	let templateId = $derived($page.params.id);
	let basePath = $derived(`/${$page.params.orgSlug}`);

	// Main state
	let loading = $state(true);
	let template = $state<EmailDesignInfo | null>(null);
	let error = $state('');

	// Edit state
	let editName = $state('');
	let editSlug = $state('');
	let editDescription = $state('');
	let editCategory = $state('general');
	let editHtmlBody = $state('');
	let editPlainText = $state('');
	let editIsActive = $state(true);
	let saving = $state(false);
	let saved = $state(false);

	// Preview state
	let showPreview = $state(false);

	// Delete state
	let showDeleteConfirm = $state(false);
	let deleting = $state(false);

	// Category options
	const categories = [
		{ value: 'general', label: 'General' },
		{ value: 'newsletter', label: 'Newsletter' },
		{ value: 'promotional', label: 'Promotional' },
		{ value: 'transactional', label: 'Transactional' }
	];

	// Check if there are unsaved changes
	let hasChanges = $derived(
		template !== null && (
			editName !== template.name ||
			editDescription !== (template.description || '') ||
			editCategory !== template.category ||
			editHtmlBody !== template.html_body ||
			editPlainText !== (template.plain_text || '') ||
			editIsActive !== template.is_active
		)
	);

	$effect(() => {
		loadTemplate();
	});

	async function loadTemplate() {
		if (!templateId) return;
		loading = true;
		error = '';
		try {
			template = await getEmailDesign({}, templateId);
			editName = template.name;
			editSlug = template.slug;
			editDescription = template.description || '';
			editCategory = template.category;
			editHtmlBody = template.html_body;
			editPlainText = template.plain_text || '';
			editIsActive = template.is_active;
		} catch (err) {
			console.error('Failed to load template:', err);
			error = 'Failed to load template';
		} finally {
			loading = false;
		}
	}

	async function saveTemplate() {
		if (!template || !templateId || !editName.trim()) return;
		saving = true;
		saved = false;
		error = '';
		try {
			const updated = await updateEmailDesign({}, {
				name: editName.trim(),
				description: editDescription.trim() || undefined,
				category: editCategory,
				html_body: editHtmlBody,
				plain_text: editPlainText || undefined,
				is_active: editIsActive
			}, templateId);
			template = updated;
			saved = true;
			setTimeout(() => { saved = false; }, 2000);
		} catch (err: any) {
			console.error('Failed to save template:', err);
			error = err.message || 'Failed to save template';
		} finally {
			saving = false;
		}
	}

	async function executeDelete() {
		if (!templateId) return;
		deleting = true;
		try {
			await deleteEmailDesign({}, templateId);
			goto(`${basePath}/templates`);
		} catch (err: any) {
			error = err.message || 'Failed to delete template';
		} finally {
			deleting = false;
		}
	}

	function getCategoryVariant(category: string): 'default' | 'info' | 'success' | 'warning' {
		switch (category) {
			case 'newsletter': return 'info';
			case 'promotional': return 'warning';
			case 'transactional': return 'success';
			default: return 'default';
		}
	}

	function getCategoryLabel(category: string): string {
		const found = categories.find(c => c.value === category);
		return found?.label || category;
	}
</script>

<svelte:head>
	<title>{template?.name || 'Template'} | Outlet</title>
</svelte:head>

<div class="p-6 max-w-6xl mx-auto">
	{#if loading}
		<div class="flex justify-center py-12">
			<LoadingSpinner size="large" />
		</div>
	{:else if !template}
		<Alert type="error" title="Template not found">
			<p>The requested template could not be found.</p>
			<Button type="primary" class="mt-3" onclick={() => goto(`${basePath}/templates`)}>
				Back to Templates
			</Button>
		</Alert>
	{:else}
		<!-- Header -->
		<div class="flex items-center gap-4 mb-6">
			<a href="{basePath}/templates" class="p-2 rounded-lg hover:bg-bg-secondary transition-colors text-text-muted hover:text-text">
				<ArrowLeft class="h-5 w-5" />
			</a>
			<div class="flex-1 flex items-center gap-3">
				<div class="data-table-icon">
					<FileText class="h-5 w-5" />
				</div>
				<div>
					<h1 class="text-2xl font-semibold text-text">{template.name}</h1>
					<div class="flex items-center gap-2 mt-1">
						<Badge variant={getCategoryVariant(template.category)} size="sm">
							{getCategoryLabel(template.category)}
						</Badge>
						{#if template.is_active}
							<Badge variant="success" size="sm">Active</Badge>
						{:else}
							<Badge variant="default" size="sm">Inactive</Badge>
						{/if}
					</div>
				</div>
			</div>
			<div class="flex items-center gap-2">
				<Button type="secondary" onclick={() => showPreview = true}>
					<Eye class="mr-2 h-4 w-4" />
					Preview
				</Button>
				<Button
					type="primary"
					onclick={saveTemplate}
					disabled={saving || !hasChanges}
				>
					<Save class="mr-2 h-4 w-4" />
					{#if saving}
						Saving...
					{:else if saved}
						Saved!
					{:else}
						Save
					{/if}
				</Button>
			</div>
		</div>

		{#if error}
			<Alert type="error" title="Error" class="mb-4">
				<p>{error}</p>
			</Alert>
		{/if}

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<!-- Left: Settings -->
			<div class="lg:col-span-1 space-y-4">
				<Card>
					<h2 class="text-lg font-medium text-text mb-4">Settings</h2>
					<div class="space-y-4">
						<div>
							<label for="template-name" class="form-label">Name</label>
							<Input id="template-name" type="text" bind:value={editName} />
						</div>
						<div>
							<label class="form-label">Slug</label>
							<p class="text-sm text-text-muted bg-bg-secondary px-3 py-2 rounded-md">{editSlug}</p>
							<p class="mt-1 text-xs text-text-muted">Cannot be changed after creation</p>
						</div>
						<div>
							<label for="template-description" class="form-label">Description</label>
							<Input id="template-description" type="text" bind:value={editDescription} placeholder="Optional description" />
						</div>
						<div>
							<label for="template-category" class="form-label">Category</label>
							<Select
								id="template-category"
								bind:value={editCategory}
								options={categories}
							/>
						</div>
						<div class="flex items-center justify-between">
							<span class="text-sm text-text">Active</span>
							<Toggle bind:checked={editIsActive} />
						</div>
					</div>
				</Card>

				<Card>
					<h2 class="text-lg font-medium text-red-600 mb-4">Danger Zone</h2>
					<p class="text-sm text-text-muted mb-4">
						Permanently delete this template. This action cannot be undone.
					</p>
					<Button type="danger" onclick={() => showDeleteConfirm = true}>
						<Trash2 class="mr-2 h-4 w-4" />
						Delete Template
					</Button>
				</Card>
			</div>

			<!-- Right: Editor -->
			<div class="lg:col-span-2">
				<Card class="h-full">
					<h2 class="text-lg font-medium text-text mb-4">Email Content</h2>
					<HtmlEditor
						bind:value={editHtmlBody}
						placeholder="Start designing your email..."
						minHeight="400px"
					/>
				</Card>
			</div>
		</div>
	{/if}
</div>

<!-- Preview Modal -->
{#if showPreview && template}
	<div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
		<div class="bg-bg-primary rounded-xl shadow-xl max-w-3xl w-full max-h-[90vh] flex flex-col">
			<div class="flex items-center justify-between p-4 border-b border-border">
				<h3 class="text-lg font-medium text-text">Preview: {editName}</h3>
				<Button type="secondary" onclick={() => showPreview = false}>Close</Button>
			</div>
			<div class="flex-1 overflow-auto p-4 bg-white">
				{@html editHtmlBody}
			</div>
		</div>
	</div>
{/if}

<AlertDialog
	bind:open={showDeleteConfirm}
	title="Delete Template"
	description={`Are you sure you want to delete "${template?.name}"? This action cannot be undone.`}
	actionLabel={deleting ? 'Deleting...' : 'Delete'}
	actionType="danger"
	onAction={executeDelete}
	onCancel={() => showDeleteConfirm = false}
/>
