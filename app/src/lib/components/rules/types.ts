// Rule Builder Types
// These types define the structure for the visual rule builder

// =============================================================================
// FACTS - Available data that rules can evaluate
// =============================================================================

export interface FactField {
	name: string; // Display name
	path: string; // Grule expression path (e.g., "Ticket.Status")
	type: 'string' | 'number' | 'boolean' | 'date';
	options?: { value: string; label: string }[]; // For enum-like fields
	description?: string;
}

export interface FactCategory {
	name: string;
	description: string;
	fields: FactField[];
}

// All available facts organized by category
export const FACT_CATEGORIES: FactCategory[] = [
	{
		name: 'Ticket',
		description: 'Support ticket information',
		fields: [
			{
				name: 'Status',
				path: 'Ticket.Status',
				type: 'string',
				options: [
					{ value: 'open', label: 'Open' },
					{ value: 'pending', label: 'Pending' },
					{ value: 'resolved', label: 'Resolved' },
					{ value: 'closed', label: 'Closed' }
				]
			},
			{
				name: 'Priority',
				path: 'Ticket.Priority',
				type: 'string',
				options: [
					{ value: 'low', label: 'Low' },
					{ value: 'normal', label: 'Normal' },
					{ value: 'high', label: 'High' },
					{ value: 'urgent', label: 'Urgent' }
				]
			},
			{
				name: 'Category',
				path: 'Ticket.Category',
				type: 'string',
				options: [
					{ value: 'billing', label: 'Billing' },
					{ value: 'technical', label: 'Technical' },
					{ value: 'general', label: 'General' },
					{ value: 'feature_request', label: 'Feature Request' }
				]
			},
			{ name: 'Subject', path: 'Ticket.Subject', type: 'string' },
			{ name: 'Hours Since Created', path: 'HoursSinceTicketCreated', type: 'number' },
			{ name: 'Hours Since Updated', path: 'HoursSinceTicketUpdated', type: 'number' }
		]
	},
	{
		name: 'Customer',
		description: 'Customer information',
		fields: [
			{ name: 'Email', path: 'Customer.Email', type: 'string' },
			{ name: 'Lifetime Value', path: 'Customer.Ltv', type: 'number' },
			{
				name: 'Plan',
				path: 'Customer.Plan',
				type: 'string',
				options: [
					{ value: 'free', label: 'Free' },
					{ value: 'pro', label: 'Pro' },
					{ value: 'team', label: 'Team' },
					{ value: 'enterprise', label: 'Enterprise' }
				]
			}
		]
	},
	{
		name: 'Contact',
		description: 'Contact/subscriber information',
		fields: [
			{ name: 'Email', path: 'Contact.Email', type: 'string' },
			{ name: 'Is New', path: 'Contact.IsNew', type: 'boolean' },
			{ name: 'Is Subscribed', path: 'Contact.IsSubscribed', type: 'boolean' }
		]
	},
	{
		name: 'Subscription',
		description: 'Subscription information',
		fields: [
			{
				name: 'Status',
				path: 'Subscription.Status',
				type: 'string',
				options: [
					{ value: 'active', label: 'Active' },
					{ value: 'past_due', label: 'Past Due' },
					{ value: 'canceled', label: 'Canceled' },
					{ value: 'trialing', label: 'Trialing' },
					{ value: 'paused', label: 'Paused' }
				]
			},
			{ name: 'Days Since Start', path: 'DaysSinceSubscriptionStart', type: 'number' }
		]
	},
	{
		name: 'Payment',
		description: 'Payment information',
		fields: [
			{
				name: 'Status',
				path: 'Payment.Status',
				type: 'string',
				options: [
					{ value: 'succeeded', label: 'Succeeded' },
					{ value: 'failed', label: 'Failed' },
					{ value: 'pending', label: 'Pending' },
					{ value: 'refunded', label: 'Refunded' }
				]
			},
			{ name: 'Days Since Payment', path: 'DaysSinceLastPayment', type: 'number' }
		]
	},
	{
		name: 'Invoice',
		description: 'Invoice information',
		fields: [
			{
				name: 'Status',
				path: 'Invoice.Status',
				type: 'string',
				options: [
					{ value: 'draft', label: 'Draft' },
					{ value: 'open', label: 'Open' },
					{ value: 'paid', label: 'Paid' },
					{ value: 'void', label: 'Void' },
					{ value: 'uncollectible', label: 'Uncollectible' }
				]
			}
		]
	},
	{
		name: 'Email Activity',
		description: 'Email delivery metrics',
		fields: [
			{ name: 'Emails Sent (24h)', path: 'EmailsSentLast24Hours', type: 'number' },
			{ name: 'Bounces (30d)', path: 'BounceCountLast30Days', type: 'number' }
		]
	},
	{
		name: 'Trial & Onboarding',
		description: 'Trial and onboarding status',
		fields: [
			{ name: 'Trial Days Remaining', path: 'TrialDaysRemaining', type: 'number' },
			{ name: 'Is Trialing', path: 'IsTrialing', type: 'boolean' },
			{ name: 'Has Completed Onboarding', path: 'HasCompletedOnboarding', type: 'boolean' },
			{
				name: 'Onboarding Step',
				path: 'OnboardingStep',
				type: 'string',
				options: [
					{ value: 'signup', label: 'Signup' },
					{ value: 'profile', label: 'Profile Setup' },
					{ value: 'integration', label: 'Integration' },
					{ value: 'first_action', label: 'First Action' },
					{ value: 'complete', label: 'Complete' }
				]
			}
		]
	},
	{
		name: 'Usage & Product',
		description: 'Product usage metrics',
		fields: [
			{ name: 'API Calls (This Month)', path: 'ApiCallsThisMonth', type: 'number' },
			{ name: 'Storage Used (MB)', path: 'StorageUsedMb', type: 'number' },
			{ name: 'Active Team Members', path: 'ActiveTeamMembers', type: 'number' },
			{ name: 'Contacts Count', path: 'ContactsCount', type: 'number' },
			{ name: 'Emails Sent (This Month)', path: 'EmailsSentThisMonth', type: 'number' },
			{ name: 'Has Exceeded Limit', path: 'HasExceededLimit', type: 'boolean' },
			{ name: 'Usage Percent', path: 'UsagePercent', type: 'number' }
		]
	},
	{
		name: 'Engagement',
		description: 'User engagement metrics',
		fields: [
			{ name: 'Days Since Last Login', path: 'DaysSinceLastLogin', type: 'number' },
			{ name: 'Sessions This Week', path: 'SessionsThisWeek', type: 'number' },
			{ name: 'Days Since Last Activity', path: 'DaysSinceLastActivity', type: 'number' },
			{ name: 'Is Active User', path: 'IsActiveUser', type: 'boolean' },
			{
				name: 'Engagement Level',
				path: 'EngagementLevel',
				type: 'string',
				options: [
					{ value: 'inactive', label: 'Inactive' },
					{ value: 'low', label: 'Low' },
					{ value: 'medium', label: 'Medium' },
					{ value: 'high', label: 'High' },
					{ value: 'power_user', label: 'Power User' }
				]
			}
		]
	},
	{
		name: 'Order',
		description: 'Order and purchase information',
		fields: [
			{
				name: 'Status',
				path: 'Order.Status',
				type: 'string',
				options: [
					{ value: 'pending', label: 'Pending' },
					{ value: 'processing', label: 'Processing' },
					{ value: 'completed', label: 'Completed' },
					{ value: 'cancelled', label: 'Cancelled' },
					{ value: 'refunded', label: 'Refunded' }
				]
			},
			{ name: 'Total Amount', path: 'Order.TotalAmount', type: 'number' },
			{ name: 'Is First Order', path: 'IsFirstOrder', type: 'boolean' },
			{ name: 'Has Discount Applied', path: 'HasDiscountApplied', type: 'boolean' },
			{ name: 'Refund Requested', path: 'RefundRequested', type: 'boolean' }
		]
	},
	{
		name: 'Dispute',
		description: 'Dispute and chargeback information',
		fields: [
			{ name: 'Has Open Dispute', path: 'HasOpenDispute', type: 'boolean' },
			{
				name: 'Dispute Status',
				path: 'Dispute.Status',
				type: 'string',
				options: [
					{ value: 'needs_response', label: 'Needs Response' },
					{ value: 'under_review', label: 'Under Review' },
					{ value: 'won', label: 'Won' },
					{ value: 'lost', label: 'Lost' }
				]
			},
			{
				name: 'Dispute Reason',
				path: 'Dispute.Reason',
				type: 'string',
				options: [
					{ value: 'fraudulent', label: 'Fraudulent' },
					{ value: 'duplicate', label: 'Duplicate' },
					{ value: 'product_not_received', label: 'Product Not Received' },
					{ value: 'product_unacceptable', label: 'Product Unacceptable' },
					{ value: 'subscription_canceled', label: 'Subscription Canceled' },
					{ value: 'unrecognized', label: 'Unrecognized' },
					{ value: 'other', label: 'Other' }
				]
			},
			{ name: 'Dispute Amount', path: 'Dispute.Amount', type: 'number' }
		]
	},
	{
		name: 'Account',
		description: 'Account-level information',
		fields: [
			{ name: 'Account Age (Days)', path: 'AccountAgeDays', type: 'number' },
			{ name: 'Total Orders', path: 'TotalOrders', type: 'number' },
			{ name: 'Failed Payment Count', path: 'FailedPaymentCount', type: 'number' },
			{ name: 'Has Payment Method', path: 'HasPaymentMethod', type: 'boolean' },
			{ name: 'Is VIP', path: 'IsVip', type: 'boolean' },
			{
				name: 'Risk Level',
				path: 'RiskLevel',
				type: 'string',
				options: [
					{ value: 'low', label: 'Low' },
					{ value: 'medium', label: 'Medium' },
					{ value: 'high', label: 'High' }
				]
			}
		]
	}
];

// Helper to get all fields as a flat list
export function getAllFactFields(): FactField[] {
	return FACT_CATEGORIES.flatMap((cat) => cat.fields);
}

// Helper to find a field by path
export function getFactFieldByPath(path: string): FactField | undefined {
	return getAllFactFields().find((f) => f.path === path);
}

// =============================================================================
// OPERATORS - Comparison operators for conditions
// =============================================================================

export interface Operator {
	value: string; // Grule operator
	label: string; // Display label
	types: ('string' | 'number' | 'boolean' | 'date')[]; // Compatible types
}

export const OPERATORS: Operator[] = [
	{ value: '==', label: 'equals', types: ['string', 'number', 'boolean'] },
	{ value: '!=', label: 'does not equal', types: ['string', 'number', 'boolean'] },
	{ value: '>', label: 'greater than', types: ['number', 'date'] },
	{ value: '>=', label: 'greater than or equal', types: ['number', 'date'] },
	{ value: '<', label: 'less than', types: ['number', 'date'] },
	{ value: '<=', label: 'less than or equal', types: ['number', 'date'] },
	{ value: 'contains', label: 'contains', types: ['string'] },
	{ value: 'startsWith', label: 'starts with', types: ['string'] },
	{ value: 'endsWith', label: 'ends with', types: ['string'] }
];

// Get operators compatible with a field type
export function getOperatorsForType(
	type: 'string' | 'number' | 'boolean' | 'date'
): Operator[] {
	return OPERATORS.filter((op) => op.types.includes(type));
}

// =============================================================================
// ACTIONS - Available actions that rules can trigger
// =============================================================================

export interface ActionParam {
	name: string;
	key: string;
	type: 'string' | 'number' | 'boolean' | 'select' | 'email_template' | 'sequence';
	required: boolean;
	options?: { value: string; label: string }[];
	placeholder?: string;
	description?: string;
}

export interface ActionType {
	type: string;
	name: string;
	description: string;
	category: 'communication' | 'ticket' | 'customer' | 'subscription' | 'system';
	params: ActionParam[];
}

export const ACTION_TYPES: ActionType[] = [
	{
		type: 'escalate_ticket',
		name: 'Escalate Ticket',
		description: 'Escalate a support ticket to higher priority or team',
		category: 'ticket',
		params: [
			{
				name: 'Priority',
				key: 'priority',
				type: 'select',
				required: false,
				options: [
					{ value: 'high', label: 'High' },
					{ value: 'urgent', label: 'Urgent' }
				]
			},
			{
				name: 'Team',
				key: 'team',
				type: 'string',
				required: false,
				placeholder: 'e.g., billing, engineering'
			}
		]
	},
	{
		type: 'send_email',
		name: 'Send Email',
		description: 'Send an email to the customer',
		category: 'communication',
		params: [
			{ name: 'To', key: 'to', type: 'string', required: true, placeholder: 'Customer.Email' },
			{
				name: 'Template',
				key: 'template_id',
				type: 'email_template',
				required: false,
				description: 'Email template to use'
			},
			{ name: 'Subject', key: 'subject', type: 'string', required: false }
		]
	},
	{
		type: 'start_sequence',
		name: 'Start Sequence',
		description: 'Enroll contact in an email sequence',
		category: 'communication',
		params: [
			{
				name: 'Sequence',
				key: 'sequence_id',
				type: 'sequence',
				required: true,
				description: 'Email sequence to start'
			}
		]
	},
	{
		type: 'notify_admin',
		name: 'Notify Admin',
		description: 'Send notification to admin users',
		category: 'communication',
		params: [
			{
				name: 'Message',
				key: 'message',
				type: 'string',
				required: true,
				placeholder: 'Notification message'
			},
			{
				name: 'Channel',
				key: 'channel',
				type: 'select',
				required: false,
				options: [
					{ value: 'email', label: 'Email' },
					{ value: 'slack', label: 'Slack' },
					{ value: 'dashboard', label: 'Dashboard' }
				]
			},
			{
				name: 'Urgency',
				key: 'urgency',
				type: 'select',
				required: false,
				options: [
					{ value: 'low', label: 'Low' },
					{ value: 'normal', label: 'Normal' },
					{ value: 'high', label: 'High' }
				]
			}
		]
	},
	{
		type: 'tag_customer',
		name: 'Tag Customer',
		description: 'Add or remove a tag from customer',
		category: 'customer',
		params: [
			{ name: 'Tag', key: 'tag', type: 'string', required: true, placeholder: 'e.g., vip, churn-risk' },
			{
				name: 'Action',
				key: 'action',
				type: 'select',
				required: false,
				options: [
					{ value: 'add', label: 'Add tag' },
					{ value: 'remove', label: 'Remove tag' }
				]
			}
		]
	},
	{
		type: 'pause_subscription',
		name: 'Pause Subscription',
		description: "Pause customer's subscription",
		category: 'subscription',
		params: [
			{
				name: 'Reason',
				key: 'reason',
				type: 'string',
				required: false,
				placeholder: 'Reason for pausing'
			}
		]
	},
	{
		type: 'create_task',
		name: 'Create Task',
		description: 'Create a follow-up task',
		category: 'system',
		params: [
			{ name: 'Title', key: 'title', type: 'string', required: true },
			{ name: 'Assignee', key: 'assignee', type: 'string', required: false },
			{
				name: 'Priority',
				key: 'priority',
				type: 'select',
				required: false,
				options: [
					{ value: 'low', label: 'Low' },
					{ value: 'normal', label: 'Normal' },
					{ value: 'high', label: 'High' }
				]
			},
			{ name: 'Due Date', key: 'due_date', type: 'string', required: false }
		]
	},
	{
		type: 'log_event',
		name: 'Log Event',
		description: 'Log an event for audit/analytics',
		category: 'system',
		params: [
			{
				name: 'Event Type',
				key: 'event_type',
				type: 'string',
				required: true,
				placeholder: 'e.g., rule_fired'
			},
			{ name: 'Message', key: 'message', type: 'string', required: false }
		]
	},
	{
		type: 'set_field',
		name: 'Set Field',
		description: 'Modify a field on an entity',
		category: 'system',
		params: [
			{
				name: 'Entity Type',
				key: 'entity_type',
				type: 'select',
				required: true,
				options: [
					{ value: 'ticket', label: 'Ticket' },
					{ value: 'customer', label: 'Customer' },
					{ value: 'contact', label: 'Contact' }
				]
			},
			{ name: 'Field', key: 'field', type: 'string', required: true, placeholder: 'Field name' },
			{ name: 'Value', key: 'value', type: 'string', required: true, placeholder: 'New value' }
		]
	},
	// Sequence Management
	{
		type: 'stop_sequence',
		name: 'Stop Sequence',
		description: 'Remove contact from an email sequence',
		category: 'communication',
		params: [
			{
				name: 'Sequence',
				key: 'sequence_id',
				type: 'sequence',
				required: false,
				description: 'Leave empty to stop all sequences'
			}
		]
	},
	{
		type: 'switch_sequence',
		name: 'Switch Sequence',
		description: 'Move contact from one sequence to another',
		category: 'communication',
		params: [
			{
				name: 'From Sequence',
				key: 'from_sequence_id',
				type: 'sequence',
				required: false,
				description: 'Current sequence (leave empty for any)'
			},
			{
				name: 'To Sequence',
				key: 'to_sequence_id',
				type: 'sequence',
				required: true,
				description: 'New sequence to enroll in'
			}
		]
	},
	// Billing Actions
	{
		type: 'apply_discount',
		name: 'Apply Discount',
		description: 'Apply a discount or coupon to the customer',
		category: 'subscription',
		params: [
			{
				name: 'Discount Type',
				key: 'discount_type',
				type: 'select',
				required: true,
				options: [
					{ value: 'percent', label: 'Percentage' },
					{ value: 'fixed', label: 'Fixed Amount' }
				]
			},
			{
				name: 'Amount',
				key: 'amount',
				type: 'number',
				required: true,
				placeholder: '10 for 10% or $10'
			},
			{
				name: 'Duration',
				key: 'duration',
				type: 'select',
				required: false,
				options: [
					{ value: 'once', label: 'Once' },
					{ value: 'forever', label: 'Forever' },
					{ value: 'repeating', label: 'Multiple Months' }
				]
			},
			{
				name: 'Months',
				key: 'duration_months',
				type: 'number',
				required: false,
				placeholder: 'Number of months (if repeating)'
			}
		]
	},
	{
		type: 'extend_trial',
		name: 'Extend Trial',
		description: 'Extend the trial period for a customer',
		category: 'subscription',
		params: [
			{
				name: 'Days',
				key: 'days',
				type: 'number',
				required: true,
				placeholder: 'Number of days to extend'
			},
			{
				name: 'Reason',
				key: 'reason',
				type: 'string',
				required: false,
				placeholder: 'Reason for extension'
			}
		]
	},
	{
		type: 'cancel_subscription',
		name: 'Cancel Subscription',
		description: 'Cancel the customer subscription',
		category: 'subscription',
		params: [
			{
				name: 'When',
				key: 'cancel_at',
				type: 'select',
				required: true,
				options: [
					{ value: 'immediately', label: 'Immediately' },
					{ value: 'period_end', label: 'At Period End' }
				]
			},
			{
				name: 'Reason',
				key: 'reason',
				type: 'string',
				required: false,
				placeholder: 'Cancellation reason'
			}
		]
	},
	{
		type: 'change_plan',
		name: 'Change Plan',
		description: 'Upgrade or downgrade customer plan',
		category: 'subscription',
		params: [
			{
				name: 'New Plan',
				key: 'plan',
				type: 'select',
				required: true,
				options: [
					{ value: 'free', label: 'Free' },
					{ value: 'pro', label: 'Pro' },
					{ value: 'team', label: 'Team' },
					{ value: 'enterprise', label: 'Enterprise' }
				]
			},
			{
				name: 'Prorate',
				key: 'prorate',
				type: 'select',
				required: false,
				options: [
					{ value: 'true', label: 'Yes' },
					{ value: 'false', label: 'No' }
				]
			}
		]
	},
	// Account Actions
	{
		type: 'suspend_account',
		name: 'Suspend Account',
		description: 'Temporarily suspend the customer account',
		category: 'customer',
		params: [
			{
				name: 'Reason',
				key: 'reason',
				type: 'select',
				required: true,
				options: [
					{ value: 'non_payment', label: 'Non-Payment' },
					{ value: 'abuse', label: 'Terms Violation' },
					{ value: 'fraud', label: 'Fraud Suspected' },
					{ value: 'other', label: 'Other' }
				]
			},
			{
				name: 'Message',
				key: 'message',
				type: 'string',
				required: false,
				placeholder: 'Additional details'
			}
		]
	},
	{
		type: 'flag_for_review',
		name: 'Flag for Review',
		description: 'Flag customer for manual review',
		category: 'customer',
		params: [
			{
				name: 'Flag Type',
				key: 'flag_type',
				type: 'select',
				required: true,
				options: [
					{ value: 'churn_risk', label: 'Churn Risk' },
					{ value: 'upgrade_opportunity', label: 'Upgrade Opportunity' },
					{ value: 'fraud_risk', label: 'Fraud Risk' },
					{ value: 'vip', label: 'VIP Candidate' },
					{ value: 'support_needed', label: 'Needs Support' }
				]
			},
			{
				name: 'Notes',
				key: 'notes',
				type: 'string',
				required: false,
				placeholder: 'Review notes'
			}
		]
	},
	{
		type: 'update_risk_level',
		name: 'Update Risk Level',
		description: 'Update customer risk assessment',
		category: 'customer',
		params: [
			{
				name: 'Risk Level',
				key: 'risk_level',
				type: 'select',
				required: true,
				options: [
					{ value: 'low', label: 'Low' },
					{ value: 'medium', label: 'Medium' },
					{ value: 'high', label: 'High' }
				]
			}
		]
	},
	// Ticket Actions
	{
		type: 'assign_ticket',
		name: 'Assign Ticket',
		description: 'Assign ticket to a team or agent',
		category: 'ticket',
		params: [
			{
				name: 'Assignee',
				key: 'assignee',
				type: 'string',
				required: false,
				placeholder: 'Agent email or ID'
			},
			{
				name: 'Team',
				key: 'team',
				type: 'select',
				required: false,
				options: [
					{ value: 'support', label: 'Support' },
					{ value: 'billing', label: 'Billing' },
					{ value: 'technical', label: 'Technical' },
					{ value: 'sales', label: 'Sales' }
				]
			}
		]
	},
	{
		type: 'close_ticket',
		name: 'Close Ticket',
		description: 'Close the support ticket',
		category: 'ticket',
		params: [
			{
				name: 'Resolution',
				key: 'resolution',
				type: 'select',
				required: true,
				options: [
					{ value: 'resolved', label: 'Resolved' },
					{ value: 'no_response', label: 'No Response' },
					{ value: 'duplicate', label: 'Duplicate' },
					{ value: 'wont_fix', label: "Won't Fix" }
				]
			},
			{
				name: 'Comment',
				key: 'comment',
				type: 'string',
				required: false,
				placeholder: 'Closing comment'
			}
		]
	},
	{
		type: 'add_ticket_note',
		name: 'Add Ticket Note',
		description: 'Add an internal note to the ticket',
		category: 'ticket',
		params: [
			{
				name: 'Note',
				key: 'note',
				type: 'string',
				required: true,
				placeholder: 'Internal note content'
			},
			{
				name: 'Visibility',
				key: 'visibility',
				type: 'select',
				required: false,
				options: [
					{ value: 'internal', label: 'Internal Only' },
					{ value: 'public', label: 'Visible to Customer' }
				]
			}
		]
	},
	// Integration Actions
	{
		type: 'trigger_webhook',
		name: 'Trigger Webhook',
		description: 'Send data to an external webhook URL',
		category: 'system',
		params: [
			{
				name: 'Webhook URL',
				key: 'url',
				type: 'string',
				required: true,
				placeholder: 'https://...'
			},
			{
				name: 'Include Data',
				key: 'include_data',
				type: 'select',
				required: false,
				options: [
					{ value: 'all', label: 'All Facts' },
					{ value: 'customer', label: 'Customer Only' },
					{ value: 'event', label: 'Event Only' }
				]
			}
		]
	},
	{
		type: 'sync_to_crm',
		name: 'Sync to CRM',
		description: 'Sync customer data to external CRM',
		category: 'system',
		params: [
			{
				name: 'CRM',
				key: 'crm',
				type: 'select',
				required: true,
				options: [
					{ value: 'hubspot', label: 'HubSpot' },
					{ value: 'salesforce', label: 'Salesforce' },
					{ value: 'pipedrive', label: 'Pipedrive' }
				]
			},
			{
				name: 'Action',
				key: 'action',
				type: 'select',
				required: true,
				options: [
					{ value: 'create', label: 'Create/Update Contact' },
					{ value: 'update_stage', label: 'Update Deal Stage' },
					{ value: 'add_note', label: 'Add Note' }
				]
			}
		]
	},
	// In-App Actions
	{
		type: 'send_in_app_notification',
		name: 'In-App Notification',
		description: 'Send notification within the application',
		category: 'communication',
		params: [
			{
				name: 'Title',
				key: 'title',
				type: 'string',
				required: true,
				placeholder: 'Notification title'
			},
			{
				name: 'Message',
				key: 'message',
				type: 'string',
				required: true,
				placeholder: 'Notification message'
			},
			{
				name: 'Type',
				key: 'notification_type',
				type: 'select',
				required: false,
				options: [
					{ value: 'info', label: 'Info' },
					{ value: 'warning', label: 'Warning' },
					{ value: 'success', label: 'Success' },
					{ value: 'error', label: 'Error' }
				]
			},
			{
				name: 'Action URL',
				key: 'action_url',
				type: 'string',
				required: false,
				placeholder: 'Link when clicked'
			}
		]
	},
	// Refund Actions
	{
		type: 'process_refund',
		name: 'Process Refund',
		description: 'Initiate a refund for the customer',
		category: 'subscription',
		params: [
			{
				name: 'Amount',
				key: 'amount',
				type: 'select',
				required: true,
				options: [
					{ value: 'full', label: 'Full Refund' },
					{ value: 'partial', label: 'Partial Refund' },
					{ value: 'prorated', label: 'Prorated' }
				]
			},
			{
				name: 'Partial Amount',
				key: 'partial_amount',
				type: 'number',
				required: false,
				placeholder: 'Amount if partial'
			},
			{
				name: 'Reason',
				key: 'reason',
				type: 'string',
				required: false,
				placeholder: 'Refund reason'
			}
		]
	}
];

// Get action type by name
export function getActionType(type: string): ActionType | undefined {
	return ACTION_TYPES.find((a) => a.type === type);
}

// Group actions by category
export function getActionsByCategory(): Record<string, ActionType[]> {
	return ACTION_TYPES.reduce(
		(acc, action) => {
			if (!acc[action.category]) {
				acc[action.category] = [];
			}
			acc[action.category].push(action);
			return acc;
		},
		{} as Record<string, ActionType[]>
	);
}

// =============================================================================
// RULE CONDITION - A single condition in a rule
// =============================================================================

export interface RuleCondition {
	id: string;
	field: string; // Fact path
	operator: string; // Comparison operator
	value: string | number | boolean; // Value to compare against
}

export interface ConditionGroup {
	id: string;
	logic: 'AND' | 'OR';
	conditions: RuleCondition[];
}

// =============================================================================
// RULE ACTION - A single action in a rule
// =============================================================================

export interface RuleActionConfig {
	id: string;
	type: string;
	params: Record<string, string | number | boolean>;
}

// =============================================================================
// RULE DEFINITION - Complete rule structure
// =============================================================================

export interface RuleDefinition {
	name: string;
	description?: string;
	salience: number;
	conditionGroups: ConditionGroup[]; // Multiple groups joined by AND
	actions: RuleActionConfig[];
}

// =============================================================================
// CONVERSION UTILITIES
// =============================================================================

// Convert visual rule definition to Grule JSON
export function ruleDefinitionToGruleJson(rule: RuleDefinition): string {
	const whenClause = buildWhenClause(rule.conditionGroups);
	const thenActions = rule.actions.map((a) => ({
		action: a.type,
		params: a.params
	}));

	const gruleRule = {
		name: rule.name,
		description: rule.description || '',
		salience: rule.salience,
		when: whenClause,
		then: thenActions
	};

	return JSON.stringify(gruleRule, null, 2);
}

function buildWhenClause(groups: ConditionGroup[]): string {
	if (groups.length === 0) return 'true';

	const groupClauses = groups.map((group) => {
		if (group.conditions.length === 0) return 'true';

		const conditionClauses = group.conditions.map((cond) => {
			const field = getFactFieldByPath(cond.field);
			const value = formatValue(cond.value, field?.type || 'string');

			// Handle special operators
			if (cond.operator === 'contains') {
				return `Facts.Contains(${cond.field}, ${value})`;
			}
			if (cond.operator === 'startsWith') {
				return `Facts.StartsWith(${cond.field}, ${value})`;
			}
			if (cond.operator === 'endsWith') {
				return `Facts.EndsWith(${cond.field}, ${value})`;
			}

			return `${cond.field} ${cond.operator} ${value}`;
		});

		const joinOp = group.logic === 'AND' ? ' && ' : ' || ';
		return group.conditions.length > 1
			? `(${conditionClauses.join(joinOp)})`
			: conditionClauses[0];
	});

	return groupClauses.join(' && ');
}

function formatValue(value: string | number | boolean, type: string): string {
	if (type === 'string') {
		return `"${String(value).replace(/"/g, '\\"')}"`;
	}
	if (type === 'boolean') {
		return String(value);
	}
	return String(value);
}

// Parse Grule JSON back to visual rule definition (best effort)
export function gruleJsonToRuleDefinition(json: string): RuleDefinition | null {
	try {
		const parsed = JSON.parse(json);
		return {
			name: parsed.name || '',
			description: parsed.description || '',
			salience: parsed.salience || 0,
			conditionGroups: parseWhenClause(parsed.when || 'true'),
			actions: (parsed.then || []).map(
				(a: { action: string; params: Record<string, unknown> }, i: number) => ({
					id: `action-${i}`,
					type: a.action,
					params: a.params || {}
				})
			)
		};
	} catch {
		return null;
	}
}

function parseWhenClause(when: string): ConditionGroup[] {
	// Simple parsing - just create one group with the raw condition
	// More sophisticated parsing could be added later
	if (when === 'true' || !when) {
		return [{ id: 'group-0', logic: 'AND', conditions: [] }];
	}

	// For now, just show the raw condition and let users edit via JSON
	return [
		{
			id: 'group-0',
			logic: 'AND',
			conditions: [
				{
					id: 'cond-0',
					field: '',
					operator: '==',
					value: when // Store raw for display
				}
			]
		}
	];
}

// Generate a unique ID
export function generateId(): string {
	return Math.random().toString(36).substring(2, 11);
}

// =============================================================================
// HUMAN-READABLE SUMMARIES
// =============================================================================

// Get human-readable label for a field path
export function getFieldLabel(path: string): string {
	const field = getFactFieldByPath(path);
	if (field) {
		// Find the category for context
		for (const cat of FACT_CATEGORIES) {
			if (cat.fields.includes(field)) {
				return `${cat.name} ${field.name}`;
			}
		}
		return field.name;
	}
	// Fallback: convert path to readable form
	return path.replace(/\./g, ' ').replace(/([A-Z])/g, ' $1').trim();
}

// Get human-readable label for an operator
export function getOperatorLabel(op: string): string {
	const operator = OPERATORS.find((o) => o.value === op);
	return operator?.label || op;
}

// Get human-readable value display
export function getValueLabel(value: string | number | boolean, fieldPath: string): string {
	const field = getFactFieldByPath(fieldPath);
	if (field?.options) {
		const option = field.options.find((o) => o.value === String(value));
		if (option) return option.label;
	}
	if (typeof value === 'boolean') {
		return value ? 'Yes' : 'No';
	}
	return String(value);
}

// Generate human-readable summary for a single condition
export function conditionToText(condition: RuleCondition): string {
	if (!condition.field) return '';

	const field = getFieldLabel(condition.field);
	const op = getOperatorLabel(condition.operator);
	const val = getValueLabel(condition.value, condition.field);

	return `${field} ${op} "${val}"`;
}

// Generate human-readable summary for a condition group
export function conditionGroupToText(group: ConditionGroup): string {
	if (group.conditions.length === 0) return 'Always';

	const conditionTexts = group.conditions
		.filter((c) => c.field)
		.map((c) => conditionToText(c));

	if (conditionTexts.length === 0) return 'Always';

	const joiner = group.logic === 'AND' ? ' AND ' : ' OR ';
	return conditionTexts.join(joiner);
}

// Generate full human-readable summary for all condition groups
export function conditionGroupsToText(groups: ConditionGroup[]): string {
	const groupTexts = groups
		.map((g) => conditionGroupToText(g))
		.filter((t) => t && t !== 'Always');

	if (groupTexts.length === 0) return 'Always matches';
	if (groupTexts.length === 1) return groupTexts[0];

	return groupTexts.map((t) => `(${t})`).join(' AND ');
}

// Generate human-readable summary for an action
export function actionToText(action: RuleActionConfig): string {
	const actionDef = getActionType(action.type);
	if (!actionDef) return action.type;

	let text = actionDef.name;

	// Add key params to summary
	const keyParams = Object.entries(action.params)
		.filter(([, v]) => v !== '' && v !== undefined)
		.slice(0, 2); // Show first 2 params max

	if (keyParams.length > 0) {
		const paramTexts = keyParams.map(([, v]) => String(v));
		text += `: ${paramTexts.join(', ')}`;
	}

	return text;
}

// Generate summary of all actions
export function actionsToText(actions: RuleActionConfig[]): string {
	if (actions.length === 0) return 'No actions';
	return actions.map((a) => actionToText(a)).join('; ');
}
