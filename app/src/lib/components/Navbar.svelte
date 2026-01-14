<script lang="ts">
	import { page } from '$app/stores';

	let { hasBuildInPublicPosts = false }: { hasBuildInPublicPosts?: boolean } = $props();
	let mobileMenuOpen = $state(false);

	const allNavLinks: Array<{ href: string; label: string; conditional?: boolean }> = [
		{ href: '/books', label: 'Books' },
		{ href: '/articles', label: 'Articles' },
		{ href: '/projects', label: 'Projects' },
		{ href: '/build-in-public', label: 'Build in Public', conditional: true },
		{ href: '/consulting', label: 'Consulting' },
		{ href: '/about', label: 'About' },
		{ href: '/contact', label: 'Contact' }
	];

	const navLinks = $derived(
		allNavLinks.filter((link) => !link.conditional || hasBuildInPublicPosts)
	);
</script>

<nav class="fixed top-0 z-50 w-full border-b border-gray-200 bg-white/80 backdrop-blur-md">
	<div class="mx-auto max-w-7xl px-6 md:px-8">
		<div class="flex h-20 items-center justify-between">
			<!-- Logo -->
			<a href="/" class="text-2xl font-semibold tracking-tight text-black">
				Your Platform
			</a>

			<!-- Desktop Navigation -->
			<div class="hidden items-center gap-8 md:flex">
				{#each navLinks as link}
					<a
						href={link.href}
						class="text-sm font-medium transition-colors"
						class:text-black={$page.url.pathname === link.href || $page.url.pathname.startsWith(link.href + '/')}
						class:text-gray-600={$page.url.pathname !== link.href && !$page.url.pathname.startsWith(link.href + '/')}
						class:hover:text-black={$page.url.pathname !== link.href}
					>
						{link.label}
					</a>
				{/each}
			</div>

			<!-- Mobile Menu Button -->
			<button
				class="flex size-10 items-center justify-center rounded-lg transition-colors hover:bg-gray-100 md:hidden"
				onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
				aria-label="Toggle menu"
			>
				{#if mobileMenuOpen}
					<svg class="size-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				{:else}
					<svg class="size-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
					</svg>
				{/if}
			</button>
		</div>

		<!-- Mobile Menu -->
		{#if mobileMenuOpen}
			<div class="border-t border-gray-200 py-4 md:hidden">
				{#each navLinks as link}
					<a
						href={link.href}
						class="block px-4 py-3 text-base font-medium transition-colors"
						class:text-black={$page.url.pathname === link.href || $page.url.pathname.startsWith(link.href + '/')}
						class:bg-gray-50={$page.url.pathname === link.href || $page.url.pathname.startsWith(link.href + '/')}
						class:text-gray-600={$page.url.pathname !== link.href && !$page.url.pathname.startsWith(link.href + '/')}
						onclick={() => (mobileMenuOpen = false)}
					>
						{link.label}
					</a>
				{/each}
			</div>
		{/if}
	</div>
</nav>

<!-- Spacer to prevent content from going under fixed navbar -->
<div class="h-20"></div>
