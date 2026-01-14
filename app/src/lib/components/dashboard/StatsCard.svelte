<script lang="ts">
	let {
		icon = '',
		label = '',
		value = 0,
		trend = null,
		trendValue = '',
		bgColor = 'bg-white',
		iconColor = 'text-blue-600'
	}: {
		icon?: string;
		label?: string;
		value?: string | number;
		trend?: 'up' | 'down' | 'neutral' | null;
		trendValue?: string;
		bgColor?: string;
		iconColor?: string;
	} = $props();
</script>

<div class="{bgColor} rounded-lg shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow">
	<div class="flex items-center justify-between">
		<div class="flex-1">
			<p class="text-sm font-medium text-gray-600 mb-1">{label}</p>
			<p class="text-3xl font-bold text-gray-900">{value}</p>
			{#if trend && trendValue}
				<div class="flex items-center mt-2">
					{#if trend === 'up'}
						<svg class="w-4 h-4 text-green-600 mr-1" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
						</svg>
						<span class="text-sm font-medium text-green-600">{trendValue}</span>
					{:else if trend === 'down'}
						<svg class="w-4 h-4 text-red-600 mr-1" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z" clip-rule="evenodd" />
						</svg>
						<span class="text-sm font-medium text-red-600">{trendValue}</span>
					{:else}
						<span class="text-sm font-medium text-gray-600">{trendValue}</span>
					{/if}
				</div>
			{/if}
		</div>
		{#if icon}
			<div class="{iconColor} opacity-80">
				{@html icon}
			</div>
		{/if}
	</div>
</div>
