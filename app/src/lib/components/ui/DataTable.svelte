<script lang="ts">
	interface Column {
		key: string;
		label: string;
		sortable?: boolean;
		type?: 'text' | 'number' | 'currency' | 'percentage' | 'date';
		width?: string;
		align?: 'left' | 'center' | 'right';
		format?: (value: any) => string;
	}

	interface Props {
		columns: Column[];
		data: Record<string, any>[];
		loading?: boolean;
		error?: string;
		sortKey?: string;
		sortOrder?: 'asc' | 'desc';
		pageSize?: number;
		showPagination?: boolean;
	}

	let { 
		columns, 
		data, 
		loading = false,
		error,
		sortKey = $bindable(),
		sortOrder = $bindable('asc'),
		pageSize = 10,
		showPagination = true
	}: Props = $props();

	let currentPage = $state(1);

	// Sort data
	const sortedData = $derived.by(() => {
		if (!sortKey) return data;

		const key = sortKey; // capture in local variable for TypeScript
		return [...data].sort((a, b) => {
			const aVal = a[key];
			const bVal = b[key];

			// Handle null/undefined values
			if (aVal == null && bVal == null) return 0;
			if (aVal == null) return 1;
			if (bVal == null) return -1;

			// Determine sort direction
			const modifier = sortOrder === 'desc' ? -1 : 1;

			// Handle different data types
			if (typeof aVal === 'number' && typeof bVal === 'number') {
				return (aVal - bVal) * modifier;
			}

			// Default string comparison
			return String(aVal).localeCompare(String(bVal)) * modifier;
		});
	});

	// Paginate data
	const paginatedData = $derived.by(() => {
		if (!showPagination) return sortedData;
		const start = (currentPage - 1) * pageSize;
		const end = start + pageSize;
		return sortedData.slice(start, end);
	});

	const totalPages = $derived(Math.ceil(sortedData.length / pageSize));

	function handleSort(columnKey: string) {
		const column = columns.find(col => col.key === columnKey);
		if (!column?.sortable) return;

		if (sortKey === columnKey) {
			sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
		} else {
			sortKey = columnKey;
			sortOrder = 'asc';
		}
	}

	function formatCellValue(value: any, column: Column): string {
		if (value == null) return '-';
		
		if (column.format) {
			return column.format(value);
		}

		switch (column.type) {
			case 'currency':
				return new Intl.NumberFormat('en-US', {
					style: 'currency',
					currency: 'USD'
				}).format(value);
			case 'percentage':
				return `${(value * 100).toFixed(2)}%`;
			case 'number':
				return new Intl.NumberFormat('en-US').format(value);
			case 'date':
				return new Date(value).toLocaleDateString();
			default:
				return String(value);
		}
	}

	function getColumnAlignment(column: Column): string {
		switch (column.align || (column.type === 'number' || column.type === 'currency' ? 'right' : 'left')) {
			case 'center':
				return 'text-center';
			case 'right':
				return 'text-right';
			default:
				return 'text-left';
		}
	}

	function goToPage(page: number) {
		currentPage = Math.max(1, Math.min(page, totalPages));
	}
</script>

<div class="flex flex-col">
	<div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
		<div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
			<div class="overflow-hidden shadow ring-1 ring-border sm:rounded-lg">
				{#if loading}
					<div class="bg-bg px-4 py-12 text-center">
						<div class="inline-flex items-center">
							<div class="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent"></div>
							<span class="ml-2 text-sm text-text-muted">Loading...</span>
						</div>
					</div>
				{:else if error}
					<div class="bg-bg px-4 py-12 text-center">
						<div class="text-sm text-error">Error: {error}</div>
					</div>
				{:else if data.length === 0}
					<div class="bg-bg px-4 py-12 text-center">
						<div class="text-sm text-text-muted">No data available</div>
					</div>
				{:else}
					<table class="min-w-full divide-y divide-border">
						<thead class="bg-bg-secondary">
							<tr>
								{#each columns as column}
									<th
										scope="col"
										class="px-6 py-3 text-xs font-medium uppercase tracking-wide text-text-muted {getColumnAlignment(column)} {column.sortable ? 'cursor-pointer hover:bg-border/50' : ''}"
										class:w-[{column.width}]={column.width}
										onclick={() => column.sortable && handleSort(column.key)}
									>
										<div class="flex items-center {column.align === 'right' ? 'justify-end' : column.align === 'center' ? 'justify-center' : 'justify-start'}">
											{column.label}
											{#if column.sortable}
												<span class="ml-2">
													{#if sortKey === column.key}
														{#if sortOrder === 'asc'}
															<svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
																<path fill-rule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clip-rule="evenodd" />
															</svg>
														{:else}
															<svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
																<path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
															</svg>
														{/if}
													{:else}
														<svg class="h-3 w-3 text-text-muted" fill="currentColor" viewBox="0 0 20 20">
															<path d="M5 12a1 1 0 102 0V6.414l1.293 1.293a1 1 0 001.414-1.414l-3-3a1 1 0 00-1.414 0l-3 3a1 1 0 001.414 1.414L5 6.414V12zM15 8a1 1 0 10-2 0v5.586l-1.293-1.293a1 1 0 00-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L15 13.586V8z" />
														</svg>
													{/if}
												</span>
											{/if}
										</div>
									</th>
								{/each}
							</tr>
						</thead>
						<tbody class="divide-y divide-border bg-bg">
							{#each paginatedData as row, index}
								<tr class="hover:bg-bg-secondary">
									{#each columns as column}
										<td class="whitespace-nowrap px-6 py-4 text-sm text-text {getColumnAlignment(column)}">
											{formatCellValue(row[column.key], column)}
										</td>
									{/each}
								</tr>
							{/each}
						</tbody>
					</table>

					{#if showPagination && totalPages > 1}
						<div class="bg-bg px-4 py-3 flex items-center justify-between border-t border-border sm:px-6">
							<div class="flex flex-1 justify-between sm:hidden">
								<button
									onclick={() => goToPage(currentPage - 1)}
									disabled={currentPage === 1}
									class="btn-secondary"
								>
									Previous
								</button>
								<button
									onclick={() => goToPage(currentPage + 1)}
									disabled={currentPage === totalPages}
									class="btn-secondary ml-3"
								>
									Next
								</button>
							</div>
							<div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
								<div>
									<p class="text-sm text-text-muted">
										Showing
										<span class="font-medium text-text">{(currentPage - 1) * pageSize + 1}</span>
										to
										<span class="font-medium text-text">{Math.min(currentPage * pageSize, sortedData.length)}</span>
										of
										<span class="font-medium text-text">{sortedData.length}</span>
										results
									</p>
								</div>
								<div>
									<nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
										<button
											onclick={() => goToPage(currentPage - 1)}
											disabled={currentPage === 1}
											class="btn-secondary rounded-l-xl rounded-r-none px-2 py-2"
											aria-label="Go to previous page"
										>
											<svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
												<path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
											</svg>
										</button>

										{#each Array.from({length: totalPages}, (_, i) => i + 1) as page}
											{#if page === currentPage}
												<button
													class="btn-primary rounded-none px-4 py-2 text-sm"
												>
													{page}
												</button>
											{:else if Math.abs(page - currentPage) <= 2 || page === 1 || page === totalPages}
												<button
													onclick={() => goToPage(page)}
													class="btn-secondary rounded-none px-4 py-2 text-sm"
												>
													{page}
												</button>
											{:else if Math.abs(page - currentPage) === 3}
												<span class="relative inline-flex items-center px-4 py-2 border-2 border-border bg-bg-secondary text-sm font-medium text-text-muted">
													...
												</span>
											{/if}
										{/each}

										<button
											onclick={() => goToPage(currentPage + 1)}
											disabled={currentPage === totalPages}
											class="btn-secondary rounded-l-none rounded-r-xl px-2 py-2"
											aria-label="Go to next page"
										>
											<svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
												<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
											</svg>
										</button>
									</nav>
								</div>
							</div>
						</div>
					{/if}
				{/if}
			</div>
		</div>
	</div>
</div>

<style>
	@reference "$src/app.css";
	@layer components.data-table {
		/* Data table uses utility classes */
	}
</style>
