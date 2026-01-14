<!--
  Badge Component
  For status indicators, counts, labels
-->

<script lang="ts">
	type BadgeVariant = 'default' | 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'info' | 'neutral';

	interface Props {
		variant?: BadgeVariant;
		type?: BadgeVariant; // alias for variant
		size?: 'sm' | 'md' | 'lg';
		rounded?: boolean;
		class?: string;
		style?: string;
		children: any;
	}

	let {
		variant = 'default',
		type,
		size = 'md',
		rounded = false,
		class: extraClass = '',
		style,
		children
	}: Props = $props();

	// type is alias for variant
	const effectiveVariant = $derived(type ?? variant);
</script>

<span
	class="badge badge-{effectiveVariant} badge-{size} {rounded ? 'badge-rounded' : 'badge-square'}"
	{style}
>
	{@render children()}
</span>

<style>
	@reference "$src/app.css";
	@layer components.badge {
		.badge {
			@apply inline-flex items-center gap-1 font-semibold;
		}

		.badge-rounded { @apply rounded-full; }
		.badge-square { @apply rounded-lg; }

		.badge-sm { @apply px-2 py-0.5 text-xs; }
		.badge-md { @apply px-2.5 py-1 text-sm; }
		.badge-lg { @apply px-3 py-1.5 text-base; }

		.badge-default {
			@apply bg-bg-secondary text-text-secondary;
			border: 1px solid var(--color-border);
		}

		.badge-primary {
			@apply text-white;
			background: linear-gradient(to right, var(--color-primary), var(--color-secondary));
		}

		.badge-success {
			background: color-mix(in srgb, var(--color-success) 15%, var(--color-bg));
			color: color-mix(in srgb, var(--color-success) 80%, black);
			@apply ring-1;
			ring-color: color-mix(in srgb, var(--color-success) 40%, transparent);
		}

		.badge-warning {
			background: color-mix(in srgb, var(--color-warning) 15%, var(--color-bg));
			color: color-mix(in srgb, var(--color-warning) 80%, black);
			@apply ring-1;
			ring-color: color-mix(in srgb, var(--color-warning) 40%, transparent);
		}

		.badge-error {
			background: color-mix(in srgb, var(--color-error) 15%, var(--color-bg));
			color: color-mix(in srgb, var(--color-error) 80%, black);
			@apply ring-1;
			ring-color: color-mix(in srgb, var(--color-error) 40%, transparent);
		}

		.badge-info {
			background: color-mix(in srgb, var(--color-info) 15%, var(--color-bg));
			color: color-mix(in srgb, var(--color-info) 80%, black);
			@apply ring-1;
			ring-color: color-mix(in srgb, var(--color-info) 40%, transparent);
		}
	}
</style>
