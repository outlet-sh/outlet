<script lang="ts">
	let { activities = [] }: {
		activities?: Array<{
			id: string;
			type: 'lead_assigned' | 'status_change' | 'meeting_scheduled' | 'meeting_completed' | 'lead_created';
			title: string;
			description: string;
			timestamp: Date;
			actor?: string;
		}>
	} = $props();

	function getActivityIcon(type: string) {
		switch (type) {
			case 'lead_assigned':
				return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
				</svg>`;
			case 'status_change':
				return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>`;
			case 'meeting_scheduled':
				return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
				</svg>`;
			case 'meeting_completed':
				return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
				</svg>`;
			case 'lead_created':
				return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
				</svg>`;
			default:
				return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>`;
		}
	}

	function getActivityColor(type: string) {
		switch (type) {
			case 'lead_assigned':
				return 'bg-blue-100 text-blue-600';
			case 'status_change':
				return 'bg-green-100 text-green-600';
			case 'meeting_scheduled':
				return 'bg-purple-100 text-purple-600';
			case 'meeting_completed':
				return 'bg-teal-100 text-teal-600';
			case 'lead_created':
				return 'bg-orange-100 text-orange-600';
			default:
				return 'bg-gray-100 text-gray-600';
		}
	}

	function formatTimestamp(timestamp: Date) {
		const now = new Date();
		const diff = now.getTime() - new Date(timestamp).getTime();
		const minutes = Math.floor(diff / 60000);
		const hours = Math.floor(diff / 3600000);
		const days = Math.floor(diff / 86400000);

		if (minutes < 1) return 'Just now';
		if (minutes < 60) return `${minutes}m ago`;
		if (hours < 24) return `${hours}h ago`;
		if (days < 7) return `${days}d ago`;
		return new Date(timestamp).toLocaleDateString();
	}
</script>

<div class="bg-white rounded-lg shadow-sm border border-gray-200">
	<div class="p-6 border-b border-gray-200">
		<h3 class="text-lg font-semibold text-gray-900">Recent Activity</h3>
	</div>
	<div class="divide-y divide-gray-200">
		{#if activities.length === 0}
			<div class="p-6 text-center text-gray-500">
				<p>No recent activity</p>
			</div>
		{:else}
			{#each activities as activity (activity.id)}
				<div class="p-4 hover:bg-gray-50 transition-colors">
					<div class="flex items-start space-x-3">
						<div class="flex-shrink-0">
							<div class="w-10 h-10 rounded-full {getActivityColor(activity.type)} flex items-center justify-center">
								{@html getActivityIcon(activity.type)}
							</div>
						</div>
						<div class="flex-1 min-w-0">
							<p class="text-sm font-medium text-gray-900">{activity.title}</p>
							<p class="text-sm text-gray-600">{activity.description}</p>
							{#if activity.actor}
								<p class="text-xs text-gray-500 mt-1">by {activity.actor}</p>
							{/if}
						</div>
						<div class="flex-shrink-0">
							<p class="text-xs text-gray-500">{formatTimestamp(activity.timestamp)}</p>
						</div>
					</div>
				</div>
			{/each}
		{/if}
	</div>
</div>
