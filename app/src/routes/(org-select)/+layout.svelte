<script lang="ts">
	import { page } from '$app/stores';
	import { getCurrentUser, logout } from '$lib/auth';
	import AuthGuard from '$lib/components/admin/AuthGuard.svelte';

	const { children } = $props();

	let user = $state(getCurrentUser());
	let showUserMenu = $state(false);

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

<svelte:window onclick={handleClickOutside} />

<svelte:head>
	<title>Outlet</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

<AuthGuard>
<div class="min-h-screen bg-gray-50">
	<!-- Minimal navbar for org selection / global settings -->
	<nav class="navbar">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex h-16 justify-between">
				<div class="flex items-center">
					<a href="/" class="navbar-brand">Outlet</a>
				</div>

				<div class="flex items-center space-x-4">
					<!-- Global Settings Link -->
					<a
						href="/settings/email"
						class="navbar-link {$page.url.pathname.startsWith('/settings') ? 'active' : ''}"
					>
						<svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
						</svg>
						Settings
					</a>

					<!-- User Menu -->
					<div class="user-menu-container relative ml-3">
						<div>
							<button
								onclick={(e) => { e.stopPropagation(); showUserMenu = !showUserMenu; }}
								class="navbar-user-btn"
								aria-expanded={showUserMenu}
								aria-haspopup="true"
							>
								<span class="sr-only">Open user menu</span>
								<div class="navbar-avatar">
									<span class="text-sm font-medium text-white">
										{(user?.name || user?.email || 'U')[0].toUpperCase()}
									</span>
								</div>
							</button>
						</div>

						{#if showUserMenu}
							<div class="navbar-dropdown" role="menu" aria-orientation="vertical">
								<div class="px-4 py-3 border-b border-border">
									<p class="text-sm font-medium text-text truncate">{user?.name || user?.email}</p>
									<p class="text-xs text-text-muted capitalize">{user?.role}</p>
								</div>
								<button
									onclick={handleLogout}
									class="navbar-dropdown-item w-full text-left"
									role="menuitem"
								>
									Sign out
								</button>
							</div>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</nav>

	<main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
		{@render children()}
	</main>
</div>
</AuthGuard>

<style>
	@reference "$src/app.css";
	@layer components.navbar {
		.navbar {
			@apply border-b shadow-sm;
			border-color: var(--color-border);
			background-color: var(--color-bg);
		}

		.navbar-brand {
			@apply text-xl font-bold;
			color: var(--color-primary);
		}

		.navbar-link {
			@apply inline-flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors;
			color: var(--color-text-muted);
		}

		.navbar-link:hover {
			background-color: var(--color-bg-secondary);
			color: var(--color-text);
		}

		.navbar-link.active {
			background-color: var(--color-bg-secondary);
			color: var(--color-text);
		}

		.navbar-user-btn {
			@apply flex items-center rounded-full text-sm transition-colors;
		}

		.navbar-user-btn:focus {
			@apply outline-none;
		}

		.navbar-avatar {
			@apply flex h-8 w-8 items-center justify-center rounded-full;
			background-color: var(--color-primary);
		}

		.navbar-dropdown {
			@apply absolute right-0 z-50 mt-2 w-48 origin-top-right rounded-md py-1 shadow-lg;
			background-color: var(--color-bg);
			box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
			border: 1px solid var(--color-border);
		}

		.navbar-dropdown-item {
			@apply block px-4 py-2 text-sm transition-colors;
			color: var(--color-text);
		}

		.navbar-dropdown-item:hover {
			background-color: var(--color-bg-secondary);
		}
	}
</style>
