<script lang="ts">
	import { Badge, Card, Tooltip } from '$lib/components/ui';
	import { ArrowDown, Zap, Filter, ChevronRight } from 'lucide-svelte';
	import {
		gruleJsonToRuleDefinition,
		conditionGroupsToText,
		actionsToText,
		getFactFieldByPath,
		getActionType,
		type ConditionGroup,
		type RuleActionConfig
	} from './types';
	import type { RuleInfo } from '$lib/api';

	let {
		rules,
		onRuleClick
	}: {
		rules: RuleInfo[];
		onRuleClick?: (rule: RuleInfo) => void;
	} = $props();

	// Sort rules by salience (higher first) and parse them
	let sortedRules = $derived(
		[...rules]
			.filter((r) => r.enabled)
			.sort((a, b) => b.salience - a.salience)
			.map((rule) => {
				const parsed = gruleJsonToRuleDefinition(rule.rule_json);
				return {
					...rule,
					parsed,
					conditionSummary: parsed ? conditionGroupsToText(parsed.conditionGroups) : 'Unable to parse',
					actionSummary: parsed ? actionsToText(parsed.actions) : 'Unable to parse',
					factsUsed: parsed ? extractFactsUsed(parsed.conditionGroups) : [],
					actionsUsed: parsed ? extractActionsUsed(parsed.actions) : []
				};
			})
	);

	// Group by category
	let rulesByCategory = $derived(
		sortedRules.reduce(
			(acc, rule) => {
				if (!acc[rule.category]) {
					acc[rule.category] = [];
				}
				acc[rule.category].push(rule);
				return acc;
			},
			{} as Record<string, typeof sortedRules>
		)
	);

	// Extract unique facts used in conditions
	function extractFactsUsed(groups: ConditionGroup[]): string[] {
		const facts = new Set<string>();
		for (const group of groups) {
			for (const condition of group.conditions) {
				if (condition.field) {
					// Get category from path (e.g., "Ticket.Status" -> "Ticket")
					const category = condition.field.split('.')[0];
					facts.add(category);
				}
			}
		}
		return Array.from(facts);
	}

	// Extract unique actions used
	function extractActionsUsed(actions: RuleActionConfig[]): string[] {
		return actions.map((a) => {
			const actionDef = getActionType(a.type);
			return actionDef?.name || a.type;
		});
	}

	function getCategoryColor(category: string): string {
		switch (category) {
			case 'support':
				return 'border-blue-500 bg-blue-50 dark:bg-blue-900/20';
			case 'billing':
				return 'border-green-500 bg-green-50 dark:bg-green-900/20';
			case 'email':
				return 'border-yellow-500 bg-yellow-50 dark:bg-yellow-900/20';
			case 'customer':
				return 'border-purple-500 bg-purple-50 dark:bg-purple-900/20';
			default:
				return 'border-gray-500 bg-gray-50 dark:bg-gray-900/20';
		}
	}

	function getFactColor(fact: string): string {
		switch (fact) {
			case 'Ticket':
				return 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300';
			case 'Customer':
				return 'bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300';
			case 'Contact':
				return 'bg-pink-100 text-pink-700 dark:bg-pink-900 dark:text-pink-300';
			case 'Subscription':
				return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300';
			case 'Payment':
				return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300';
			case 'Invoice':
				return 'bg-orange-100 text-orange-700 dark:bg-orange-900 dark:text-orange-300';
			default:
				return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
		}
	}

	const categoryLabels: Record<string, string> = {
		support: 'Support Rules',
		billing: 'Billing Rules',
		email: 'Email Rules',
		customer: 'Customer Rules'
	};
</script>

<div class="space-y-8">
	<!-- Legend -->
	<div class="flex flex-wrap items-center gap-4 text-xs text-text-muted">
		<span class="font-medium">Facts:</span>
		{#each ['Ticket', 'Customer', 'Contact', 'Subscription', 'Payment'] as fact}
			<span class="px-2 py-0.5 rounded {getFactColor(fact)}">{fact}</span>
		{/each}
	</div>

	<!-- Rules by Category -->
	{#each Object.entries(rulesByCategory) as [category, categoryRules]}
		<div>
			<h3 class="text-sm font-medium text-text mb-3 flex items-center gap-2">
				<span class="w-3 h-3 rounded-full {getCategoryColor(category).split(' ')[0].replace('border-', 'bg-')}"></span>
				{categoryLabels[category] || category}
				<Badge variant="default">{categoryRules.length}</Badge>
			</h3>

			<!-- Flow visualization -->
			<div class="relative">
				{#each categoryRules as rule, index}
					<div class="relative">
						<!-- Connector line -->
						{#if index > 0}
							<div class="absolute left-6 -top-4 w-0.5 h-4 bg-border"></div>
						{/if}

						<!-- Rule node -->
						<button
							type="button"
							class="w-full text-left p-4 border-l-4 rounded-lg transition-all hover:shadow-md {getCategoryColor(
								category
							)} mb-2"
							onclick={() => onRuleClick?.(rule)}
						>
							<div class="flex items-start justify-between gap-4">
								<div class="flex-1 min-w-0">
									<!-- Header -->
									<div class="flex items-center gap-2 mb-2">
										<span class="text-xs font-mono text-text-muted">#{rule.salience}</span>
										<span class="font-medium text-text">{rule.name}</span>
									</div>

									<!-- Facts used -->
									{#if rule.factsUsed.length > 0}
										<div class="flex items-center gap-2 mb-2">
											<Filter class="h-3 w-3 text-text-muted" />
											<div class="flex flex-wrap gap-1">
												{#each rule.factsUsed as fact}
													<span class="text-xs px-1.5 py-0.5 rounded {getFactColor(fact)}">{fact}</span>
												{/each}
											</div>
										</div>
									{/if}

									<!-- Condition summary -->
									<div class="text-xs text-text-muted mb-2 line-clamp-1">
										<span class="font-medium">When:</span> {rule.conditionSummary}
									</div>

									<!-- Actions -->
									{#if rule.actionsUsed.length > 0}
										<div class="flex items-center gap-2">
											<Zap class="h-3 w-3 text-primary" />
											<div class="flex flex-wrap gap-1">
												{#each rule.actionsUsed as action}
													<span class="text-xs px-1.5 py-0.5 rounded bg-primary/10 text-primary"
														>{action}</span
													>
												{/each}
											</div>
										</div>
									{/if}
								</div>

								<ChevronRight class="h-4 w-4 text-text-muted flex-shrink-0" />
							</div>
						</button>

						<!-- Arrow to next rule -->
						{#if index < categoryRules.length - 1}
							<div class="flex items-center justify-center py-1">
								<ArrowDown class="h-4 w-4 text-text-muted" />
							</div>
						{/if}
					</div>
				{/each}
			</div>
		</div>
	{/each}

	{#if sortedRules.length === 0}
		<div class="text-center py-8 text-text-muted">
			<p>No enabled rules to display</p>
		</div>
	{/if}

	<!-- Execution Order Note -->
	{#if sortedRules.length > 0}
		<Card>
			<div class="text-sm text-text-muted">
				<h4 class="font-medium text-text mb-2">How Rules Execute</h4>
				<ul class="space-y-1 list-disc list-inside">
					<li>Rules execute in <strong>salience order</strong> (higher numbers first)</li>
					<li>Multiple rules can fire for the same event</li>
					<li>Actions from all matching rules are executed</li>
					<li>Rules checking the same facts may interact via forward chaining</li>
				</ul>
			</div>
		</Card>
	{/if}
</div>
