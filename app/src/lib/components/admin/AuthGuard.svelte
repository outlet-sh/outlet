<script lang="ts">
	import { goto } from '$app/navigation';
	import { isAuthenticated, getCurrentUser, isAdmin } from '$lib/auth';

	let { requireAdmin = false, redirectTo = '/auth/login' }: { requireAdmin?: boolean; redirectTo?: string } = $props();

	let isLoading = $state(true);
	let isAuthorized = $state(false);
	let user = $state<any>(null);

	$effect(() => {
		checkAuth();
	});

	function checkAuth() {
		// Check if user is authenticated
		if (!isAuthenticated()) {
			goto(redirectTo);
			return;
		}

		// Get current user
		user = getCurrentUser();

		// Check admin requirement
		if (requireAdmin && !isAdmin()) {
			goto('/'); // Redirect non-admins to admin
			return;
		}

		isAuthorized = true;
		isLoading = false;
	}
</script>

{#if isLoading}
	<div class="flex items-center justify-center min-h-screen bg-gray-50">
		<div class="text-center">
			<div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-600">Loading...</p>
		</div>
	</div>
{:else if isAuthorized}
	<slot {user} />
{/if}
