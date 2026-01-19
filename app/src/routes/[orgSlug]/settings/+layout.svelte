<script lang="ts">
	import { page } from '$app/stores';

	const { children } = $props();

	// Get orgSlug from URL for building paths
	let orgSlug = $derived($page.params.orgSlug);

	let tabs = $derived([
		{ name: 'General', href: `/${orgSlug}/settings` },
		{ name: 'Email', href: `/${orgSlug}/settings/email` },
		{ name: 'Webhooks', href: `/${orgSlug}/settings/webhooks` },
		{ name: 'Privacy', href: `/${orgSlug}/settings/privacy` }
	]);

	function isActive(href: string): boolean {
		return $page.url.pathname === href;
	}
</script>

<div class="py-6">
	<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="mb-6">
			<h1 class="text-2xl font-semibold text-text">Workspace Settings</h1>
			<p class="mt-1 text-sm text-text-muted">Manage your workspace configuration</p>
		</div>

		<div class="border-b border-border">
			<nav class="-mb-px flex space-x-8">
				{#each tabs as tab}
					<a
						href={tab.href}
						class="whitespace-nowrap border-b-2 py-4 px-1 text-sm font-medium {isActive(tab.href)
							? 'border-primary text-primary'
							: 'border-transparent text-text-muted hover:border-border hover:text-text'}"
					>
						{tab.name}
					</a>
				{/each}
			</nav>
		</div>

		<div class="mt-6">
			{@render children()}
		</div>
	</div>
</div>
