<script lang="ts">
	import { Button, Select, Input } from '$lib/components/ui';
	import { Plus, Trash2, GripVertical, Zap } from 'lucide-svelte';
	import {
		type RuleActionConfig,
		ACTION_TYPES,
		getActionType,
		getActionsByCategory,
		generateId
	} from './types';

	let {
		actions = $bindable<RuleActionConfig[]>([])
	}: {
		actions: RuleActionConfig[];
	} = $props();

	function createEmptyAction(): RuleActionConfig {
		return {
			id: generateId(),
			type: '',
			params: {}
		};
	}

	function addAction() {
		actions = [...actions, createEmptyAction()];
	}

	function removeAction(index: number) {
		actions = actions.filter((_, i) => i !== index);
	}

	function handleActionTypeChange(index: number, newType: string) {
		const actionDef = getActionType(newType);
		// Reset params when type changes
		actions[index].type = newType;
		actions[index].params = {};

		// Set default values for required params
		if (actionDef) {
			for (const param of actionDef.params) {
				if (param.required) {
					actions[index].params[param.key] = '';
				}
			}
		}
	}

	const actionsByCategory = getActionsByCategory();
	const categoryLabels: Record<string, string> = {
		communication: 'Communication',
		ticket: 'Tickets',
		customer: 'Customer',
		subscription: 'Subscription',
		system: 'System'
	};
</script>

<div class="space-y-4">
	<div class="flex items-center gap-2 mb-2">
		<Zap class="h-4 w-4 text-text-muted" />
		<span class="text-sm font-medium text-text">Then execute these actions</span>
	</div>

	{#each actions as action, index (action.id)}
		{@const actionDef = getActionType(action.type)}
		<div class="border border-border rounded-lg p-4 bg-bg-secondary/50">
			<div class="flex items-start gap-2">
				<GripVertical class="h-4 w-4 text-text-muted flex-shrink-0 cursor-grab mt-2" />

				<div class="flex-1 space-y-3">
					<!-- Action Type Select -->
					<div class="flex items-center gap-2">
						<Select
							value={action.type}
							size="sm"
							class="flex-1"
							onchange={() => {}}
						>
							<option value="">Select action...</option>
							{#each Object.entries(actionsByCategory) as [category, categoryActions]}
								<optgroup label={categoryLabels[category] || category}>
									{#each categoryActions as act}
										<option value={act.type}>{act.name}</option>
									{/each}
								</optgroup>
							{/each}
						</Select>

						<Button type="ghost" size="sm" onclick={() => removeAction(index)}>
							<Trash2 class="h-4 w-4 text-text-muted hover:text-red-500" />
						</Button>
					</div>

					<!-- Hack to handle select change properly -->
					<select
						class="hidden"
						bind:value={action.type}
						onchange={(e) => handleActionTypeChange(index, (e.target as HTMLSelectElement).value)}
					>
						<option value="">Select action...</option>
						{#each ACTION_TYPES as act}
							<option value={act.type}>{act.name}</option>
						{/each}
					</select>

					<!-- Action Description -->
					{#if actionDef}
						<p class="text-xs text-text-muted">{actionDef.description}</p>

						<!-- Action Parameters -->
						{#if actionDef.params.length > 0}
							<div class="grid grid-cols-2 gap-3 pt-2 border-t border-border/50">
								{#each actionDef.params as param}
									<div>
										<label class="block text-xs font-medium text-text-muted mb-1">
											{param.name}
											{#if param.required}
												<span class="text-red-500">*</span>
											{/if}
										</label>

										{#if param.type === 'select' && param.options}
											<Select
												bind:value={action.params[param.key]}
												size="sm"
												options={param.options}
											>
												<option value="">Select...</option>
												{#each param.options as opt}
													<option value={opt.value}>{opt.label}</option>
												{/each}
											</Select>
										{:else if param.type === 'number'}
											<Input
												type="number"
												bind:value={action.params[param.key]}
												size="sm"
												placeholder={param.placeholder}
											/>
										{:else}
											<Input
												type="text"
												bind:value={action.params[param.key]}
												size="sm"
												placeholder={param.placeholder}
											/>
										{/if}

										{#if param.description}
											<p class="text-xs text-text-muted mt-0.5">{param.description}</p>
										{/if}
									</div>
								{/each}
							</div>
						{/if}
					{/if}
				</div>
			</div>
		</div>
	{/each}

	<!-- Add Action Button -->
	<Button type="secondary" size="sm" onclick={addAction}>
		<Plus class="h-4 w-4 mr-1" />
		Add action
	</Button>

	{#if actions.length === 0}
		<p class="text-sm text-text-muted text-center py-4">
			No actions configured. Add at least one action for this rule.
		</p>
	{/if}
</div>
