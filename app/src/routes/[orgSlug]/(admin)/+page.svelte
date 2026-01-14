<script lang="ts">
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { Card, Badge, LoadingSpinner, PageHeader, DateFilterPopover } from '$lib/components/ui';
	import { MetricCard, AreaChart, BarChart, SparklineChart, ChartCard } from '$lib/components/charts';
	import { getDashboardStats, type DashboardStatsResponse } from '$lib/api';
	import { getCurrentUser } from '$lib/auth';
	import { TrendingUp, TrendingDown, CheckCircle, Circle, ArrowRight, Users, ShoppingBag, DollarSign, Mail, AlertCircle, CreditCard, Bot, Activity, BarChart3, RefreshCw } from 'lucide-svelte';

	interface Props {
		data: {
			orgSlug?: string;
		};
	}

	const { data }: Props = $props();
	let basePath = $derived(data.orgSlug ? `/${data.orgSlug}` : '');
	let user = $state(getCurrentUser());

	let loading = $state(true);
	let stats = $state<DashboardStatsResponse | null>(null);
	let orgName = $state('');

	// Date range for filtering
	let dateRange = $state({
		from: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
		to: new Date().toISOString().split('T')[0]
	});

	function handleDateChange() {
		loadDashboardData();
	}

	$effect(() => {
		loadDashboardData();
	});

	async function loadDashboardData() {
		loading = true;

		try {
			const orgId = browser ? localStorage.getItem('currentOrgId') : null;
			if (!orgId) {
				loading = false;
				return;
			}

			stats = await getDashboardStats({ from: dateRange.from, to: dateRange.to }, orgId);

			// Get org name from localStorage or stats
			const storedOrgName = browser ? localStorage.getItem('currentOrgName') : null;
			orgName = storedOrgName || 'Your Business';

		} catch (err) {
			console.error('Failed to load dashboard data:', err);
		} finally {
			loading = false;
		}
	}

	function formatCurrency(cents: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',
			minimumFractionDigits: 0,
			maximumFractionDigits: 0
		}).format(cents / 100);
	}

	function formatNumber(num: number): string {
		if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
		if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
		return num.toString();
	}

	function formatDate(dateString: string): string {
		if (!dateString) return 'N/A';
		const date = new Date(dateString);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		const diffHours = Math.floor(diffMs / 3600000);
		const diffDays = Math.floor(diffMs / 86400000);

		if (diffMins < 60) return `${diffMins}m ago`;
		if (diffHours < 24) return `${diffHours}h ago`;
		if (diffDays < 7) return `${diffDays}d ago`;
		return date.toLocaleDateString();
	}

	// Calculate setup progress
	let setupProgress = $derived(() => {
		if (!stats) return 0;
		const checks = [
			stats.stripe_configured,
			stats.has_products,
			stats.has_email_lists,
			stats.has_mcp_configured
		];
		return checks.filter(Boolean).length;
	});

	let isSetupComplete = $derived(setupProgress() === 4);

	// Prepare revenue trend data for chart - format for AreaChart
	const revenueTrendData = $derived.by(() => {
		if (!stats?.revenue_trend || stats.revenue_trend.length === 0) return null;
		return {
			labels: stats.revenue_trend.map((t: { month: string }) => t.month),
			datasets: [{
				label: 'Revenue',
				data: stats.revenue_trend.map((t: { revenue_cents: number }) => t.revenue_cents / 100),
				backgroundColor: '#10b981',
				borderColor: '#10b981'
			}]
		};
	});

	// Calculate MRR (Monthly Recurring Revenue) - use revenue MTD as proxy
	let mrrDisplay = $derived(() => {
		if (!stats) return { value: '$0', trend: 0 };
		const mrr = stats.revenue_mtd_cents / 100;
		return {
			value: formatCurrency(stats.revenue_mtd_cents),
			trend: stats.revenue_growth_pct
		};
	});
</script>

<svelte:head>
	<title>Dashboard - {orgName || 'Admin'}</title>
</svelte:head>

<div class="p-6 space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-semibold text-text">Welcome back, {user?.name || user?.email?.split('@')[0] || 'there'}</h1>
			<p class="mt-1 text-sm text-text-muted">{orgName}</p>
		</div>
		<div class="flex items-center gap-3">
			<DateFilterPopover bind:dateRange onchange={handleDateChange} />
			<button
				onclick={loadDashboardData}
				disabled={loading}
				class="p-2 rounded-lg hover:bg-bg-secondary transition-colors"
				title="Refresh"
			>
				<RefreshCw class="h-5 w-5 text-text-muted {loading ? 'animate-spin' : ''}" />
			</button>
		</div>
	</div>

		{#if loading}
			<div class="flex items-center justify-center h-64">
				<LoadingSpinner size="large" />
			</div>
		{:else if stats}
			<!-- Setup Checklist (only show if not complete) -->
			{#if !isSetupComplete}
				<Card class="mb-6">
					<div class="flex items-center justify-between mb-4">
						<div>
							<h2 class="text-lg font-semibold text-text">Get Started</h2>
							<p class="text-sm text-text-muted">{setupProgress()} of 4 tasks complete</p>
						</div>
						<div class="flex items-center gap-2">
							<div class="w-32 h-2 bg-gray-200 rounded-full overflow-hidden">
								<div
									class="h-full bg-primary rounded-full transition-all"
									style="width: {(setupProgress() / 4) * 100}%"
								></div>
							</div>
						</div>
					</div>
					<div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
						<a
							href="{basePath}/settings/billing"
							class="flex items-center gap-2 p-3 rounded-lg border border-border hover:bg-bg-secondary transition-colors {stats.stripe_configured ? 'opacity-60' : ''}"
						>
							{#if stats.stripe_configured}
								<CheckCircle class="h-4 w-4 text-green-500 flex-shrink-0" />
							{:else}
								<Circle class="h-4 w-4 text-text-muted flex-shrink-0" />
							{/if}
							<div class="min-w-0">
								<p class="text-sm font-medium text-text truncate">Connect Stripe</p>
								<p class="text-xs text-text-muted">Accept payments</p>
							</div>
						</a>
						<a
							href="{basePath}/commerce/products"
							class="flex items-center gap-2 p-3 rounded-lg border border-border hover:bg-bg-secondary transition-colors {stats.has_products ? 'opacity-60' : ''}"
						>
							{#if stats.has_products}
								<CheckCircle class="h-4 w-4 text-green-500 flex-shrink-0" />
							{:else}
								<Circle class="h-4 w-4 text-text-muted flex-shrink-0" />
							{/if}
							<div class="min-w-0">
								<p class="text-sm font-medium text-text truncate">Add product</p>
								<p class="text-xs text-text-muted">Create catalog</p>
							</div>
						</a>
						<a
							href="{basePath}/email"
							class="flex items-center gap-2 p-3 rounded-lg border border-border hover:bg-bg-secondary transition-colors {stats.has_email_lists ? 'opacity-60' : ''}"
						>
							{#if stats.has_email_lists}
								<CheckCircle class="h-4 w-4 text-green-500 flex-shrink-0" />
							{:else}
								<Circle class="h-4 w-4 text-text-muted flex-shrink-0" />
							{/if}
							<div class="min-w-0">
								<p class="text-sm font-medium text-text truncate">Email list</p>
								<p class="text-xs text-text-muted">Build audience</p>
							</div>
						</a>
						<a
							href="{basePath}/settings/mcp"
							class="flex items-center gap-2 p-3 rounded-lg border border-border hover:bg-bg-secondary transition-colors {stats.has_mcp_configured ? 'opacity-60' : ''}"
						>
							{#if stats.has_mcp_configured}
								<CheckCircle class="h-4 w-4 text-green-500 flex-shrink-0" />
							{:else}
								<Circle class="h-4 w-4 text-text-muted flex-shrink-0" />
							{/if}
							<div class="min-w-0">
								<p class="text-sm font-medium text-text truncate">Configure MCP</p>
								<p class="text-xs text-text-muted">AI assistant</p>
							</div>
						</a>
					</div>
				</Card>
			{/if}

			<!-- Needs Attention (show if any issues) -->
			{#if stats.needs_attention && stats.needs_attention.length > 0}
				<Card class="mb-6 border-l-4 border-l-warning">
					<h2 class="text-lg font-semibold text-text mb-4 flex items-center gap-2">
						<AlertCircle class="h-5 w-5 text-warning" />
						Needs Attention
					</h2>
					<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
						{#each stats.needs_attention as item}
							<a
								href="{basePath}/{item.type === 'failed_payments' || item.type === 'pending_orders' ? 'commerce/orders' : 'llm-logs'}"
								class="flex items-center justify-between p-3 bg-bg-secondary rounded-lg hover:bg-bg-tertiary transition-colors"
							>
								<div>
									<p class="text-sm font-medium text-text">{item.count} {item.type.replace('_', ' ')}</p>
									<p class="text-xs text-text-muted">{item.description}</p>
								</div>
								<ArrowRight class="h-4 w-4 text-text-muted" />
							</a>
						{/each}
					</div>
				</Card>
			{/if}

			<!-- KPI Cards Row - Primary Metrics -->
			<div class="grid grid-cols-2 lg:grid-cols-5 gap-4 mb-6">
				<!-- MRR / Revenue MTD -->
				<Card class="relative overflow-hidden">
					<div class="flex items-start justify-between">
						<div class="flex-1 min-w-0">
							<p class="text-xs font-medium text-text-muted uppercase tracking-wide">Revenue MTD</p>
							<p class="text-2xl font-bold text-text mt-1 truncate">{formatCurrency(stats.revenue_mtd_cents)}</p>
							{#if stats.revenue_growth_pct !== 0}
								<p class="text-xs mt-1 flex items-center gap-1 {stats.revenue_growth_pct > 0 ? 'text-green-600' : 'text-red-600'}">
									{#if stats.revenue_growth_pct > 0}
										<TrendingUp class="h-3 w-3" />
									{:else}
										<TrendingDown class="h-3 w-3" />
									{/if}
									{stats.revenue_growth_pct > 0 ? '+' : ''}{stats.revenue_growth_pct.toFixed(1)}% vs last month
								</p>
							{:else}
								<p class="text-xs text-text-muted mt-1">vs last month</p>
							{/if}
						</div>
						<div class="p-2 bg-green-100 rounded-lg flex-shrink-0">
							<DollarSign class="h-5 w-5 text-green-600" />
						</div>
					</div>
				</Card>

				<!-- Revenue YTD -->
				<Card class="relative overflow-hidden">
					<div class="flex items-start justify-between">
						<div class="flex-1 min-w-0">
							<p class="text-xs font-medium text-text-muted uppercase tracking-wide">Revenue YTD</p>
							<p class="text-2xl font-bold text-text mt-1 truncate">{formatCurrency(stats.revenue_ytd_cents)}</p>
							<p class="text-xs text-text-muted mt-1">{stats.paid_orders_count} orders</p>
						</div>
						<div class="p-2 bg-blue-100 rounded-lg flex-shrink-0">
							<BarChart3 class="h-5 w-5 text-blue-600" />
						</div>
					</div>
				</Card>

				<!-- Customers -->
				<Card class="relative overflow-hidden">
					<div class="flex items-start justify-between">
						<div class="flex-1 min-w-0">
							<p class="text-xs font-medium text-text-muted uppercase tracking-wide">Customers</p>
							<p class="text-2xl font-bold text-text mt-1">{formatNumber(stats.total_customers)}</p>
							{#if stats.new_customers_30d > 0}
								<p class="text-xs text-green-600 mt-1 flex items-center gap-1">
									<TrendingUp class="h-3 w-3" />
									+{stats.new_customers_30d} this month
								</p>
							{:else}
								<p class="text-xs text-text-muted mt-1">This month: 0 new</p>
							{/if}
						</div>
						<div class="p-2 bg-purple-100 rounded-lg flex-shrink-0">
							<Users class="h-5 w-5 text-purple-600" />
						</div>
					</div>
				</Card>

				<!-- Subscriptions / Churn -->
				<Card class="relative overflow-hidden">
					<div class="flex items-start justify-between">
						<div class="flex-1 min-w-0">
							<p class="text-xs font-medium text-text-muted uppercase tracking-wide">Active Subs</p>
							<p class="text-2xl font-bold text-text mt-1">{stats.active_subscriptions}</p>
							{#if stats.churn_rate > 0}
								<p class="text-xs text-amber-600 mt-1">{stats.churn_rate.toFixed(1)}% churn rate</p>
							{:else}
								<p class="text-xs text-green-600 mt-1">0% churn</p>
							{/if}
						</div>
						<div class="p-2 bg-indigo-100 rounded-lg flex-shrink-0">
							<CreditCard class="h-5 w-5 text-indigo-600" />
						</div>
					</div>
				</Card>

				<!-- LLM Usage -->
				<Card class="relative overflow-hidden">
					<div class="flex items-start justify-between">
						<div class="flex-1 min-w-0">
							<p class="text-xs font-medium text-text-muted uppercase tracking-wide">LLM Tokens</p>
							<p class="text-2xl font-bold text-text mt-1">{formatNumber(stats.llm_total_tokens)}</p>
							<p class="text-xs text-text-muted mt-1">${(stats.llm_total_cost_micros / 1000000).toFixed(2)} cost</p>
						</div>
						<div class="p-2 bg-cyan-100 rounded-lg flex-shrink-0">
							<Bot class="h-5 w-5 text-cyan-600" />
						</div>
					</div>
				</Card>
			</div>

			<!-- Charts Row -->
			<div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
				<!-- Revenue Trend Chart -->
				<ChartCard title="Revenue Trend" subtitle="Last 6 months">
					{#if revenueTrendData}
						<AreaChart
							data={revenueTrendData}
							height={200}
							showGrid={true}
						/>
					{:else}
						<div class="flex items-center justify-center h-[200px] text-text-muted">
							<p>No revenue data yet</p>
						</div>
					{/if}
				</ChartCard>

				<!-- Recent Orders -->
				<Card>
					<div class="flex items-center justify-between mb-4">
						<h2 class="text-lg font-semibold text-text">Recent Orders</h2>
						<a href="{basePath}/commerce/orders" class="text-sm font-medium text-primary hover:underline flex items-center gap-1">
							View all
							<ArrowRight class="h-4 w-4" />
						</a>
					</div>
					<div class="space-y-2">
						{#if stats.recent_orders && stats.recent_orders.length > 0}
							{#each stats.recent_orders as order}
								<div class="flex items-center justify-between py-2 border-b border-border last:border-0">
									<div class="min-w-0 flex-1">
										<p class="text-sm font-medium text-text truncate">
											{order.customer_email || 'Guest'}
										</p>
										<p class="text-xs text-text-muted">{formatDate(order.created_at)}</p>
									</div>
									<div class="flex items-center gap-3 ml-4">
										<span class="text-sm font-medium text-text">
											{formatCurrency(order.total_cents || 0)}
										</span>
										<Badge
											variant={order.payment_status === 'paid' ? 'success' : order.payment_status === 'failed' ? 'error' : 'warning'}
											size="sm"
										>
											{order.payment_status || 'pending'}
										</Badge>
									</div>
								</div>
							{/each}
						{:else}
							<div class="text-center py-8">
								<ShoppingBag class="h-8 w-8 text-text-muted mx-auto mb-2" />
								<p class="text-sm text-text-muted">No orders yet</p>
								<a href="{basePath}/commerce/products" class="text-sm text-primary hover:underline mt-1 inline-block">
									Add products to get started
								</a>
							</div>
						{/if}
					</div>
				</Card>
			</div>

			<!-- Secondary Metrics Row -->
			<div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
				<!-- Email Stats -->
				<Card>
					<div class="flex items-center gap-3">
						<div class="p-2 bg-blue-100 rounded-lg">
							<Mail class="h-4 w-4 text-blue-600" />
						</div>
						<div>
							<p class="text-xs text-text-muted">Emails Sent (30d)</p>
							<p class="text-lg font-semibold text-text">{formatNumber(stats.emails_sent_30d)}</p>
						</div>
					</div>
					{#if stats.emails_sent_30d > 0}
						<div class="mt-3 pt-3 border-t border-border">
							<div class="flex justify-between text-xs">
								<span class="text-text-muted">Open rate</span>
								<span class="font-medium text-text">{stats.email_open_rate.toFixed(1)}%</span>
							</div>
						</div>
					{/if}
				</Card>

				<!-- Orders This Month -->
				<Card>
					<div class="flex items-center gap-3">
						<div class="p-2 bg-green-100 rounded-lg">
							<ShoppingBag class="h-4 w-4 text-green-600" />
						</div>
						<div>
							<p class="text-xs text-text-muted">Orders (MTD)</p>
							<p class="text-lg font-semibold text-text">{stats.orders_mtd_count}</p>
						</div>
					</div>
					{#if stats.paid_orders_count > 0}
						<div class="mt-3 pt-3 border-t border-border">
							<div class="flex justify-between text-xs">
								<span class="text-text-muted">Avg order value</span>
								<span class="font-medium text-text">{formatCurrency(stats.total_revenue_cents / stats.paid_orders_count)}</span>
							</div>
						</div>
					{/if}
				</Card>

				<!-- LLM Performance -->
				<Card>
					<div class="flex items-center gap-3">
						<div class="p-2 bg-cyan-100 rounded-lg">
							<Activity class="h-4 w-4 text-cyan-600" />
						</div>
						<div>
							<p class="text-xs text-text-muted">LLM Requests (30d)</p>
							<p class="text-lg font-semibold text-text">{formatNumber(stats.llm_total_requests)}</p>
						</div>
					</div>
					{#if stats.llm_total_requests > 0}
						<div class="mt-3 pt-3 border-t border-border">
							<div class="flex justify-between text-xs">
								<span class="text-text-muted">Success rate</span>
								<span class="font-medium text-text">{((stats.llm_success_count / stats.llm_total_requests) * 100).toFixed(1)}%</span>
							</div>
						</div>
					{/if}
				</Card>

				<!-- New Subscriptions -->
				<Card>
					<div class="flex items-center gap-3">
						<div class="p-2 bg-purple-100 rounded-lg">
							<TrendingUp class="h-4 w-4 text-purple-600" />
						</div>
						<div>
							<p class="text-xs text-text-muted">New Subs (MTD)</p>
							<p class="text-lg font-semibold text-text">{stats.new_subscriptions_mtd}</p>
						</div>
					</div>
					{#if stats.trialing_count > 0}
						<div class="mt-3 pt-3 border-t border-border">
							<div class="flex justify-between text-xs">
								<span class="text-text-muted">In trial</span>
								<span class="font-medium text-text">{stats.trialing_count}</span>
							</div>
						</div>
					{/if}
				</Card>
			</div>
		{:else}
			<Card>
				<div class="text-center py-12">
					<Activity class="h-12 w-12 text-text-muted mx-auto mb-4" />
					<h3 class="text-lg font-semibold text-text mb-2">No data available</h3>
					<p class="text-text-muted mb-4">Select an organization to view your dashboard</p>
				</div>
			</Card>
		{/if}
</div>
