<script lang="ts">
	import { Card, Button, Input, Alert, LoadingSpinner, SaveButton } from '$lib/components/ui';
	import { onMount } from 'svelte';
	import { getPlatformSettingsByCategory, updateEmailSettings } from '$lib/api/generate/outlet';

	let saving = $state(false);
	let saved = $state(false);
	let loading = $state(true);
	let error = $state('');

	let smtpHost = $state('');
	let smtpPort = $state(587);
	let smtpUser = $state('');
	let smtpPassword = $state('');
	let hasExistingPassword = $state(false);

	onMount(async () => {
		await loadSettings();
	});

	async function loadSettings() {
		try {
			const response = await getPlatformSettingsByCategory({}, 'email');
			for (const setting of response.settings) {
				switch (setting.key) {
					case 'smtp_host':
						smtpHost = setting.value || '';
						break;
					case 'smtp_port':
						smtpPort = setting.value ? parseInt(setting.value) : 587;
						break;
					case 'smtp_user':
						smtpUser = setting.value || '';
						break;
					case 'smtp_password':
						hasExistingPassword = !!setting.value;
						break;
				}
			}
		} catch (err: any) {
			error = err.message || 'Failed to load settings';
		} finally {
			loading = false;
		}
	}

	async function saveSettings() {
		saving = true;
		saved = false;
		error = '';

		try {
			await updateEmailSettings({
				smtp_host: smtpHost,
				smtp_port: smtpPort,
				smtp_user: smtpUser,
				smtp_password: smtpPassword || undefined
			});
			saved = true;
			setTimeout(() => {
				saved = false;
			}, 2000);
			if (smtpPassword) {
				hasExistingPassword = true;
				smtpPassword = '';
			}
		} catch (err: any) {
			error = err.message || 'Failed to save settings';
		} finally {
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>Email Settings - Outlet</title>
</svelte:head>

{#if error}
	<Alert type="error" title="Error">
		<p>{error}</p>
	</Alert>
{/if}

{#if loading}
	<div class="flex justify-center py-8">
		<LoadingSpinner />
	</div>
{:else}
	<Card title="Email Configuration" subtitle="Configure SMTP settings for sending emails">
		<div class="space-y-6">
			<div class="grid grid-cols-1 gap-6 sm:grid-cols-6">
				<div class="sm:col-span-4">
					<label for="smtp-host" class="form-label">SMTP Host</label>
					<Input type="text" id="smtp-host" bind:value={smtpHost} placeholder="smtp.example.com" />
				</div>

				<div class="sm:col-span-2">
					<label for="smtp-port" class="form-label">Port</label>
					<Input type="number" id="smtp-port" bind:value={smtpPort} />
				</div>

				<div class="sm:col-span-3">
					<label for="smtp-user" class="form-label">Username</label>
					<Input type="text" id="smtp-user" bind:value={smtpUser} />
				</div>

				<div class="sm:col-span-3">
					<label for="smtp-password" class="form-label">Password</label>
					<Input
						type="password"
						id="smtp-password"
						bind:value={smtpPassword}
						placeholder={hasExistingPassword ? '••••••••••••••••' : ''}
					/>
					{#if hasExistingPassword}
						<p class="mt-1 text-sm text-text-muted">Leave blank to keep existing password</p>
					{/if}
				</div>
			</div>

			<p class="text-sm text-text-muted">
				From Name and From Email are configured per organization.
			</p>

			<div class="flex justify-end">
				<SaveButton label="Save Settings" {saving} {saved} onclick={saveSettings} />
			</div>
		</div>
	</Card>
{/if}
