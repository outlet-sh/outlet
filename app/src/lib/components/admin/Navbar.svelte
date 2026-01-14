<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getCurrentUser, isAdmin, logout } from '$lib/auth';
	import { adminListCustomers, listOrders, listBillingProducts, type CustomerInfo, type OrderInfo, type BillingProductInfo, type OrgInfo } from '$lib/api';
	import { Search, X } from 'lucide-svelte';

	interface Props {
		org?: OrgInfo | null;
		orgSlug?: string;
	}

	const { org, orgSlug }: Props = $props();

	let showUserMenu = $state(false);
	let showMobileMenu = $state(false);
	let showSearch = $state(false);
	let searchQuery = $state('');
	let searchResults = $state<{
		customers: CustomerInfo[];
		orders: OrderInfo[];
		products: BillingProductInfo[];
	}>({ customers: [], orders: [], products: [] });
	let searchLoading = $state(false);
	let searchTimeout: ReturnType<typeof setTimeout> | null = null;

	let user = $state(getCurrentUser());
	let userIsAdmin = $state(isAdmin());

	// Build base path for this org
	let basePath = $derived(orgSlug ? `/${orgSlug}` : '');

	function switchOrg() {
		// Clear current org and go to org selection
		localStorage.removeItem('currentOrgId');
		localStorage.removeItem('currentOrgName');
		localStorage.removeItem('currentOrgSlug');
		goto('/');
	}

	function toggleUserMenu() {
		showUserMenu = !showUserMenu;
	}

	function toggleMobileMenu() {
		showMobileMenu = !showMobileMenu;
	}

	function handleLogout() {
		logout();
	}

	// Close menus when clicking outside
	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.user-menu-container')) {
			showUserMenu = false;
		}
		if (!target.closest('.mobile-menu-container') && !target.closest('.mobile-menu-button')) {
			showMobileMenu = false;
		}
		if (!target.closest('.navbar-search-dropdown') && !target.closest('.navbar-search-btn')) {
			showSearch = false;
		}
	}

	function isActive(path: string, exact = false): boolean {
		const fullPath = `${basePath}${path}`;
		if (exact) {
			return $page.url.pathname === fullPath;
		}
		return $page.url.pathname.startsWith(fullPath);
	}

	// Check if any of the paths are active (for grouped nav items)
	function isAnyActive(paths: string[]): boolean {
		return paths.some(path => $page.url.pathname.startsWith(`${basePath}${path}`));
	}

	// Global search functionality
	function toggleSearch() {
		showSearch = !showSearch;
		if (showSearch) {
			searchQuery = '';
			searchResults = { customers: [], orders: [], products: [] };
		}
	}

	async function handleSearch() {
		if (!searchQuery.trim()) {
			searchResults = { customers: [], orders: [], products: [] };
			return;
		}

		// Debounce search
		if (searchTimeout) {
			clearTimeout(searchTimeout);
		}

		searchTimeout = setTimeout(async () => {
			searchLoading = true;
			try {
				const [customersRes, ordersRes, productsRes] = await Promise.all([
					adminListCustomers({ search: searchQuery, page_size: 5 }),
					listOrders({ limit: 5 }),
					listBillingProducts({})
				]);

				// Filter orders client-side since API might not support search
				const filteredOrders = (ordersRes.orders || []).filter(o =>
					o.email?.toLowerCase().includes(searchQuery.toLowerCase()) ||
					o.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
					o.order_number?.toLowerCase().includes(searchQuery.toLowerCase())
				).slice(0, 5);

				// Filter products client-side since API might not support search
				const filteredProducts = (productsRes.products || []).filter(p =>
					p.name.toLowerCase().includes(searchQuery.toLowerCase())
				).slice(0, 5);

				searchResults = {
					customers: customersRes.customers || [],
					orders: filteredOrders,
					products: filteredProducts
				};
			} catch (err) {
				console.error('Search failed:', err);
			} finally {
				searchLoading = false;
			}
		}, 300);
	}

	function navigateToResult(type: string, id: string) {
		showSearch = false;
		searchQuery = '';
		if (type === 'customer') {
			goto(`${basePath}/customers/${id}`);
		} else if (type === 'order') {
			goto(`${basePath}/commerce/orders/${id}`);
		} else if (type === 'product') {
			goto(`${basePath}/commerce/products/${id}`);
		}
	}

	const hasResults = $derived(
		searchResults.customers.length > 0 ||
		searchResults.orders.length > 0 ||
		searchResults.products.length > 0
	);
</script>

<svelte:window onclick={handleClickOutside} />

<nav class="navbar">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="flex h-16 justify-between">
			<!-- Mobile menu button -->
			<div class="flex items-center sm:hidden">
				<button
					onclick={(e) => {
						e.stopPropagation();
						toggleMobileMenu();
					}}
					class="mobile-menu-button navbar-mobile-btn"
					aria-expanded={showMobileMenu}
				>
					<span class="sr-only">Open main menu</span>
					{#if showMobileMenu}
						<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
						</svg>
					{:else}
						<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
						</svg>
					{/if}
				</button>
			</div>

			<!-- Logo and main navigation -->
			<div class="flex">
				<div class="flex flex-shrink-0 items-center">
					<a href="{basePath}/" class="navbar-brand">{org?.name || 'Outlet'}</a>
				</div>
				<div class="hidden sm:ml-6 sm:flex sm:space-x-8">
					<a href="{basePath}/customers" class="navbar-link" class:active={isActive('/customers')}>
						Customers
					</a>
					{#if userIsAdmin}
						<a href="{basePath}/marketing" class="navbar-link" class:active={isActive('/marketing')}>
							Marketing
						</a>
						<a href="{basePath}/commerce" class="navbar-link" class:active={isActive('/commerce')}>
							Commerce
						</a>
						<a href="{basePath}/ai" class="navbar-link" class:active={isActive('/ai')}>
							AI
						</a>
					{/if}
				</div>
			</div>

			<!-- Search and User menu (desktop only) -->
			<div class="hidden sm:flex items-center gap-2">
				<!-- Global Search -->
				<div class="relative">
					<button
						onclick={(e) => {
							e.stopPropagation();
							toggleSearch();
						}}
						class="navbar-search-btn"
						aria-label="Search"
					>
						<Search class="h-5 w-5" />
					</button>

					{#if showSearch}
						<div class="navbar-search-dropdown" onclick={(e) => e.stopPropagation()}>
							<div class="relative">
								<Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-text-muted" />
								<input
									type="text"
									placeholder="Search customers, orders, products..."
									bind:value={searchQuery}
									oninput={handleSearch}
									class="navbar-search-input"
									autofocus
								/>
								{#if searchQuery}
									<button
										onclick={() => { searchQuery = ''; searchResults = { customers: [], orders: [], products: [] }; }}
										class="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text"
									>
										<X class="h-4 w-4" />
									</button>
								{/if}
							</div>

							{#if searchLoading}
								<div class="p-4 text-center text-text-muted text-sm">
									Searching...
								</div>
							{:else if searchQuery && !hasResults}
								<div class="p-4 text-center text-text-muted text-sm">
									No results found for "{searchQuery}"
								</div>
							{:else if hasResults}
								<div class="navbar-search-results">
									{#if searchResults.customers.length > 0}
										<div class="navbar-search-section">
											<p class="navbar-search-section-title">Customers</p>
											{#each searchResults.customers as customer}
												<button
													onclick={() => navigateToResult('customer', customer.id)}
													class="navbar-search-result"
												>
													<span class="font-medium">{customer.name || customer.email}</span>
													<span class="text-text-muted">{customer.email}</span>
												</button>
											{/each}
										</div>
									{/if}

									{#if searchResults.orders.length > 0}
										<div class="navbar-search-section">
											<p class="navbar-search-section-title">Orders</p>
											{#each searchResults.orders as order}
												<button
													onclick={() => navigateToResult('order', String(order.id))}
													class="navbar-search-result"
												>
													<span class="font-medium">{order.order_number}</span>
													<span class="text-text-muted">{order.name || order.email}</span>
												</button>
											{/each}
										</div>
									{/if}

									{#if searchResults.products.length > 0}
										<div class="navbar-search-section">
											<p class="navbar-search-section-title">Products</p>
											{#each searchResults.products as product}
												<button
													onclick={() => navigateToResult('product', product.id)}
													class="navbar-search-result"
												>
													<span class="font-medium">{product.name}</span>
													<span class="text-text-muted">{product.type}</span>
												</button>
											{/each}
										</div>
									{/if}
								</div>
							{:else}
								<div class="p-4 text-center text-text-muted text-sm">
									Type to search across your site
								</div>
							{/if}
						</div>
					{/if}
				</div>

				<div class="user-menu-container relative ml-3">
					<div>
						<button
							onclick={(e) => {
								e.stopPropagation();
								toggleUserMenu();
							}}
							class="navbar-user-btn"
							id="user-menu-button"
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
						<div
							class="navbar-dropdown"
							role="menu"
							aria-orientation="vertical"
							aria-labelledby="user-menu-button"
						>
							<!-- User info header -->
							<div class="px-4 py-3 border-b border-border">
								<p class="text-sm font-medium text-text truncate">
									{user?.email || 'User'}
								</p>
								<p class="text-xs capitalize text-text-muted">{user?.role || 'agent'}</p>
							</div>
							<a
								href="{basePath}/settings"
								class="navbar-dropdown-item"
								role="menuitem"
								onclick={() => (showUserMenu = false)}
							>
								Settings
							</a>
							<div class="navbar-dropdown-divider"></div>
							{#if org}
								<button
									onclick={() => { showUserMenu = false; switchOrg(); }}
									class="navbar-dropdown-item w-full text-left"
									role="menuitem"
								>
									Switch Site
								</button>
								<div class="navbar-dropdown-divider"></div>
							{/if}
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

	<!-- Mobile menu -->
	{#if showMobileMenu}
		<div class="mobile-menu-container sm:hidden">
			<div class="space-y-1 pb-3 pt-2">
				<a
					href="{basePath}/customers"
					onclick={() => (showMobileMenu = false)}
					class="navbar-mobile-link"
					class:active={isActive('/customers')}
				>
					Customers
				</a>
				{#if userIsAdmin}
					<a
						href="{basePath}/marketing"
						onclick={() => (showMobileMenu = false)}
						class="navbar-mobile-link"
						class:active={isActive('/marketing')}
					>
						Marketing
					</a>
					<a
						href="{basePath}/commerce"
						onclick={() => (showMobileMenu = false)}
						class="navbar-mobile-link"
						class:active={isActive('/commerce')}
					>
						Commerce
					</a>
					<a
						href="{basePath}/ai"
						onclick={() => (showMobileMenu = false)}
						class="navbar-mobile-link"
						class:active={isActive('/ai')}
					>
						AI
					</a>
				{/if}
			</div>
			<!-- Mobile user section -->
			<div class="navbar-mobile-user">
				<div class="flex items-center px-4">
					<div class="navbar-avatar-lg">
						<span class="text-sm font-medium text-white">
							{(user?.name || user?.email || 'U')[0].toUpperCase()}
						</span>
					</div>
					<div class="ml-3">
						<div class="text-base font-medium text-text">{user?.name || user?.email || 'User'}</div>
						<div class="text-sm font-medium text-text-muted capitalize">{user?.role || 'agent'}</div>
					</div>
				</div>
				<div class="mt-3 space-y-1">
					<a
						href="{basePath}/settings"
						onclick={() => (showMobileMenu = false)}
						class="navbar-mobile-menu-item"
					>
						Settings
					</a>
					<div class="navbar-mobile-divider"></div>
					{#if org}
						<button
							onclick={() => { showMobileMenu = false; switchOrg(); }}
							class="navbar-mobile-menu-item w-full text-left"
						>
							Switch Site
						</button>
						<div class="navbar-mobile-divider"></div>
					{/if}
					<button
						onclick={() => { showMobileMenu = false; handleLogout(); }}
						class="navbar-mobile-menu-item w-full text-left"
					>
						Sign out
					</button>
				</div>
			</div>
		</div>
	{/if}
</nav>

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

		.navbar-mobile-btn {
			@apply inline-flex items-center justify-center rounded-md p-2 transition-colors;
			color: var(--color-text-muted);
		}

		.navbar-mobile-btn:hover {
			background-color: var(--color-bg-secondary);
			color: var(--color-text);
		}

		.navbar-mobile-btn:focus {
			@apply outline-none ring-2 ring-inset;
			--tw-ring-color: var(--color-primary);
		}

		.navbar-link {
			@apply inline-flex items-center border-b-2 px-1 pt-1 text-sm font-medium transition-colors;
			border-color: transparent;
			color: var(--color-text-muted);
		}

		.navbar-link:hover {
			border-color: var(--color-border);
			color: var(--color-text);
		}

		.navbar-link.active {
			border-color: var(--color-primary);
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

		.navbar-avatar-lg {
			@apply flex h-10 w-10 items-center justify-center rounded-full;
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

		.navbar-dropdown-divider {
			@apply my-1;
			border-top: 1px solid var(--color-border);
		}

		.navbar-mobile-link {
			@apply block border-l-4 py-2 pl-3 pr-4 text-base font-medium transition-colors;
			border-color: transparent;
			color: var(--color-text-muted);
		}

		.navbar-mobile-link:hover {
			border-color: var(--color-border);
			background-color: var(--color-bg-secondary);
			color: var(--color-text);
		}

		.navbar-mobile-link.active {
			border-color: var(--color-primary);
			background-color: color-mix(in srgb, var(--color-primary) 10%, transparent);
			color: var(--color-primary);
		}

		.navbar-mobile-user {
			@apply pb-3 pt-4;
			border-top: 1px solid var(--color-border);
		}

		.navbar-mobile-menu-item {
			@apply block px-4 py-2 text-base font-medium transition-colors;
			color: var(--color-text-muted);
		}

		.navbar-mobile-menu-item:hover {
			background-color: var(--color-bg-secondary);
			color: var(--color-text);
		}

		.navbar-mobile-divider {
			@apply my-1 mx-4;
			border-top: 1px solid var(--color-border);
		}

		.navbar-search-btn {
			@apply flex items-center justify-center h-9 w-9 rounded-lg transition-colors;
			color: var(--color-text-muted);
		}

		.navbar-search-btn:hover {
			background-color: var(--color-bg-secondary);
			color: var(--color-text);
		}

		.navbar-search-dropdown {
			@apply absolute right-0 z-50 mt-2 w-96 origin-top-right rounded-lg shadow-lg;
			background-color: var(--color-bg);
			border: 1px solid var(--color-border);
		}

		.navbar-search-input {
			@apply w-full pl-10 pr-10 py-3 text-sm rounded-t-lg border-0 border-b;
			background-color: var(--color-bg);
			color: var(--color-text);
			border-color: var(--color-border);
		}

		.navbar-search-input:focus {
			@apply outline-none ring-0;
		}

		.navbar-search-input::placeholder {
			color: var(--color-text-muted);
		}

		.navbar-search-results {
			@apply max-h-80 overflow-y-auto;
		}

		.navbar-search-section {
			@apply py-2;
		}

		.navbar-search-section-title {
			@apply px-4 py-1 text-xs font-medium uppercase tracking-wide;
			color: var(--color-text-muted);
		}

		.navbar-search-result {
			@apply w-full px-4 py-2 flex flex-col items-start text-left text-sm transition-colors;
			color: var(--color-text);
		}

		.navbar-search-result:hover {
			background-color: var(--color-bg-secondary);
		}
	}
</style>
