<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	interface AttributionModel {
		key: string;
		name: string;
		description: string;
		category: 'rule-based' | 'algorithmic';
	}

	interface Props {
		selectedModel: string;
		showDescriptions?: boolean;
		allowMultiple?: boolean;
		selectedModels?: string[];
	}

	let { 
		selectedModel = $bindable(), 
		showDescriptions = true,
		allowMultiple = false,
		selectedModels = $bindable([])
	}: Props = $props();

	const dispatch = createEventDispatcher<{
		change: string;
		multiChange: string[];
	}>();

	const models: AttributionModel[] = [
		{
			key: 'first-touch',
			name: 'First-Touch',
			description: 'Credits the first touchpoint in the customer journey',
			category: 'rule-based'
		},
		{
			key: 'last-touch',
			name: 'Last-Touch',
			description: 'Credits the last touchpoint before conversion',
			category: 'rule-based'
		},
		{
			key: 'linear',
			name: 'Linear',
			description: 'Distributes credit equally across all touchpoints',
			category: 'rule-based'
		},
		{
			key: 'position-based',
			name: 'Position-Based (U-Shaped)',
			description: 'Credits first and last touchpoints more heavily',
			category: 'rule-based'
		},
		{
			key: 'data-driven',
			name: 'Data-Driven Attribution',
			description: 'Uses machine learning to determine credit distribution',
			category: 'algorithmic'
		},
		{
			key: 'shapley',
			name: 'Shapley Attribution',
			description: 'Game theory approach to fair credit distribution',
			category: 'algorithmic'
		},
		{
			key: 'bayesian',
			name: 'Bayesian Attribution',
			description: 'Probabilistic approach with continuous updating',
			category: 'algorithmic'
		},
		{
			key: 'markov',
			name: 'Markov Attribution',
			description: 'Removal effect modeling using Markov chains',
			category: 'algorithmic'
		}
	];

	function handleModelSelect(modelKey: string) {
		if (allowMultiple) {
			const newSelection = selectedModels.includes(modelKey)
				? selectedModels.filter(m => m !== modelKey)
				: [...selectedModels, modelKey];
			selectedModels = newSelection;
			dispatch('multiChange', newSelection);
		} else {
			selectedModel = modelKey;
			dispatch('change', modelKey);
		}
	}

	const ruleBasedModels = models.filter(m => m.category === 'rule-based');
	const algorithmicModels = models.filter(m => m.category === 'algorithmic');
</script>

<div class="space-y-6">
	<!-- Rule-based Models -->
	<div>
		<h3 class="text-sm font-medium text-text mb-3">Rule-Based Models</h3>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
			{#each ruleBasedModels as model}
				<div
					class="relative cursor-pointer rounded-lg border p-4 transition-all {allowMultiple
						? (selectedModels.includes(model.key) ? 'border-primary bg-primary/5 ring-2 ring-primary' : 'border-border hover:border-text-muted')
						: (selectedModel === model.key ? 'border-primary bg-primary/5 ring-2 ring-primary' : 'border-border hover:border-text-muted')}"
					onclick={() => handleModelSelect(model.key)}
					role="button"
					tabindex="0"
					onkeydown={(e) => e.key === 'Enter' && handleModelSelect(model.key)}
				>
					<div class="flex items-start">
						{#if allowMultiple}
							<input
								type="checkbox"
								id="rule-model-{model.key}"
								checked={selectedModels.includes(model.key)}
								class="mt-1 h-4 w-4 rounded border-border text-primary focus:ring-primary"
								readonly
							/>
						{:else}
							<input
								type="radio"
								id="rule-model-{model.key}"
								name="attribution-model"
								value={model.key}
								checked={selectedModel === model.key}
								class="mt-1 h-4 w-4 border-border text-primary focus:ring-primary"
								readonly
							/>
						{/if}
						<div class="ml-3">
							<label for="rule-model-{model.key}" class="cursor-pointer text-sm font-medium text-text">
								{model.name}
							</label>
							{#if showDescriptions}
								<p class="mt-1 text-sm text-text-muted">{model.description}</p>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>

	<!-- Algorithmic Models -->
	<div>
		<h3 class="text-sm font-medium text-text mb-3">Algorithmic Models</h3>
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
			{#each algorithmicModels as model}
				<div
					class="relative cursor-pointer rounded-lg border p-4 transition-all {allowMultiple
						? (selectedModels.includes(model.key) ? 'border-primary bg-primary/5 ring-2 ring-primary' : 'border-border hover:border-text-muted')
						: (selectedModel === model.key ? 'border-primary bg-primary/5 ring-2 ring-primary' : 'border-border hover:border-text-muted')}"
					onclick={() => handleModelSelect(model.key)}
					role="button"
					tabindex="0"
					onkeydown={(e) => e.key === 'Enter' && handleModelSelect(model.key)}
				>
					<div class="flex items-start">
						{#if allowMultiple}
							<input
								type="checkbox"
								id="algo-model-{model.key}"
								checked={selectedModels.includes(model.key)}
								class="mt-1 h-4 w-4 rounded border-border text-primary focus:ring-primary"
								readonly
							/>
						{:else}
							<input
								type="radio"
								id="algo-model-{model.key}"
								name="attribution-model"
								value={model.key}
								checked={selectedModel === model.key}
								class="mt-1 h-4 w-4 border-border text-primary focus:ring-primary"
								readonly
							/>
						{/if}
						<div class="ml-3">
							<label for="algo-model-{model.key}" class="cursor-pointer text-sm font-medium text-text">
								{model.name}
							</label>
							{#if showDescriptions}
								<p class="mt-1 text-sm text-text-muted">{model.description}</p>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>

