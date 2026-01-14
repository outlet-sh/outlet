<script lang="ts">
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import {
		LayoutDashboard,
		Users,
		Send,
		Workflow,
		FileText,
		Settings,
		ChevronLeft,
		ChevronRight,
		Mail
	} from 'lucide-svelte';

	interface Props {
		orgSlug: string;
		orgName?: string;
	}

	const { orgSlug, orgName = 'Outlet' }: Props = $props();

	const STORAGE_KEY = 'sidebar-collapsed';

	// Load collapsed state from localStorage
	let collapsed = $state(false);

	$effect(() => {
		if (browser) {
			const saved = localStorage.getItem(STORAGE_KEY);
			collapsed = saved === 'true';
		}
	});

	function toggleCollapse() {
		collapsed = !collapsed;
		if (browser) {
			localStorage.setItem(STORAGE_KEY, String(collapsed));
		}
	}

	let basePath = $derived(`/${orgSlug}`);

	// Navigation items - flat, email-focused like Sendy
	const navItems = [
		{ id: 'dashboard', label: 'Dashboard', href: '', icon: LayoutDashboard, exact: true },
		{ id: 'lists', label: 'Lists', href: '/lists', icon: Users },
		{ id: 'campaigns', label: 'Campaigns', href: '/campaigns', icon: Send },
		{ id: 'sequences', label: 'Sequences', href: '/sequences', icon: Workflow },
		{ id: 'templates', label: 'Templates', href: '/templates', icon: FileText },
	];

	const bottomItems = [
		{ id: 'settings', label: 'Settings', href: '/settings', icon: Settings },
	];

	function isActive(href: string, exact = false): boolean {
		const fullPath = `${basePath}${href}`;
		if (exact) {
			return $page.url.pathname === fullPath || $page.url.pathname === fullPath + '/';
		}
		return $page.url.pathname.startsWith(fullPath) && fullPath !== basePath;
	}
</script>

<aside class="sidebar" class:collapsed>
	<!-- Logo / Brand -->
	<div class="sidebar-header">
		{#if !collapsed}
			<a href={basePath} class="sidebar-brand">
				<Mail class="h-6 w-6 text-primary" />
				<span class="sidebar-brand-text">{orgName}</span>
			</a>
		{:else}
			<a href={basePath} class="sidebar-brand-icon">
				<Mail class="h-6 w-6 text-primary" />
			</a>
		{/if}
	</div>

	<!-- Main Navigation -->
	<nav class="sidebar-nav">
		{#each navItems as item}
			<a
				href="{basePath}{item.href}"
				class="sidebar-link"
				class:active={isActive(item.href, item.exact)}
				title={collapsed ? item.label : undefined}
			>
				<svelte:component this={item.icon} class="sidebar-icon" />
				{#if !collapsed}
					<span class="sidebar-label">{item.label}</span>
				{/if}
			</a>
		{/each}
	</nav>

	<!-- Spacer -->
	<div class="flex-1"></div>

	<!-- Bottom Navigation -->
	<nav class="sidebar-nav sidebar-nav-bottom">
		{#each bottomItems as item}
			<a
				href="{basePath}{item.href}"
				class="sidebar-link"
				class:active={isActive(item.href)}
				title={collapsed ? item.label : undefined}
			>
				<svelte:component this={item.icon} class="sidebar-icon" />
				{#if !collapsed}
					<span class="sidebar-label">{item.label}</span>
				{/if}
			</a>
		{/each}

		<!-- Collapse Toggle -->
		<button
			onclick={toggleCollapse}
			class="sidebar-collapse-btn"
			title={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
		>
			{#if collapsed}
				<ChevronRight class="h-4 w-4" />
			{:else}
				<ChevronLeft class="h-4 w-4" />
				<span class="sidebar-label">Collapse</span>
			{/if}
		</button>
	</nav>
</aside>

<style>
	@reference "$src/app.css";
	@layer components.sidebar {
		.sidebar {
			@apply flex flex-col h-full border-r;
			width: 220px;
			background-color: var(--color-bg);
			border-color: var(--color-border);
			transition: width 0.2s ease;
		}

		.sidebar.collapsed {
			width: 64px;
		}

		.sidebar-header {
			@apply flex items-center h-16 px-4 border-b;
			border-color: var(--color-border);
		}

		.sidebar-brand {
			@apply flex items-center gap-3 font-semibold;
			color: var(--color-text);
		}

		.sidebar-brand-text {
			@apply truncate;
			max-width: 140px;
		}

		.sidebar-brand-icon {
			@apply flex items-center justify-center w-full;
		}

		.sidebar-nav {
			@apply flex flex-col gap-1 p-3;
		}

		.sidebar-nav-bottom {
			@apply border-t pt-3;
			border-color: var(--color-border);
		}

		.sidebar-link {
			@apply flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all;
			color: var(--color-text-muted);
		}

		.sidebar-link:hover {
			background-color: var(--color-bg-secondary);
			color: var(--color-text);
		}

		.sidebar-link.active {
			background-color: color-mix(in srgb, var(--color-primary) 12%, transparent);
			color: var(--color-primary);
		}

		.sidebar.collapsed .sidebar-link {
			@apply justify-center px-0;
		}

		.sidebar-icon {
			@apply h-5 w-5 flex-shrink-0;
		}

		.sidebar-label {
			@apply truncate;
		}

		.sidebar-collapse-btn {
			@apply flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all w-full;
			color: var(--color-text-muted);
		}

		.sidebar-collapse-btn:hover {
			background-color: var(--color-bg-secondary);
			color: var(--color-text);
		}

		.sidebar.collapsed .sidebar-collapse-btn {
			@apply justify-center px-0;
		}
	}
</style>
