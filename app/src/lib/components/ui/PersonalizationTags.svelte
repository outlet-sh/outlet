<script lang="ts">
	interface Variable {
		name: string;
		label: string;
		required?: boolean;
	}

	let {
		variables = [],
		insertVariable = null,
		title = 'Personalization Tags',
		noBorder = false
	}: {
		variables: Variable[];
		insertVariable?: ((variable: string) => void) | null;
		title?: string;
		noBorder?: boolean;
	} = $props();
</script>

<div class={noBorder ? '' : 'border-t border-base-300 pt-6'}>
	<h4 class="text-sm font-medium text-base-content mb-3">{title}</h4>
	<p class="text-xs text-base-content/60 mb-3">
		Click to insert into email:
	</p>
	<div class="space-y-2">
		{#each variables as variable}
			<button
				type="button"
				onclick={() => insertVariable?.(`{{${variable.name}}}`)}
				class="w-full text-left p-2 rounded text-sm transition-colors {variable.required ? 'bg-warning/20 hover:bg-warning/30' : 'bg-base-200 hover:bg-base-300'}"
			>
				<code class="font-mono text-xs">{`{{${variable.name}}}`}</code>
				<span class="text-base-content/60 ml-1 text-xs">- {variable.label}</span>
				{#if variable.required}
					<span class="text-warning text-[10px] ml-1">(required)</span>
				{/if}
			</button>
		{/each}
	</div>
</div>
