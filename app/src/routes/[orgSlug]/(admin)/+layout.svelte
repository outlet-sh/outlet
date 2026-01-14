<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Sidebar from '$lib/components/admin/Sidebar.svelte';
	import AuthGuard from '$lib/components/admin/AuthGuard.svelte';
	import { getOrganizationBySlug, type OrgInfo } from '$lib/api';
	import { getCurrentUser, logout } from '$lib/auth';
	import { Search, X, Menu } from 'lucide-svelte';

	interface Props {
		data: {
			orgSlug: string;
		};
		children: import('svelte').Snippet;
	}

	const { data, children }: Props = $props();

	let org = $state<OrgInfo | null>(null);
	let loading = $state(true);
	let error = $state('');
	let user = $state(getCurrentUser());

	// Mobile sidebar state
	let mobileMenuOpen = $state(false);

	// User menu state
	let showUserMenu = $state(false);

	$effect(() => {
		loadOrg();
	});

	// Close mobile menu on navigation
	$effect(() => {
		$page.url.pathname;
		mobileMenuOpen = false;
	});

	async function loadOrg() {
		loading = true;
		error = '';
		try {
			org = await getOrganizationBySlug({}, data.orgSlug);
			localStorage.setItem('currentOrgId', org.id);
			localStorage.setItem('currentOrgName', org.name);
			localStorage.setItem('currentOrgSlug', org.slug);
		} catch (err: any) {
			console.error('Failed to load organization:', err);
			error = `Organization "${data.orgSlug}" not found`;
			goto('/');
		} finally {
			loading = false;
		}
	}

	function switchOrg() {
		localStorage.removeItem('currentOrgId');
		localStorage.removeItem('currentOrgName');
		localStorage.removeItem('currentOrgSlug');
		goto('/');
	}

	function handleLogout() {
		logout();
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.user-menu-container')) {
			showUserMenu = false;
		}
	}
</script>

<svelte:head>
	<title>{org?.name || 'Admin'} | Outlet</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<svelte:window onclick={handleClickOutside} />

<AuthGuard>
	{#if loading}
		<div class="min-h-screen bg-bg-secondary flex items-center justify-center">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
		</div>
	{:else if error}
		<div class="min-h-screen bg-bg-secondary flex items-center justify-center">
			<div class="text-error">{error}</div>
		</div>
	{:else}
		<div class="admin-layout">
			<!-- Desktop Sidebar -->
			<div class="admin-sidebar hidden md:block">
				<Sidebar orgSlug={data.orgSlug} orgName={org?.name} />
			</div>

			<!-- Mobile Header -->
			<div class="admin-mobile-header md:hidden">
				<button
					onclick={() => mobileMenuOpen = !mobileMenuOpen}
					class="admin-mobile-menu-btn"
				>
					{#if mobileMenuOpen}
						<X class="h-6 w-6" />
					{:else}
						<Menu class="h-6 w-6" />
					{/if}
				</button>
				<span class="font-semibold text-text">{org?.name || 'Outlet'}</span>
				<div class="user-menu-container relative">
					<button
						onclick={(e) => { e.stopPropagation(); showUserMenu = !showUserMenu; }}
						class="admin-user-btn"
					>
						<div class="admin-avatar">
							{(user?.name || user?.email || 'U')[0].toUpperCase()}
						</div>
					</button>
					{#if showUserMenu}
						<div class="admin-user-dropdown">
							<div class="px-4 py-3 border-b border-border">
								<p class="text-sm font-medium text-text truncate">{user?.email}</p>
							</div>
							<a href="/{data.orgSlug}/settings" class="admin-dropdown-item" onclick={() => showUserMenu = false}>
								Settings
							</a>
							<div class="border-t border-border my-1"></div>
							<button onclick={switchOrg} class="admin-dropdown-item w-full text-left">
								Switch Site
							</button>
							<button onclick={handleLogout} class="admin-dropdown-item w-full text-left">
								Sign out
							</button>
						</div>
					{/if}
				</div>
			</div>

			<!-- Mobile Sidebar Overlay -->
			{#if mobileMenuOpen}
				<div
					class="admin-mobile-overlay md:hidden"
					onclick={() => mobileMenuOpen = false}
					onkeydown={(e) => e.key === 'Escape' && (mobileMenuOpen = false)}
					role="button"
					tabindex="-1"
				></div>
				<div class="admin-mobile-sidebar md:hidden">
					<Sidebar orgSlug={data.orgSlug} orgName={org?.name} />
				</div>
			{/if}

			<!-- Main Content -->
			<div class="admin-main">
				<!-- Top Bar (Desktop) -->
				<header class="admin-topbar hidden md:flex">
					<div class="flex-1"></div>
					<div class="flex items-center gap-4">
						<!-- User Menu -->
						<div class="user-menu-container relative">
							<button
								onclick={(e) => { e.stopPropagation(); showUserMenu = !showUserMenu; }}
								class="admin-user-btn"
							>
								<div class="admin-avatar">
									{(user?.name || user?.email || 'U')[0].toUpperCase()}
								</div>
								<span class="text-sm text-text-muted hidden lg:block">{user?.email}</span>
							</button>
							{#if showUserMenu}
								<div class="admin-user-dropdown">
									<div class="px-4 py-3 border-b border-border">
										<p class="text-sm font-medium text-text truncate">{user?.email}</p>
										<p class="text-xs text-text-muted capitalize">{user?.role || 'user'}</p>
									</div>
									<a href="/{data.orgSlug}/settings" class="admin-dropdown-item" onclick={() => showUserMenu = false}>
										Settings
									</a>
									<div class="border-t border-border my-1"></div>
									<button onclick={switchOrg} class="admin-dropdown-item w-full text-left">
										Switch Site
									</button>
									<button onclick={handleLogout} class="admin-dropdown-item w-full text-left">
										Sign out
									</button>
								</div>
							{/if}
						</div>
					</div>
				</header>

				<!-- Page Content -->
				<main class="admin-content">
					{@render children()}
				</main>
			</div>
		</div>
	{/if}
</AuthGuard>

<style>
	@reference "$src/app.css";
	@layer components.admin-layout {
		.admin-layout {
			@apply flex min-h-screen;
			background-color: var(--color-bg-secondary);
		}

		.admin-sidebar {
			@apply flex-shrink-0 fixed inset-y-0 left-0 z-30;
		}

		.admin-main {
			@apply flex-1 flex flex-col min-h-screen;
			margin-left: 220px;
		}

		/* Adjust for collapsed sidebar */
		:global(.sidebar.collapsed) ~ .admin-main {
			margin-left: 64px;
		}

		.admin-topbar {
			@apply h-14 px-6 items-center justify-between border-b sticky top-0 z-20;
			background-color: var(--color-bg);
			border-color: var(--color-border);
		}

		.admin-content {
			@apply flex-1;
		}

		/* Mobile styles */
		.admin-mobile-header {
			@apply fixed top-0 left-0 right-0 h-14 px-4 flex items-center justify-between z-40 border-b;
			background-color: var(--color-bg);
			border-color: var(--color-border);
		}

		.admin-mobile-menu-btn {
			@apply p-2 rounded-lg transition-colors;
			color: var(--color-text-muted);
		}

		.admin-mobile-menu-btn:hover {
			background-color: var(--color-bg-secondary);
		}

		.admin-mobile-overlay {
			@apply fixed inset-0 z-40;
			background-color: rgba(0, 0, 0, 0.5);
		}

		.admin-mobile-sidebar {
			@apply fixed inset-y-0 left-0 z-50 w-[220px];
		}

		@media (max-width: 767px) {
			.admin-main {
				margin-left: 0;
				padding-top: 56px;
			}
		}

		/* User menu */
		.admin-user-btn {
			@apply flex items-center gap-2 p-1 rounded-lg transition-colors;
		}

		.admin-user-btn:hover {
			background-color: var(--color-bg-secondary);
		}

		.admin-avatar {
			@apply flex items-center justify-center h-8 w-8 rounded-full text-sm font-medium text-white;
			background-color: var(--color-primary);
		}

		.admin-user-dropdown {
			@apply absolute right-0 top-full mt-2 w-56 rounded-lg shadow-lg py-1 z-50;
			background-color: var(--color-bg);
			border: 1px solid var(--color-border);
		}

		.admin-dropdown-item {
			@apply block px-4 py-2 text-sm transition-colors;
			color: var(--color-text);
		}

		.admin-dropdown-item:hover {
			background-color: var(--color-bg-secondary);
		}
	}
</style>
