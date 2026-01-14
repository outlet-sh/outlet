<script lang="ts">
	import * as api from '$lib/api';
	import { Button, Input, Select, Card, Alert, Badge, Collapsible } from '$lib/components/ui';
	import { Play, CheckCircle, XCircle, Clock, Zap, AlertTriangle } from 'lucide-svelte';
	import { FACT_CATEGORIES, ACTION_TYPES } from './types';

	let {
		ruleId,
		onClose
	}: {
		ruleId: string;
		onClose?: () => void;
	} = $props();

	// Test data state - flattened structure for easy binding
	let ticketData = $state({
		enabled: false,
		status: '',
		priority: '',
		category: '',
		subject: ''
	});

	let customerData = $state({
		enabled: false,
		email: '',
		ltv: '',
		plan: ''
	});

	let contactData = $state({
		enabled: false,
		email: '',
		is_new: 'false',
		is_subscribed: 'true'
	});

	let subscriptionData = $state({
		enabled: false,
		status: ''
	});

	let paymentData = $state({
		enabled: false,
		status: ''
	});

	let invoiceData = $state({
		enabled: false,
		status: ''
	});

	let computedData = $state({
		hours_since_ticket_created: '',
		hours_since_ticket_updated: '',
		days_since_last_payment: '',
		days_since_subscription_start: '',
		emails_sent_last_24_hours: '',
		bounce_count_last_30_days: ''
	});

	// UI state
	let testing = $state(false);
	let testResult = $state<api.TestRuleResponse | null>(null);
	let testError = $state('');

	// Build test data JSON from form state
	function buildTestData(): string {
		const facts: Record<string, unknown> = {};

		if (ticketData.enabled) {
			facts.ticket = {
				status: ticketData.status || undefined,
				priority: ticketData.priority || undefined,
				category: ticketData.category || undefined,
				subject: ticketData.subject || undefined
			};
		}

		if (customerData.enabled) {
			facts.customer = {
				email: customerData.email || undefined,
				ltv: customerData.ltv ? parseFloat(customerData.ltv) : undefined,
				plan: customerData.plan || undefined
			};
		}

		if (contactData.enabled) {
			facts.contact = {
				email: contactData.email || undefined,
				is_new: contactData.is_new === 'true',
				is_subscribed: contactData.is_subscribed === 'true'
			};
		}

		if (subscriptionData.enabled) {
			facts.subscription = {
				status: subscriptionData.status || undefined
			};
		}

		if (paymentData.enabled) {
			facts.payment = {
				status: paymentData.status || undefined
			};
		}

		if (invoiceData.enabled) {
			facts.invoice = {
				status: invoiceData.status || undefined
			};
		}

		// Add computed fields
		if (computedData.hours_since_ticket_created) {
			facts.hours_since_ticket_created = parseFloat(computedData.hours_since_ticket_created);
		}
		if (computedData.hours_since_ticket_updated) {
			facts.hours_since_ticket_updated = parseFloat(computedData.hours_since_ticket_updated);
		}
		if (computedData.days_since_last_payment) {
			facts.days_since_last_payment = parseInt(computedData.days_since_last_payment);
		}
		if (computedData.days_since_subscription_start) {
			facts.days_since_subscription_start = parseInt(computedData.days_since_subscription_start);
		}
		if (computedData.emails_sent_last_24_hours) {
			facts.emails_sent_last_24_hours = parseInt(computedData.emails_sent_last_24_hours);
		}
		if (computedData.bounce_count_last_30_days) {
			facts.bounce_count_last_30_days = parseInt(computedData.bounce_count_last_30_days);
		}

		// Remove undefined values
		const cleanFacts = JSON.parse(JSON.stringify(facts));
		return JSON.stringify(cleanFacts, null, 2);
	}

	async function runTest() {
		testing = true;
		testError = '';
		testResult = null;

		try {
			const testData = buildTestData();
			const result = await api.adminTestRule({}, { test_data: testData }, ruleId);
			testResult = result;

			if (result.errors && result.errors.length > 0) {
				testError = result.errors.join('; ');
			}
		} catch (err: unknown) {
			const error = err as { message?: string };
			testError = error.message || 'Test execution failed';
		} finally {
			testing = false;
		}
	}

	function getActionLabel(actionType: string): string {
		const action = ACTION_TYPES.find((a) => a.type === actionType);
		return action?.name || actionType;
	}

	// Get options for fact fields
	function getFieldOptions(path: string) {
		for (const cat of FACT_CATEGORIES) {
			for (const field of cat.fields) {
				if (field.path === path) {
					return field.options;
				}
			}
		}
		return null;
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h3 class="text-lg font-medium text-text">Test Rule</h3>
			<p class="text-sm text-text-muted">Simulate facts to see which rules fire</p>
		</div>
		{#if onClose}
			<Button type="ghost" size="sm" onclick={onClose}>Close</Button>
		{/if}
	</div>

	<!-- Fact Input Sections -->
	<div class="space-y-3">
		<!-- Ticket Facts -->
		<Collapsible title="Ticket" defaultOpen={ticketData.enabled}>
			{#snippet header()}
				<div class="flex items-center gap-2">
					<input
						type="checkbox"
						bind:checked={ticketData.enabled}
						class="rounded border-border"
						onclick={(e) => e.stopPropagation()}
					/>
					<span class="font-medium text-text">Ticket</span>
					{#if ticketData.enabled}
						<Badge variant="info">Active</Badge>
					{/if}
				</div>
			{/snippet}

			<div class="grid grid-cols-2 gap-3 pt-3">
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Status</label>
					<Select bind:value={ticketData.status} size="sm">
						<option value="">Select...</option>
						{#each getFieldOptions('Ticket.Status') || [] as opt}
							<option value={opt.value}>{opt.label}</option>
						{/each}
					</Select>
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Priority</label>
					<Select bind:value={ticketData.priority} size="sm">
						<option value="">Select...</option>
						{#each getFieldOptions('Ticket.Priority') || [] as opt}
							<option value={opt.value}>{opt.label}</option>
						{/each}
					</Select>
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Category</label>
					<Select bind:value={ticketData.category} size="sm">
						<option value="">Select...</option>
						{#each getFieldOptions('Ticket.Category') || [] as opt}
							<option value={opt.value}>{opt.label}</option>
						{/each}
					</Select>
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Subject</label>
					<Input type="text" bind:value={ticketData.subject} size="sm" placeholder="Ticket subject" />
				</div>
			</div>
		</Collapsible>

		<!-- Customer Facts -->
		<Collapsible title="Customer" defaultOpen={customerData.enabled}>
			{#snippet header()}
				<div class="flex items-center gap-2">
					<input
						type="checkbox"
						bind:checked={customerData.enabled}
						class="rounded border-border"
						onclick={(e) => e.stopPropagation()}
					/>
					<span class="font-medium text-text">Customer</span>
					{#if customerData.enabled}
						<Badge variant="info">Active</Badge>
					{/if}
				</div>
			{/snippet}

			<div class="grid grid-cols-2 gap-3 pt-3">
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Email</label>
					<Input type="email" bind:value={customerData.email} size="sm" placeholder="customer@example.com" />
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Lifetime Value</label>
					<Input type="number" bind:value={customerData.ltv} size="sm" placeholder="0" />
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Plan</label>
					<Select bind:value={customerData.plan} size="sm">
						<option value="">Select...</option>
						{#each getFieldOptions('Customer.Plan') || [] as opt}
							<option value={opt.value}>{opt.label}</option>
						{/each}
					</Select>
				</div>
			</div>
		</Collapsible>

		<!-- Contact Facts -->
		<Collapsible title="Contact" defaultOpen={contactData.enabled}>
			{#snippet header()}
				<div class="flex items-center gap-2">
					<input
						type="checkbox"
						bind:checked={contactData.enabled}
						class="rounded border-border"
						onclick={(e) => e.stopPropagation()}
					/>
					<span class="font-medium text-text">Contact</span>
					{#if contactData.enabled}
						<Badge variant="info">Active</Badge>
					{/if}
				</div>
			{/snippet}

			<div class="grid grid-cols-2 gap-3 pt-3">
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Email</label>
					<Input type="email" bind:value={contactData.email} size="sm" placeholder="contact@example.com" />
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Is New</label>
					<Select bind:value={contactData.is_new} size="sm">
						<option value="true">Yes</option>
						<option value="false">No</option>
					</Select>
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Is Subscribed</label>
					<Select bind:value={contactData.is_subscribed} size="sm">
						<option value="true">Yes</option>
						<option value="false">No</option>
					</Select>
				</div>
			</div>
		</Collapsible>

		<!-- Subscription Facts -->
		<Collapsible title="Subscription" defaultOpen={subscriptionData.enabled}>
			{#snippet header()}
				<div class="flex items-center gap-2">
					<input
						type="checkbox"
						bind:checked={subscriptionData.enabled}
						class="rounded border-border"
						onclick={(e) => e.stopPropagation()}
					/>
					<span class="font-medium text-text">Subscription</span>
					{#if subscriptionData.enabled}
						<Badge variant="info">Active</Badge>
					{/if}
				</div>
			{/snippet}

			<div class="grid grid-cols-2 gap-3 pt-3">
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Status</label>
					<Select bind:value={subscriptionData.status} size="sm">
						<option value="">Select...</option>
						{#each getFieldOptions('Subscription.Status') || [] as opt}
							<option value={opt.value}>{opt.label}</option>
						{/each}
					</Select>
				</div>
			</div>
		</Collapsible>

		<!-- Payment Facts -->
		<Collapsible title="Payment" defaultOpen={paymentData.enabled}>
			{#snippet header()}
				<div class="flex items-center gap-2">
					<input
						type="checkbox"
						bind:checked={paymentData.enabled}
						class="rounded border-border"
						onclick={(e) => e.stopPropagation()}
					/>
					<span class="font-medium text-text">Payment</span>
					{#if paymentData.enabled}
						<Badge variant="info">Active</Badge>
					{/if}
				</div>
			{/snippet}

			<div class="grid grid-cols-2 gap-3 pt-3">
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Status</label>
					<Select bind:value={paymentData.status} size="sm">
						<option value="">Select...</option>
						{#each getFieldOptions('Payment.Status') || [] as opt}
							<option value={opt.value}>{opt.label}</option>
						{/each}
					</Select>
				</div>
			</div>
		</Collapsible>

		<!-- Computed Fields -->
		<Collapsible title="Computed Fields">
			{#snippet header()}
				<span class="font-medium text-text">Computed Fields</span>
			{/snippet}

			<div class="grid grid-cols-2 gap-3 pt-3">
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Hours Since Ticket Created</label>
					<Input type="number" bind:value={computedData.hours_since_ticket_created} size="sm" placeholder="0" />
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Hours Since Ticket Updated</label>
					<Input type="number" bind:value={computedData.hours_since_ticket_updated} size="sm" placeholder="0" />
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Days Since Last Payment</label>
					<Input type="number" bind:value={computedData.days_since_last_payment} size="sm" placeholder="0" />
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Days Since Subscription Start</label>
					<Input type="number" bind:value={computedData.days_since_subscription_start} size="sm" placeholder="0" />
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Emails Sent (24h)</label>
					<Input type="number" bind:value={computedData.emails_sent_last_24_hours} size="sm" placeholder="0" />
				</div>
				<div>
					<label class="block text-xs font-medium text-text-muted mb-1">Bounces (30d)</label>
					<Input type="number" bind:value={computedData.bounce_count_last_30_days} size="sm" placeholder="0" />
				</div>
			</div>
		</Collapsible>
	</div>

	<!-- Run Test Button -->
	<div class="flex justify-end">
		<Button type="primary" onclick={runTest} disabled={testing}>
			<Play class="h-4 w-4 mr-2" />
			{testing ? 'Running...' : 'Run Test'}
		</Button>
	</div>

	<!-- Test Results -->
	{#if testError && !testResult}
		<Alert type="error" title="Test Failed">
			<p>{testError}</p>
		</Alert>
	{/if}

	{#if testResult}
		<Card>
			<div class="space-y-4">
				<!-- Summary Stats -->
				<div class="flex items-center gap-6">
					<div class="flex items-center gap-2">
						{#if testResult.rules_fired > 0}
							<CheckCircle class="h-5 w-5 text-green-500" />
						{:else}
							<XCircle class="h-5 w-5 text-text-muted" />
						{/if}
						<span class="text-sm">
							<span class="font-medium">{testResult.rules_fired}</span> of {testResult.rules_evaluated} rules fired
						</span>
					</div>
					<div class="flex items-center gap-2 text-sm text-text-muted">
						<Clock class="h-4 w-4" />
						<span>{testResult.duration_ms}ms</span>
					</div>
				</div>

				<!-- Errors -->
				{#if testResult.errors && testResult.errors.length > 0}
					<div class="p-3 bg-red-50 dark:bg-red-900/20 rounded-lg">
						<div class="flex items-center gap-2 text-red-700 dark:text-red-400 mb-2">
							<AlertTriangle class="h-4 w-4" />
							<span class="font-medium text-sm">Errors</span>
						</div>
						<ul class="text-sm text-red-600 dark:text-red-300 space-y-1">
							{#each testResult.errors as error}
								<li>{error}</li>
							{/each}
						</ul>
					</div>
				{/if}

				<!-- Actions -->
				{#if testResult.actions && testResult.actions.length > 0}
					<div>
						<h4 class="text-xs font-medium text-text uppercase tracking-wider mb-2">Actions That Would Execute</h4>
						<div class="space-y-2">
							{#each testResult.actions as action}
								<div class="flex items-start gap-2 p-2 bg-bg-secondary rounded">
									<Zap class="h-4 w-4 text-primary flex-shrink-0 mt-0.5" />
									<div class="flex-1 min-w-0">
										<span class="font-medium text-sm text-text">{getActionLabel(action.type)}</span>
										{#if action.params}
											<pre class="text-xs text-text-muted mt-1 overflow-x-auto">{action.params}</pre>
										{/if}
									</div>
								</div>
							{/each}
						</div>
					</div>
				{:else if testResult.rules_fired === 0}
					<p class="text-sm text-text-muted text-center py-4">
						No rules matched the provided facts
					</p>
				{/if}
			</div>
		</Card>
	{/if}
</div>
