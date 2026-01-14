<script lang="ts">
	import { Badge, Button, Tooltip } from '$lib/components/ui';
	import { Pencil, Trash2, Power, PowerOff, ChevronDown, ChevronRight, Zap, Play } from 'lucide-svelte';
	import {
		gruleJsonToRuleDefinition,
		conditionGroupsToText,
		actionsToText
	} from './types';
	import type { RuleInfo } from '$lib/api';

	let {
		rule,
		onEdit,
		onToggle,
		onDelete,
		onTest
	}: {
		rule: RuleInfo;
		onEdit?: (rule: RuleInfo) => void;
		onToggle?: (rule: RuleInfo) => void;
		onDelete?: (rule: RuleInfo) => void;
		onTest?: (rule: RuleInfo) => void;
	} = $props();

	let expanded = $state(false);

	// Parse rule JSON for display
	let parsedRule = $derived.by(() => {
		try {
			return gruleJsonToRuleDefinition(rule.rule_json);
		} catch {
			return null;
		}
	});

	let conditionSummary = $derived(
		parsedRule ? conditionGroupsToText(parsedRule.conditionGroups) : 'Unable to parse'
	);

	let actionsSummary = $derived(
		parsedRule ? actionsToText(parsedRule.actions) : 'Unable to parse'
	);

	function getCategoryVariant(category: string): 'info' | 'success' | 'warning' | 'default' {
		switch (category) {
			case 'support':
				return 'info';
			case 'billing':
				return 'success';
			case 'email':
				return 'warning';
			default:
				return 'default';
		}
	}
</script>

<div
	class="border border-border rounded-lg transition-colors hover:border-border-hover {!rule.enabled
		? 'opacity-60'
		: ''}"
>
	<!-- Header -->
	<div class="p-4">
		<div class="flex items-start justify-between gap-4">
			<div class="flex-1 min-w-0">
				<!-- Title Row -->
				<div class="flex items-center gap-2 flex-wrap">
					<button
						type="button"
						class="flex items-center gap-1 text-left"
						onclick={() => (expanded = !expanded)}
					>
						{#if expanded}
							<ChevronDown class="h-4 w-4 text-text-muted flex-shrink-0" />
						{:else}
							<ChevronRight class="h-4 w-4 text-text-muted flex-shrink-0" />
						{/if}
						<h3 class="text-sm font-medium text-text truncate">{rule.name}</h3>
					</button>
					<Badge variant={getCategoryVariant(rule.category)}>{rule.category}</Badge>
					{#if !rule.enabled}
						<Badge variant="default">Disabled</Badge>
					{/if}
					{#if rule.validation_errors && rule.validation_errors.length > 0}
						<Tooltip content={rule.validation_errors.join(', ')}>
							<Badge variant="error">{rule.validation_errors.length} error(s)</Badge>
						</Tooltip>
					{/if}
				</div>

				<!-- Condition Preview -->
				<div class="mt-2 text-sm text-text-muted">
					<span class="font-medium text-text">When:</span>
					<span class="ml-1 line-clamp-1">{conditionSummary}</span>
				</div>

				<!-- Actions Preview -->
				<div class="mt-1 text-sm text-text-muted">
					<span class="font-medium text-text">Then:</span>
					<span class="ml-1 line-clamp-1">{actionsSummary}</span>
				</div>

				<!-- Meta -->
				<div class="mt-2 flex items-center gap-4 text-xs text-text-muted">
					<span>Priority: {rule.salience}</span>
					{#if rule.entity_type}
						<span>{rule.entity_type}: {rule.entity_id}</span>
					{/if}
				</div>
			</div>

			<!-- Action Buttons -->
			<div class="flex items-center gap-1 flex-shrink-0">
				<Tooltip content="Test rule">
					<Button type="ghost" size="icon" onclick={() => onTest?.(rule)}>
						<Play class="h-4 w-4 text-primary" />
					</Button>
				</Tooltip>
				<Tooltip content={rule.enabled ? 'Disable rule' : 'Enable rule'}>
					<Button type="ghost" size="icon" onclick={() => onToggle?.(rule)}>
						{#if rule.enabled}
							<Power class="h-4 w-4 text-green-500" />
						{:else}
							<PowerOff class="h-4 w-4 text-text-muted" />
						{/if}
					</Button>
				</Tooltip>
				<Tooltip content="Edit rule">
					<Button type="ghost" size="icon" onclick={() => onEdit?.(rule)}>
						<Pencil class="h-4 w-4" />
					</Button>
				</Tooltip>
				<Tooltip content="Delete rule">
					<Button type="ghost" size="icon" onclick={() => onDelete?.(rule)}>
						<Trash2 class="h-4 w-4 text-red-500" />
					</Button>
				</Tooltip>
			</div>
		</div>
	</div>

	<!-- Expanded Details -->
	{#if expanded}
		<div class="px-4 pb-4 pt-0 border-t border-border">
			{#if rule.description}
				<p class="text-sm text-text-muted mt-3 mb-4">{rule.description}</p>
			{/if}

			<!-- Full Conditions -->
			<div class="space-y-3">
				<div>
					<h4 class="text-xs font-medium text-text uppercase tracking-wider mb-2">Conditions</h4>
					{#if parsedRule && parsedRule.conditionGroups.length > 0}
						<div class="space-y-2">
							{#each parsedRule.conditionGroups as group, i}
								<div class="text-sm bg-bg-secondary rounded p-2">
									{#if i > 0}
										<span class="text-primary font-medium">AND </span>
									{/if}
									{#if group.conditions.length === 0}
										<span class="text-text-muted italic">Always matches</span>
									{:else}
										{#each group.conditions as condition, j}
											{#if j > 0}
												<span class="text-primary font-medium">
													{group.logic === 'AND' ? ' AND ' : ' OR '}
												</span>
											{/if}
											<span class="text-text">{condition.field || 'No field'}</span>
											<span class="text-text-muted"> {condition.operator} </span>
											<span class="text-accent">"{condition.value}"</span>
										{/each}
									{/if}
								</div>
							{/each}
						</div>
					{:else}
						<p class="text-sm text-text-muted italic">No conditions (always matches)</p>
					{/if}
				</div>

				<!-- Full Actions -->
				<div>
					<h4 class="text-xs font-medium text-text uppercase tracking-wider mb-2">Actions</h4>
					{#if parsedRule && parsedRule.actions.length > 0}
						<div class="space-y-2">
							{#each parsedRule.actions as action}
								<div class="flex items-start gap-2 text-sm bg-bg-secondary rounded p-2">
									<Zap class="h-4 w-4 text-primary flex-shrink-0 mt-0.5" />
									<div>
										<span class="font-medium text-text">{action.type}</span>
										{#if Object.keys(action.params).length > 0}
											<div class="text-text-muted mt-1">
												{#each Object.entries(action.params) as [key, value]}
													<div class="text-xs">
														<span class="text-text-muted">{key}:</span>
														<span class="text-text ml-1">{value}</span>
													</div>
												{/each}
											</div>
										{/if}
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<p class="text-sm text-text-muted italic">No actions configured</p>
					{/if}
				</div>
			</div>
		</div>
	{/if}
</div>
