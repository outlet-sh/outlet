<!--
  StatusBadge Component
  Context-aware status indicators with tooltips
-->

<script lang="ts">
	import Tooltip from './Tooltip.svelte';

	// Status config based on backend shared_models.py status definitions
	// Using semantic colors for light theme
	const statusConfig = {
		// Collection statuses
		"ACTIVE": {
			color: "bg-success",
			textColor: "text-success",
			bgColor: "bg-success/10",
			label: "Active"
		},
		"active": {
			color: "bg-success",
			textColor: "text-success",
			bgColor: "bg-success/10",
			label: "Active"
		},
		"ERROR": {
			color: "bg-error",
			textColor: "text-error",
			bgColor: "bg-error/10",
			label: "Error"
		},
		"error": {
			color: "bg-error",
			textColor: "text-error",
			bgColor: "bg-error/10",
			label: "Error"
		},
		"NEEDS SOURCE": {
			color: "bg-text-muted",
			textColor: "text-text-secondary",
			bgColor: "bg-bg-secondary",
			label: "Needs Source"
		},
		"needs source": {
			color: "bg-text-muted",
			textColor: "text-text-secondary",
			bgColor: "bg-bg-secondary",
			label: "Needs Source"
		},
		// Source connection statuses
		"IN_PROGRESS": {
			color: "bg-info",
			textColor: "text-info",
			bgColor: "bg-info/10",
			label: "Syncing"
		},
		"in_progress": {
			color: "bg-info",
			textColor: "text-info",
			bgColor: "bg-info/10",
			label: "In Progress"
		},
		"failing": {
			color: "bg-error",
			textColor: "text-error",
			bgColor: "bg-error/10",
			label: "Failing"
		},
		// Sync job statuses
		"pending": {
			color: "bg-warning",
			textColor: "text-warning",
			bgColor: "bg-warning/10",
			label: "Pending"
		},
		"completed": {
			color: "bg-success",
			textColor: "text-success",
			bgColor: "bg-success/10",
			label: "Completed"
		},
		"failed": {
			color: "bg-error",
			textColor: "text-error",
			bgColor: "bg-error/10",
			label: "Failed"
		},
		"cancelled": {
			color: "bg-error",
			textColor: "text-error",
			bgColor: "bg-error/10",
			label: "Cancelled"
		},
		// API Key statuses
		"EXPIRED": {
			color: "bg-error",
			textColor: "text-error",
			bgColor: "bg-error/10",
			label: "Expired"
		},
		"EXPIRING_SOON": {
			color: "bg-warning",
			textColor: "text-warning",
			bgColor: "bg-warning/10",
			label: "Expiring Soon"
		},
		"UNKNOWN": {
			color: "bg-text-muted",
			textColor: "text-text-secondary",
			bgColor: "bg-bg-secondary",
			label: "Unknown"
		},
		// Fallback for unknown statuses
		"default": {
			color: "bg-text-muted",
			textColor: "text-text-secondary",
			bgColor: "bg-bg-secondary",
			label: "Unknown"
		}
	};

	type TooltipContext = "collection" | "apiKey";

	interface Props {
		status: string;
		class?: string;
		showTooltip?: boolean;
		tooltipContext?: TooltipContext;
	}

	let {
		status,
		class: className = '',
		showTooltip = false,
		tooltipContext
	}: Props = $props();

	// Get status configuration or default
	function getStatusConfig(statusKey: string = "") {
		// Try exact match first
		if (statusKey in statusConfig) {
			return statusConfig[statusKey as keyof typeof statusConfig];
		}

		// Try case-insensitive match
		const lowerKey = statusKey.toLowerCase();
		for (const key in statusConfig) {
			if (key.toLowerCase() === lowerKey) {
				return statusConfig[key as keyof typeof statusConfig];
			}
		}

		// Return default if no match
		const formatted = statusKey ? statusKey.charAt(0).toUpperCase() + statusKey.slice(1).toLowerCase() : "Unknown";
		return {
			...statusConfig["default"],
			label: formatted
		};
	}

	// Context-aware status descriptions for tooltips
	function getStatusDescription(statusKey: string, context?: TooltipContext): string | null {
		const normalizedKey = statusKey.toUpperCase();

		// Collection-specific descriptions
		const collectionDescriptions: Record<string, string> = {
			"ACTIVE": "At least one source connection has completed a sync or is currently syncing. Your collection has data and is ready for queries.",
			"ERROR": "All source connections have failed their last sync. Check your connections and authentication to resolve sync issues.",
			"NEEDS SOURCE": "This collection has no authenticated connections, or connections exist but haven't successfully synced yet. Configure a source or wait for the initial sync to complete."
		};

		// API key-specific descriptions
		const apiKeyDescriptions: Record<string, string> = {
			"ACTIVE": "This API key is valid and can be used to authenticate requests. It will remain active until its expiration date.",
			"EXPIRING_SOON": "This API key will expire within 7 days. Create a new key before expiration to avoid service interruption.",
			"EXPIRED": "This API key has expired and can no longer be used. Delete this key and create a new one to continue using the API.",
			"UNKNOWN": "Unable to determine the status of this API key. Please verify the expiration date."
		};

		// Select description based on context
		if (context === "apiKey") {
			return apiKeyDescriptions[normalizedKey] || null;
		} else if (context === "collection") {
			return collectionDescriptions[normalizedKey] || null;
		}

		// No context specified - return null (no tooltip)
		return null;
	}

	const config = $derived(getStatusConfig(status));
	const description = $derived(getStatusDescription(status, tooltipContext));
</script>

{#if showTooltip && description}
	<Tooltip content={description}>
		<div class="inline-flex items-center gap-1.5 py-1 px-2.5 rounded-full {config.bgColor} {className}">
			<div class="h-2 w-2 rounded-full {config.color}"></div>
			<span class="text-xs font-medium {config.textColor}">
				{config.label}
			</span>
		</div>
	</Tooltip>
{:else}
	<div class="inline-flex items-center gap-1.5 py-1 px-2.5 rounded-full {config.bgColor} {className}">
		<div class="h-2 w-2 rounded-full {config.color}"></div>
		<span class="text-xs font-medium {config.textColor}">
			{config.label}
		</span>
	</div>
{/if}

<style>
	@reference "$src/app.css";
	@layer components.status-badge {
		/* Status badge uses utility classes */
	}
</style>
