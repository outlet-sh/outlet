<script lang="ts">
	import StatsCard from './StatsCard.svelte';
	import RecentActivity from './RecentActivity.svelte';

	let { data }: {
		data: {
			stats: {
				totalLeads: number;
				qualifiedLeads: number;
				meetingsScheduled: number;
				upcomingMeetings: number;
			};
			leads: Array<any>;
			upcomingMeetings: Array<any>;
			todaySchedule: Array<any>;
			recentActivity: Array<any>;
		}
	} = $props();

	const userIcon = `<svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
		<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
	</svg>`;

	const checkIcon = `<svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
		<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
	</svg>`;

	const calendarIcon = `<svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
		<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
	</svg>`;

	const clockIcon = `<svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
		<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
	</svg>`;

	function formatDateTime(date: Date | string) {
		const d = new Date(date);
		return d.toLocaleString('en-US', {
			month: 'short',
			day: 'numeric',
			hour: 'numeric',
			minute: '2-digit'
		});
	}

	function formatTime(date: Date | string) {
		const d = new Date(date);
		return d.toLocaleTimeString('en-US', {
			hour: 'numeric',
			minute: '2-digit'
		});
	}

	const qualificationRate = $derived(data.stats.totalLeads > 0
		? Math.round((data.stats.qualifiedLeads / data.stats.totalLeads) * 100)
		: 0);
</script>

<div class="space-y-6">
	<!-- Stats Grid -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
		<StatsCard
			icon={userIcon}
			iconColor="text-blue-600"
			label="Total Assigned Leads"
			value={data.stats.totalLeads}
		/>
		<StatsCard
			icon={checkIcon}
			iconColor="text-green-600"
			label="Qualified Leads"
			value={data.stats.qualifiedLeads}
			trend={qualificationRate >= 50 ? 'up' : 'neutral'}
			trendValue="{qualificationRate}% rate"
		/>
		<StatsCard
			icon={calendarIcon}
			iconColor="text-purple-600"
			label="Meetings Scheduled"
			value={data.stats.meetingsScheduled}
		/>
		<StatsCard
			icon={clockIcon}
			iconColor="text-orange-600"
			label="Upcoming Meetings"
			value={data.stats.upcomingMeetings}
		/>
	</div>

	<!-- Today's Schedule -->
	<div class="bg-white rounded-lg shadow-sm border border-gray-200">
		<div class="p-6 border-b border-gray-200">
			<h3 class="text-lg font-semibold text-gray-900">Today's Schedule</h3>
		</div>
		<div class="p-6">
			{#if data.todaySchedule.length === 0}
				<p class="text-gray-500 text-center py-8">No meetings scheduled for today</p>
			{:else}
				<div class="space-y-4">
					{#each data.todaySchedule as meeting}
						<div class="flex items-start space-x-4 p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
							<div class="flex-shrink-0">
								<div class="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center">
									<svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
									</svg>
								</div>
							</div>
							<div class="flex-1 min-w-0">
								<p class="text-sm font-semibold text-gray-900">{meeting.attendee_name}</p>
								<p class="text-sm text-gray-600">{meeting.attendee_email}</p>
								{#if meeting.company}
									<p class="text-xs text-gray-500">{meeting.company}</p>
								{/if}
							</div>
							<div class="flex-shrink-0 text-right">
								<p class="text-sm font-medium text-gray-900">{formatTime(meeting.meeting_time)}</p>
								<p class="text-xs text-gray-500">{meeting.duration_minutes} min</p>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>

	<!-- Two Column Layout -->
	<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
		<!-- My Assigned Leads -->
		<div class="bg-white rounded-lg shadow-sm border border-gray-200">
			<div class="p-6 border-b border-gray-200 flex items-center justify-between">
				<h3 class="text-lg font-semibold text-gray-900">My Assigned Leads</h3>
				<a href="/dashboard/agent/leads" class="text-sm font-medium text-blue-600 hover:text-blue-700">
					View All
				</a>
			</div>
			<div class="divide-y divide-gray-200">
				{#if data.leads.length === 0}
					<div class="p-6 text-center text-gray-500">
						<p>No leads assigned yet</p>
					</div>
				{:else}
					{#each data.leads.slice(0, 5) as lead}
						<div class="p-4 hover:bg-gray-50 transition-colors">
							<div class="flex items-start justify-between">
								<div class="flex-1 min-w-0">
									<p class="text-sm font-semibold text-gray-900">{lead.name}</p>
									<p class="text-sm text-gray-600">{lead.email}</p>
									{#if lead.company}
										<p class="text-xs text-gray-500">{lead.company}</p>
									{/if}
								</div>
								<div class="flex-shrink-0 ml-4">
									{#if lead.qualification_status === 'qualified'}
										<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
											Qualified
										</span>
									{:else if lead.qualification_status === 'not_qualified'}
										<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
											Not Qualified
										</span>
									{:else}
										<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
											Pending
										</span>
									{/if}
								</div>
							</div>
							{#if lead.qualification_score}
								<div class="mt-2">
									<div class="flex items-center">
										<div class="flex-1">
											<div class="h-2 bg-gray-200 rounded-full overflow-hidden">
												<div
													class="h-full bg-blue-600 rounded-full"
													style="width: {lead.qualification_score}%"
												></div>
											</div>
										</div>
										<span class="ml-2 text-xs font-medium text-gray-600">{lead.qualification_score}</span>
									</div>
								</div>
							{/if}
						</div>
					{/each}
				{/if}
			</div>
		</div>

		<!-- Upcoming Meetings -->
		<div class="bg-white rounded-lg shadow-sm border border-gray-200">
			<div class="p-6 border-b border-gray-200 flex items-center justify-between">
				<h3 class="text-lg font-semibold text-gray-900">Upcoming Meetings</h3>
				<a href="/dashboard/agent/meetings" class="text-sm font-medium text-blue-600 hover:text-blue-700">
					View All
				</a>
			</div>
			<div class="divide-y divide-gray-200">
				{#if data.upcomingMeetings.length === 0}
					<div class="p-6 text-center text-gray-500">
						<p>No upcoming meetings</p>
					</div>
				{:else}
					{#each data.upcomingMeetings.slice(0, 5) as meeting}
						<div class="p-4 hover:bg-gray-50 transition-colors">
							<div class="flex items-start justify-between">
								<div class="flex-1 min-w-0">
									<p class="text-sm font-semibold text-gray-900">{meeting.attendee_name}</p>
									<p class="text-sm text-gray-600">{meeting.attendee_email}</p>
									{#if meeting.company}
										<p class="text-xs text-gray-500">{meeting.company}</p>
									{/if}
								</div>
								<div class="flex-shrink-0 ml-4 text-right">
									<p class="text-xs font-medium text-gray-900">{formatDateTime(meeting.meeting_time)}</p>
									{#if meeting.status === 'scheduled'}
										<span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800 mt-1">
											Scheduled
										</span>
									{/if}
								</div>
							</div>
						</div>
					{/each}
				{/if}
			</div>
		</div>
	</div>

	<!-- Recent Activity -->
	<RecentActivity activities={data.recentActivity} />
</div>
