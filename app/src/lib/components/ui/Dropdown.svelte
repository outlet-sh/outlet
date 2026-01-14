<!--
  Dropdown Menu Component
  Dropdown menu with items, dividers, and sections
-->

<script lang="ts">
	interface MenuItem {
		label: string;
		icon?: string;
		onclick?: () => void;
		href?: string;
		disabled?: boolean;
		danger?: boolean;
		divider?: boolean;
	}

	interface Props {
		items: MenuItem[];
		label?: string;
		icon?: string;
		align?: 'left' | 'right';
		children?: any;
	}

	let {
		items,
		label = 'Options',
		icon = 'ellipsis-v',
		align = 'right',
		children
	}: Props = $props();

	let isOpen = $state(false);

	function toggleDropdown() {
		isOpen = !isOpen;
	}

	function closeDropdown() {
		isOpen = false;
	}

	function handleItemClick(item: MenuItem) {
		if (item.disabled) return;
		if (item.onclick) {
			item.onclick();
		}
		closeDropdown();
	}
</script>

<svelte:window onclick={closeDropdown} />

<div class="relative inline-block text-left">
	<button
		type="button"
		onclick={(e) => {
			e.stopPropagation();
			toggleDropdown();
		}}
		class="btn-secondary inline-flex items-center gap-2"
	>
		{#if children}
			{@render children()}
		{:else}
			<i class="fas fa-{icon}"></i>
			{label}
		{/if}
	</button>

	{#if isOpen}
		<div
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.stopPropagation()}
			role="menu"
			tabindex="-1"
			class="dropdown-menu {align === 'right' ? 'dropdown-right' : 'dropdown-left'}"
		>
			<div class="dropdown-content">
				{#each items as item}
					{#if item.divider}
						<div class="dropdown-divider"></div>
					{:else if item.href}
						<a
							href={item.href}
							class="dropdown-item {item.disabled ? 'dropdown-item-disabled' : item.danger ? 'dropdown-item-danger' : ''}"
						>
							{#if item.icon}
								<i class="dropdown-item-icon fas fa-{item.icon}"></i>
							{/if}
							{item.label}
						</a>
					{:else}
						<button
							type="button"
							onclick={() => handleItemClick(item)}
							disabled={item.disabled}
							class="dropdown-item {item.disabled ? 'dropdown-item-disabled' : item.danger ? 'dropdown-item-danger' : ''}"
						>
							{#if item.icon}
								<i class="dropdown-item-icon fas fa-{item.icon}"></i>
							{/if}
							{item.label}
						</button>
					{/if}
				{/each}
			</div>
		</div>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.dropdown {
		.dropdown-menu {
			@apply absolute z-50 mt-2 w-56 rounded-xl shadow-xl focus:outline-none bg-bg;
			border: 1px solid var(--color-border);
		}

		.dropdown-right {
			@apply right-0;
		}

		.dropdown-left {
			@apply left-0;
		}

		.dropdown-content {
			@apply py-1;
		}

		.dropdown-divider {
			@apply my-1 h-px;
			background: var(--color-border);
		}

		.dropdown-item {
			@apply flex items-center gap-3 px-4 py-2.5 text-base transition-colors w-full text-left text-text-secondary;

			&:hover:not(.dropdown-item-disabled) {
				@apply bg-bg-secondary text-text;
			}
		}

		.dropdown-item-icon {
			@apply w-5;
		}

		.dropdown-item-disabled {
			@apply cursor-not-allowed opacity-50;
		}

		.dropdown-item-danger {
			color: var(--color-error);

			&:hover {
				background: color-mix(in srgb, var(--color-error) 10%, transparent);
			}
		}
	}
</style>
