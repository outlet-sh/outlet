<script lang="ts">
	import { Button, Select, Input } from '$lib/components/ui';
	import { Plus, Trash2, GripVertical } from 'lucide-svelte';
	import {
		type ConditionGroup,
		type RuleCondition,
		FACT_CATEGORIES,
		getFactFieldByPath,
		getOperatorsForType,
		generateId
	} from './types';

	let {
		groups = $bindable<ConditionGroup[]>([])
	}: {
		groups: ConditionGroup[];
	} = $props();

	// Ensure at least one group exists
	$effect(() => {
		if (groups.length === 0) {
			groups = [createEmptyGroup()];
		}
	});

	function createEmptyGroup(): ConditionGroup {
		return {
			id: generateId(),
			logic: 'AND',
			conditions: []
		};
	}

	function createEmptyCondition(): RuleCondition {
		return {
			id: generateId(),
			field: '',
			operator: '==',
			value: ''
		};
	}

	function addCondition(groupIndex: number) {
		groups[groupIndex].conditions = [...groups[groupIndex].conditions, createEmptyCondition()];
	}

	function removeCondition(groupIndex: number, conditionIndex: number) {
		groups[groupIndex].conditions = groups[groupIndex].conditions.filter(
			(_, i) => i !== conditionIndex
		);
	}

	function addGroup() {
		groups = [...groups, createEmptyGroup()];
	}

	function removeGroup(groupIndex: number) {
		if (groups.length > 1) {
			groups = groups.filter((_, i) => i !== groupIndex);
		}
	}

	function getFieldOptions() {
		const options: { value: string; label: string }[] = [];
		for (const category of FACT_CATEGORIES) {
			for (const field of category.fields) {
				options.push({
					value: field.path,
					label: `${category.name} › ${field.name}`
				});
			}
		}
		return options;
	}

	function getOperatorOptions(fieldPath: string) {
		const field = getFactFieldByPath(fieldPath);
		const type = field?.type || 'string';
		return getOperatorsForType(type).map((op) => ({
			value: op.value,
			label: op.label
		}));
	}

	function getValueOptions(fieldPath: string) {
		const field = getFactFieldByPath(fieldPath);
		return field?.options || null;
	}

	function getFieldType(fieldPath: string): 'string' | 'number' | 'boolean' {
		const field = getFactFieldByPath(fieldPath);
		if (field?.type === 'number') return 'number';
		if (field?.type === 'boolean') return 'boolean';
		return 'string';
	}

	const fieldOptions = getFieldOptions();
</script>

<div class="space-y-4">
	{#each groups as group, groupIndex (group.id)}
		<div class="border border-border rounded-lg p-4 bg-bg-secondary/50">
			<!-- Group Header -->
			<div class="flex items-center justify-between mb-3">
				<div class="flex items-center gap-2">
					<span class="text-sm font-medium text-text">
						{#if groupIndex === 0}
							When
						{:else}
							AND when
						{/if}
					</span>
					{#if group.conditions.length > 1}
						<Select
							bind:value={group.logic}
							size="sm"
							class="w-20"
							options={[
								{ value: 'AND', label: 'ALL' },
								{ value: 'OR', label: 'ANY' }
							]}
						/>
						<span class="text-xs text-text-muted">of these are true</span>
					{/if}
				</div>
				{#if groups.length > 1}
					<Button type="ghost" size="sm" onclick={() => removeGroup(groupIndex)}>
						<Trash2 class="h-4 w-4 text-red-500" />
					</Button>
				{/if}
			</div>

			<!-- Conditions -->
			<div class="space-y-2">
				{#each group.conditions as condition, conditionIndex (condition.id)}
					{@const valueOptions = getValueOptions(condition.field)}
					{@const fieldType = getFieldType(condition.field)}
					<div class="flex items-center gap-2">
						<GripVertical class="h-4 w-4 text-text-muted flex-shrink-0 cursor-grab" />

						<!-- Field Select -->
						<Select
							bind:value={condition.field}
							size="sm"
							class="flex-1 min-w-[180px]"
							options={fieldOptions}
						>
							<option value="">Select field...</option>
							{#each FACT_CATEGORIES as category}
								<optgroup label={category.name}>
									{#each category.fields as field}
										<option value={field.path}>{field.name}</option>
									{/each}
								</optgroup>
							{/each}
						</Select>

						<!-- Operator Select -->
						<Select
							bind:value={condition.operator}
							size="sm"
							class="w-40"
							options={getOperatorOptions(condition.field)}
						/>

						<!-- Value Input -->
						{#if valueOptions}
							<Select
								bind:value={condition.value}
								size="sm"
								class="flex-1 min-w-[120px]"
								options={valueOptions}
							>
								<option value="">Select value...</option>
								{#each valueOptions as opt}
									<option value={opt.value}>{opt.label}</option>
								{/each}
							</Select>
						{:else if fieldType === 'boolean'}
							<Select
								bind:value={condition.value}
								size="sm"
								class="flex-1 min-w-[120px]"
								options={[
									{ value: 'true', label: 'True' },
									{ value: 'false', label: 'False' }
								]}
							/>
						{:else if fieldType === 'number'}
							<Input
								type="number"
								bind:value={condition.value}
								size="sm"
								class="flex-1 min-w-[120px]"
								placeholder="Enter value..."
							/>
						{:else}
							<Input
								type="text"
								bind:value={condition.value}
								size="sm"
								class="flex-1 min-w-[120px]"
								placeholder="Enter value..."
							/>
						{/if}

						<!-- Remove Button -->
						<Button
							type="ghost"
							size="sm"
							onclick={() => removeCondition(groupIndex, conditionIndex)}
						>
							<Trash2 class="h-4 w-4 text-text-muted hover:text-red-500" />
						</Button>
					</div>
				{/each}

				<!-- Add Condition Button -->
				<Button type="ghost" size="sm" onclick={() => addCondition(groupIndex)} class="mt-2">
					<Plus class="h-4 w-4 mr-1" />
					Add condition
				</Button>
			</div>
		</div>
	{/each}

	<!-- Add Group Button -->
	<Button type="secondary" size="sm" onclick={addGroup}>
		<Plus class="h-4 w-4 mr-1" />
		Add condition group (AND)
	</Button>
</div>
