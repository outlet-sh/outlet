<script lang="ts">
	interface DateRange {
		from: string;
		to: string;
	}

	interface Props {
		dateRange: DateRange;
		presets?: { label: string; from: string; to: string }[];
		onchange?: (range: DateRange) => void;
	}

	let { dateRange = $bindable(), presets = [], onchange }: Props = $props();

	const defaultPresets = [
		{
			label: 'Last 7 days',
			from: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
			to: new Date().toISOString().split('T')[0]
		},
		{
			label: 'Last 30 days',
			from: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
			to: new Date().toISOString().split('T')[0]
		},
		{
			label: 'Last 90 days',
			from: new Date(Date.now() - 90 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
			to: new Date().toISOString().split('T')[0]
		},
		{
			label: 'This month',
			from: new Date(new Date().getFullYear(), new Date().getMonth(), 1).toISOString().split('T')[0],
			to: new Date().toISOString().split('T')[0]
		}
	];

	const allPresets = presets.length > 0 ? presets : defaultPresets;

	function handleFromChange(event: Event) {
		const input = event.target as HTMLInputElement;
		dateRange.from = input.value;
		onchange?.(dateRange);
	}

	function handleToChange(event: Event) {
		const input = event.target as HTMLInputElement;
		dateRange.to = input.value;
		onchange?.(dateRange);
	}

	function selectPreset(preset: { from: string; to: string }) {
		dateRange.from = preset.from;
		dateRange.to = preset.to;
		onchange?.(dateRange);
	}
</script>

<div class="space-y-4">
	<div class="flex flex-wrap gap-2">
		{#each allPresets as preset}
			<button
				type="button"
				onclick={() => selectPreset(preset)}
				class="btn-secondary"
			>
				{preset.label}
			</button>
		{/each}
	</div>

	<div class="flex space-x-4">
		<div class="flex-1">
			<label for="from-date" class="block text-sm font-medium text-text-muted">From</label>
			<input
				type="date"
				id="from-date"
				bind:value={dateRange.from}
				onchange={handleFromChange}
				class="input mt-1 block w-full"
			/>
		</div>
		<div class="flex-1">
			<label for="to-date" class="block text-sm font-medium text-text-muted">To</label>
			<input
				type="date"
				id="to-date"
				bind:value={dateRange.to}
				onchange={handleToChange}
				class="input mt-1 block w-full"
			/>
		</div>
	</div>
</div>

