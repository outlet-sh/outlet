<script lang="ts">
	import { untrack } from 'svelte';
	import { Input, Textarea, Tabs, Alert, Button } from '$lib/components/ui';
	import { FileCode, Wand2 } from 'lucide-svelte';
	import ConditionBuilder from './ConditionBuilder.svelte';
	import ActionBuilder from './ActionBuilder.svelte';
	import {
		type RuleDefinition,
		type ConditionGroup,
		type RuleActionConfig,
		ruleDefinitionToGruleJson,
		gruleJsonToRuleDefinition,
		generateId
	} from './types';

	let {
		name = $bindable(''),
		description = $bindable(''),
		salience = $bindable(0),
		ruleJson = $bindable(''),
		category = $bindable('support'),
		onValidChange
	}: {
		name: string;
		description: string;
		salience: number;
		ruleJson: string;
		category: string;
		onValidChange?: (valid: boolean) => void;
	} = $props();

	// Tab state
	let activeTab = $state('visual');
	const tabs = [
		{ id: 'visual', label: 'Visual Editor', icon: Wand2 },
		{ id: 'json', label: 'JSON', icon: FileCode }
	];

	// Visual editor state
	let conditionGroups = $state<ConditionGroup[]>([
		{ id: generateId(), logic: 'AND', conditions: [] }
	]);
	let actions = $state<RuleActionConfig[]>([]);

	// JSON editor state
	let jsonError = $state('');

	// Track if we're syncing to prevent circular updates
	let isSyncing = $state(false);
	let initialized = $state(false);

	// Initialize from ruleJson on mount only (untrack to prevent circular dependency)
	$effect(() => {
		if (!initialized) {
			untrack(() => {
				if (ruleJson) {
					const parsed = gruleJsonToRuleDefinition(ruleJson);
					if (parsed) {
						isSyncing = true;
						conditionGroups = parsed.conditionGroups;
						actions = parsed.actions;
						isSyncing = false;
					}
				}
			});
			initialized = true;
		}
	});

	// Sync visual editor to JSON when switching tabs
	function handleTabChange(tabId: string) {
		if (activeTab === 'visual' && tabId === 'json') {
			// Switching from visual to JSON - generate JSON
			syncVisualToJson();
		} else if (activeTab === 'json' && tabId === 'visual') {
			// Switching from JSON to visual - parse JSON
			syncJsonToVisual();
		}
		activeTab = tabId;
	}

	function syncVisualToJson() {
		const rule: RuleDefinition = {
			name,
			description,
			salience,
			conditionGroups,
			actions
		};
		ruleJson = ruleDefinitionToGruleJson(rule);
		jsonError = '';
	}

	function syncJsonToVisual() {
		const parsed = gruleJsonToRuleDefinition(ruleJson);
		if (parsed) {
			isSyncing = true;
			conditionGroups = parsed.conditionGroups;
			actions = parsed.actions;
			isSyncing = false;
			jsonError = '';
		} else {
			jsonError = 'Invalid JSON format. Manual edits may not be fully parsed into visual editor.';
		}
	}

	// Validate JSON on change
	function handleJsonChange() {
		try {
			JSON.parse(ruleJson);
			jsonError = '';
			onValidChange?.(true);
		} catch {
			jsonError = 'Invalid JSON syntax';
			onValidChange?.(false);
		}
	}

	// Update JSON when visual editor changes (but not during sync)
	$effect(() => {
		// Track these - we want to react to changes in the visual editor
		const currentTab = activeTab;
		const currentName = name;
		const currentDesc = description;
		const currentSalience = salience;
		const currentGroups = conditionGroups;
		const currentActions = actions;

		// Don't track these - they're just guards
		const syncing = untrack(() => isSyncing);
		const isInit = untrack(() => initialized);

		if (currentTab === 'visual' && !syncing && isInit) {
			// Auto-sync visual to JSON
			const rule: RuleDefinition = {
				name: currentName,
				description: currentDesc,
				salience: currentSalience,
				conditionGroups: currentGroups,
				actions: currentActions
			};
			ruleJson = ruleDefinitionToGruleJson(rule);
		}
	});
</script>

<div class="space-y-6">
	<!-- Rule Metadata -->
	<div class="grid grid-cols-2 gap-4">
		<div>
			<label for="rule-name" class="form-label">Name</label>
			<Input id="rule-name" type="text" bind:value={name} placeholder="Rule name" />
		</div>
		<div>
			<label for="rule-salience" class="form-label">Salience (Priority)</label>
			<Input
				id="rule-salience"
				type="number"
				bind:value={salience}
				placeholder="0"
				min={-100}
				max={100}
			/>
			<p class="mt-1 text-xs text-text-muted">Higher values run first (-100 to 100)</p>
		</div>
	</div>

	<div>
		<label for="rule-description" class="form-label">Description</label>
		<Textarea
			id="rule-description"
			bind:value={description}
			placeholder="What does this rule do?"
			rows={2}
		/>
	</div>

	<!-- Editor Tabs -->
	<div class="border-t border-border pt-4">
		<Tabs {tabs} bind:activeTab onchange={handleTabChange} variant="pills" />
	</div>

	<!-- Visual Editor -->
	{#if activeTab === 'visual'}
		<div class="space-y-6">
			<!-- Conditions Section -->
			<div>
				<h3 class="text-sm font-medium text-text mb-3">Conditions</h3>
				<ConditionBuilder bind:groups={conditionGroups} />
			</div>

			<!-- Actions Section -->
			<div class="border-t border-border pt-6">
				<h3 class="text-sm font-medium text-text mb-3">Actions</h3>
				<ActionBuilder bind:actions />
			</div>
		</div>
	{/if}

	<!-- JSON Editor -->
	{#if activeTab === 'json'}
		<div class="space-y-4">
			{#if jsonError}
				<Alert type="warning" title="Warning">
					<p>{jsonError}</p>
				</Alert>
			{/if}

			<div>
				<label for="rule-json" class="form-label flex items-center gap-2">
					<FileCode class="h-4 w-4" />
					Rule Definition (JSON)
				</label>
				<Textarea
					id="rule-json"
					bind:value={ruleJson}
					placeholder="Enter rule JSON..."
					rows={16}
					class="font-mono text-sm"
				/>
				<p class="mt-1 text-xs text-text-muted">
					Define your rule using Grule JSON format with "when" condition and "then" actions.
				</p>
			</div>

			<div class="flex gap-2">
				<Button type="secondary" size="sm" onclick={syncJsonToVisual}>
					Parse to Visual Editor
				</Button>
			</div>
		</div>
	{/if}
</div>
