<script lang="ts">
	import * as api from '$lib/api';
	import type { AutomationLogEntry } from '$lib/api';
	import { Card, Alert, LoadingSpinner, Badge, Button, Select } from '$lib/components/ui';
	import { Check, X, Clock, Zap, ChevronLeft, ChevronRight, Filter, RefreshCw } from 'lucide-svelte';

	let loading = $state(true);
	let error = $state('');

	// Data
	let entries = $state<AutomationLogEntry[]>([]);
	let total = $state(0);
	let pages = $state(1);

	// Filters
	let filterEventType = $state('');
	let filterSuccess = $state('');
	let currentPage = $state(1);
	const pageSize = 20;

	$effect(() => {
		loadData();
	});

	async function loadData() {
		loading = true;
		error = '';

		try {
			await loadLog();
		} catch (err: any) {
			console.error('Failed to load data:', err);
			error = err.message || 'Failed to load data';
		} finally {
			loading = false;
		}
	}

	async function loadLog() {
		try {
			const params: api.ListAutomationLogRequestParams = {
				page: currentPage,
				page_size: pageSize
			};
			if (filterEventType) params.event_type = filterEventType;
			if (filterSuccess) params.success = filterSuccess;

			const resp = await api.adminListAutomationLog(params);
			entries = resp.entries || [];
			total = resp.total;
			pages = resp.pages;
		} catch (err: any) {
			console.error('Failed to load automation log:', err);
			error = err.message || 'Failed to load automation log';
		}
	}

	function applyFilters() {
		currentPage = 1;
		loadLog();
	}

	function clearFilters() {
		filterEventType = '';
		filterSuccess = '';
		currentPage = 1;
		loadLog();
	}

	function goToPage(page: number) {
		currentPage = page;
		loadLog();
	}

	function formatTime(ms: number): string {
		if (ms < 1000) return `${ms}ms`;
		return `${(ms / 1000).toFixed(2)}s`;
	}

	function formatDate(dateStr: string): string {
		const date = new Date(dateStr);
		return date.toLocaleString();
	}

	function getEventTypeVariant(type: string): 'info' | 'success' | 'warning' | 'default' {
		if (type.includes('ticket')) return 'info';
		if (type.includes('payment')) return 'success';
		if (type.includes('email')) return 'warning';
		return 'default';
	}

	// Expanded row state
	let expandedRowId = $state<string | null>(null);

	function toggleRow(id: string) {
		expandedRowId = expandedRowId === id ? null : id;
	}
</script>

<svelte:head>
	<title>Automation Log</title>
</svelte:head>

{#if error}
	<Alert type="error" title="Error">
		<p>{error}</p>
	</Alert>
{/if}

{#if loading}
	<LoadingSpinner size="large" />
{:else}
	<div class="space-y-6">
		<!-- Header -->
		<Card>
			<div class="flex items-center justify-between">
				<div>
					<h2 class="text-lg font-medium text-text">Automation Log</h2>
					<p class="mt-1 text-sm text-text-muted">
						View history of rule executions and their outcomes
					</p>
				</div>
				<Button type="secondary" onclick={() => loadLog()}>
					<RefreshCw class="mr-2 h-4 w-4" />
					Refresh
				</Button>
			</div>
		</Card>

		<!-- Filters -->
		<Card>
			<div class="flex flex-wrap items-end gap-4">
				<div class="flex-1 min-w-[150px]">
					<label for="filter-event" class="form-label">Event Type</label>
					<Select id="filter-event" bind:value={filterEventType}>
						<option value="">All Events</option>
						<option value="ticket.created">Ticket Created</option>
						<option value="ticket.updated">Ticket Updated</option>
						<option value="payment.succeeded">Payment Succeeded</option>
						<option value="payment.failed">Payment Failed</option>
						<option value="email.bounced">Email Bounced</option>
						<option value="email.complained">Email Complained</option>
						<option value="customer.created">Customer Created</option>
						<option value="contact.created">Contact Created</option>
					</Select>
				</div>

				<div class="flex-1 min-w-[150px]">
					<label for="filter-success" class="form-label">Status</label>
					<Select id="filter-success" bind:value={filterSuccess}>
						<option value="">All</option>
						<option value="true">Success</option>
						<option value="false">Failed</option>
					</Select>
				</div>

				<div class="flex gap-2">
					<Button type="primary" onclick={applyFilters}>
						<Filter class="mr-2 h-4 w-4" />
						Apply
					</Button>
					<Button type="secondary" onclick={clearFilters}>Clear</Button>
				</div>
			</div>
		</Card>

		<!-- Log Table -->
		<Card>
			{#if entries.length === 0}
				<div class="text-center py-12">
					<Zap class="mx-auto h-12 w-12 text-text-muted" />
					<h3 class="mt-2 text-sm font-medium text-text">No automation events</h3>
					<p class="mt-1 text-sm text-text-muted">
						Rule executions will appear here when events trigger your rules.
					</p>
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="min-w-full divide-y divide-border">
						<thead>
							<tr>
								<th class="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">
									Time
								</th>
								<th class="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">
									Event
								</th>
								<th class="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">
									Rule
								</th>
								<th class="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">
									Status
								</th>
								<th class="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">
									Duration
								</th>
								<th class="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">
									Actions
								</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-border">
							{#each entries as entry}
								<tr
									class="hover:bg-bg-secondary cursor-pointer"
									onclick={() => toggleRow(entry.id)}
								>
									<td class="px-4 py-3 text-sm text-text-muted whitespace-nowrap">
										{formatDate(entry.created_at)}
									</td>
									<td class="px-4 py-3 whitespace-nowrap">
										<Badge variant={getEventTypeVariant(entry.event_type)}>
											{entry.event_type}
										</Badge>
									</td>
									<td class="px-4 py-3 text-sm text-text">
										{entry.rule_name}
									</td>
									<td class="px-4 py-3 whitespace-nowrap">
										{#if entry.success}
											<span class="flex items-center text-green-600">
												<Check class="h-4 w-4 mr-1" />
												Success
											</span>
										{:else}
											<span class="flex items-center text-red-600">
												<X class="h-4 w-4 mr-1" />
												Failed
											</span>
										{/if}
									</td>
									<td class="px-4 py-3 text-sm text-text-muted whitespace-nowrap">
										<span class="flex items-center">
											<Clock class="h-3 w-3 mr-1" />
											{formatTime(entry.execution_time_ms)}
										</span>
									</td>
									<td class="px-4 py-3 text-sm text-text-muted">
										{entry.actions_executed?.length || 0} action(s)
									</td>
								</tr>
								{#if expandedRowId === entry.id}
									<tr class="bg-bg-secondary">
										<td colspan="6" class="px-4 py-4">
											<div class="space-y-4">
												{#if entry.error_message}
													<div>
														<h4 class="text-sm font-medium text-red-600 mb-1">Error</h4>
														<p class="text-sm text-red-500 bg-red-50 p-2 rounded">
															{entry.error_message}
														</p>
													</div>
												{/if}

												{#if entry.actions_executed && entry.actions_executed.length > 0}
													<div>
														<h4 class="text-sm font-medium text-text mb-2">Actions Executed</h4>
														<div class="space-y-2">
															{#each entry.actions_executed as action}
																<div class="bg-bg p-2 rounded border border-border">
																	<span class="text-sm font-mono text-primary">{action.type}</span>
																	{#if action.params}
																		<pre class="mt-1 text-xs text-text-muted overflow-x-auto">{JSON.stringify(action.params, null, 2)}</pre>
																	{/if}
																</div>
															{/each}
														</div>
													</div>
												{/if}

												{#if entry.event_payload}
													<div>
														<h4 class="text-sm font-medium text-text mb-2">Event Payload</h4>
														<pre class="text-xs text-text-muted bg-bg p-2 rounded border border-border overflow-x-auto max-h-48">{entry.event_payload}</pre>
													</div>
												{/if}

												<div class="text-xs text-text-muted">
													Event ID: <code class="bg-bg px-1 rounded">{entry.event_id}</code>
												</div>
											</div>
										</td>
									</tr>
								{/if}
							{/each}
						</tbody>
					</table>
				</div>

				<!-- Pagination -->
				{#if pages > 1}
					<div class="flex items-center justify-between border-t border-border px-4 py-3">
						<div class="text-sm text-text-muted">
							Showing {(currentPage - 1) * pageSize + 1} to {Math.min(currentPage * pageSize, total)} of {total} entries
						</div>
						<div class="flex items-center gap-2">
							<Button
								type="secondary"
								size="sm"
								disabled={currentPage <= 1}
								onclick={() => goToPage(currentPage - 1)}
							>
								<ChevronLeft class="h-4 w-4" />
							</Button>
							<span class="text-sm text-text">
								Page {currentPage} of {pages}
							</span>
							<Button
								type="secondary"
								size="sm"
								disabled={currentPage >= pages}
								onclick={() => goToPage(currentPage + 1)}
							>
								<ChevronRight class="h-4 w-4" />
							</Button>
						</div>
					</div>
				{/if}
			{/if}
		</Card>
	</div>
{/if}
