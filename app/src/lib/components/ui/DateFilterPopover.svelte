<script lang="ts">
	import { Filter, Calendar, X } from 'lucide-svelte';

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

	let showPopover = $state(false);

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

	function selectPreset(preset: { from: string; to: string; label: string }) {
		dateRange.from = preset.from;
		dateRange.to = preset.to;
		onchange?.(dateRange);
		showPopover = false;
	}

	function togglePopover() {
		showPopover = !showPopover;
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.date-filter-container')) {
			showPopover = false;
		}
	}

	function formatDateDisplay(from: string, to: string): string {
		const fromDate = new Date(from);
		const toDate = new Date(to);
		const options: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric' };

		// Check if it matches a preset
		for (const preset of allPresets) {
			if (preset.from === from && preset.to === to) {
				return preset.label;
			}
		}

		return `${fromDate.toLocaleDateString('en-US', options)} - ${toDate.toLocaleDateString('en-US', options)}`;
	}

	const displayText = $derived(formatDateDisplay(dateRange.from, dateRange.to));
</script>

<svelte:window onclick={handleClickOutside} />

<div class="date-filter-container relative">
	<button
		type="button"
		onclick={(e) => { e.stopPropagation(); togglePopover(); }}
		class="inline-flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-lg border border-border bg-bg hover:bg-bg-secondary transition-colors text-text"
	>
		<Filter size={14} class="text-text-muted" />
		<span>{displayText}</span>
	</button>

	{#if showPopover}
		<div
			class="absolute top-full right-0 mt-2 z-50 w-80 bg-bg border border-border rounded-lg shadow-lg"
			onclick={(e) => e.stopPropagation()}
		>
			<!-- Header -->
			<div class="flex items-center justify-between px-4 py-3 border-b border-border">
				<div class="flex items-center gap-2 text-sm font-medium text-text">
					<Calendar size={16} />
					Date Range
				</div>
				<button
					type="button"
					onclick={() => showPopover = false}
					class="text-text-muted hover:text-text transition-colors"
				>
					<X size={16} />
				</button>
			</div>

			<!-- Presets -->
			<div class="p-3 border-b border-border">
				<div class="flex flex-wrap gap-2">
					{#each allPresets as preset}
						<button
							type="button"
							onclick={() => selectPreset(preset)}
							class="px-3 py-1.5 text-xs font-medium rounded-md transition-colors {dateRange.from === preset.from && dateRange.to === preset.to
								? 'bg-primary text-white'
								: 'bg-bg-secondary text-text-muted hover:bg-bg-tertiary hover:text-text'}"
						>
							{preset.label}
						</button>
					{/each}
				</div>
			</div>

			<!-- Custom Date Inputs -->
			<div class="p-4 space-y-3">
				<div class="grid grid-cols-2 gap-3">
					<div>
						<label for="from-date" class="block text-xs font-medium text-text-muted mb-1">From</label>
						<input
							type="date"
							id="from-date"
							bind:value={dateRange.from}
							onchange={handleFromChange}
							class="input w-full text-sm"
						/>
					</div>
					<div>
						<label for="to-date" class="block text-xs font-medium text-text-muted mb-1">To</label>
						<input
							type="date"
							id="to-date"
							bind:value={dateRange.to}
							onchange={handleToChange}
							class="input w-full text-sm"
						/>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
