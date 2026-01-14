<script lang="ts">
	let {
		leads = $bindable([]),
		onView = () => {},
		onAssign = () => {},
		onEdit = () => {}
	}: {
		leads: any[];
		onView?: (lead: any) => void;
		onAssign?: (lead: any) => void;
		onEdit?: (lead: any) => void;
	} = $props();

	let sortColumn = $state<string>('created_at');
	let sortDirection = $state<'asc' | 'desc'>('desc');

	function sortBy(column: string) {
		if (sortColumn === column) {
			sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
		} else {
			sortColumn = column;
			sortDirection = 'asc';
		}

		leads = [...leads].sort((a, b) => {
			let aVal = a[column];
			let bVal = b[column];

			if (aVal === null || aVal === undefined) return 1;
			if (bVal === null || bVal === undefined) return -1;

			if (typeof aVal === 'string') aVal = aVal.toLowerCase();
			if (typeof bVal === 'string') bVal = bVal.toLowerCase();

			if (sortDirection === 'asc') {
				return aVal > bVal ? 1 : -1;
			} else {
				return aVal < bVal ? 1 : -1;
			}
		});
	}

	function getStatusColor(status: string): string {
		const colors: Record<string, string> = {
			new: 'bg-blue-100 text-blue-800',
			contacted: 'bg-yellow-100 text-yellow-800',
			qualified: 'bg-green-100 text-green-800',
			unqualified: 'bg-red-100 text-red-800',
			converted: 'bg-purple-100 text-purple-800'
		};
		return colors[status?.toLowerCase()] || 'bg-gray-100 text-gray-800';
	}

	function formatDate(dateString: string): string {
		if (!dateString) return 'N/A';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}
</script>

<div class="overflow-x-auto bg-white rounded-lg shadow">
	<table class="min-w-full divide-y divide-gray-200">
		<thead class="bg-gray-50">
			<tr>
				<th
					class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
					onclick={() => sortBy('name')}
				>
					<div class="flex items-center">
						Name
						{#if sortColumn === 'name'}
							<svg class="ml-1 w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
								{#if sortDirection === 'asc'}
									<path
										fill-rule="evenodd"
										d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"
										clip-rule="evenodd"
									/>
								{:else}
									<path
										fill-rule="evenodd"
										d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"
										clip-rule="evenodd"
									/>
								{/if}
							</svg>
						{/if}
					</div>
				</th>
				<th
					class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
					onclick={() => sortBy('email')}
				>
					Contact
				</th>
				<th
					class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
					onclick={() => sortBy('status')}
				>
					Status
				</th>
				<th
					class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
					onclick={() => sortBy('qualification_score')}
				>
					Score
				</th>
				<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
					Agent
				</th>
				<th
					class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
					onclick={() => sortBy('created_at')}
				>
					Created
				</th>
				<th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
					Actions
				</th>
			</tr>
		</thead>
		<tbody class="bg-white divide-y divide-gray-200">
			{#if leads.length === 0}
				<tr>
					<td colspan="7" class="px-6 py-12 text-center text-gray-500">
						<svg
							class="mx-auto h-12 w-12 text-gray-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
							/>
						</svg>
						<p class="mt-2 text-sm">No leads found</p>
					</td>
				</tr>
			{:else}
				{#each leads as lead (lead.id)}
					<tr class="hover:bg-gray-50">
						<td class="px-6 py-4 whitespace-nowrap">
							<div class="flex items-center">
								<div>
									<div class="text-sm font-medium text-gray-900">
										{lead.name || lead.first_name + ' ' + lead.last_name || 'N/A'}
									</div>
									<div class="text-sm text-gray-500">{lead.company || ''}</div>
								</div>
							</div>
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<div class="text-sm text-gray-900">{lead.email || 'N/A'}</div>
							<div class="text-sm text-gray-500">{lead.phone || ''}</div>
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<span
								class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full {getStatusColor(
									lead.status
								)}"
							>
								{lead.status || 'new'}
							</span>
						</td>
						<td class="px-6 py-4 whitespace-nowrap">
							<div class="flex items-center">
								<div class="text-sm text-gray-900">
									{lead.qualification_score || 0}/100
								</div>
								<div class="ml-2 w-16 bg-gray-200 rounded-full h-2">
									<div
										class="bg-blue-600 h-2 rounded-full"
										style="width: {lead.qualification_score || 0}%"
									></div>
								</div>
							</div>
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
							{lead.assigned_agent?.name || lead.agent_name || 'Unassigned'}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
							{formatDate(lead.created_at)}
						</td>
						<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
							<button
								onclick={() => onView(lead)}
								class="text-blue-600 hover:text-blue-900 mr-3"
								title="View"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
									/>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
									/>
								</svg>
							</button>
							<button
								onclick={() => onEdit(lead)}
								class="text-indigo-600 hover:text-indigo-900 mr-3"
								title="Edit"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
									/>
								</svg>
							</button>
							<button
								onclick={() => onAssign(lead)}
								class="text-green-600 hover:text-green-900"
								title="Assign"
							>
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
									/>
								</svg>
							</button>
						</td>
					</tr>
				{/each}
			{/if}
		</tbody>
	</table>
</div>
