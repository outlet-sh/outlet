<script lang="ts">
	import { goto } from '$app/navigation';
	import * as api from '$lib/api';
	import { Card, Input, Button, Alert, Select } from '$lib/components/ui';
	import { ChevronRight, ChevronLeft, Mail, Server, Check, Loader2, HelpCircle, Lightbulb, Shield, Zap } from 'lucide-svelte';

	type Step = 'welcome' | 'smtp' | 'complete';
	let currentStep = $state<Step>('welcome');
	let loading = $state(false);
	let error = $state('');
	let testSuccess = $state(false);

	// SMTP form fields
	let smtpHost = $state('');
	let smtpPort = $state('587');
	let smtpUser = $state('');
	let smtpPassword = $state('');
	let fromEmail = $state('');
	let fromName = $state('');
	let replyTo = $state('');

	const portOptions = [
		{ value: '25', label: '25 (SMTP)' },
		{ value: '465', label: '465 (SMTPS)' },
		{ value: '587', label: '587 (Submission)' }
	];

	function nextStep() {
		if (currentStep === 'welcome') {
			currentStep = 'smtp';
		} else if (currentStep === 'smtp') {
			saveAndComplete();
		}
	}

	function prevStep() {
		if (currentStep === 'smtp') {
			currentStep = 'welcome';
		}
	}

	async function saveAndComplete() {
		loading = true;
		error = '';

		try {
			await api.updateEmailSettings({
				smtp_host: smtpHost,
				smtp_port: parseInt(smtpPort),
				smtp_user: smtpUser,
				smtp_password: smtpPassword,
				from_email: fromEmail,
				from_name: fromName,
				reply_to: replyTo || fromEmail
			});

			testSuccess = true;
			currentStep = 'complete';
		} catch (err: any) {
			console.error('Failed to save email settings:', err);
			error = err.message || 'Failed to save email settings';
		} finally {
			loading = false;
		}
	}

	function goToDashboard() {
		goto('/');
	}

	function skipForNow() {
		localStorage.setItem('setupSkipped', 'true');
		goto('/');
	}
</script>

<svelte:head>
	<title>Setup - Outlet</title>
</svelte:head>

<div class="grid grid-cols-1 lg:grid-cols-5 gap-8">
	<!-- Main Content -->
	<div class="lg:col-span-3">
		<Card>
			<!-- Progress indicator -->
			<div class="mb-8">
				<div class="flex items-center justify-between text-sm">
					<span class="text-text-muted">Step {currentStep === 'welcome' ? 1 : currentStep === 'smtp' ? 2 : 3} of 3</span>
				</div>
				<div class="mt-2 h-2 bg-border rounded-full overflow-hidden">
					<div
						class="h-full bg-primary transition-all duration-300"
						style="width: {currentStep === 'welcome' ? '33%' : currentStep === 'smtp' ? '66%' : '100%'}"
					></div>
				</div>
			</div>

			{#if currentStep === 'welcome'}
				<!-- Welcome Step -->
				<div class="text-center">
					<div class="mx-auto w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mb-6">
						<Mail class="w-8 h-8 text-primary" />
					</div>
					<h1 class="text-2xl font-bold text-text mb-2">Welcome to Outlet</h1>
					<p class="text-text-muted mb-8">
						Let's get you set up to send emails. This will only take a couple of minutes.
					</p>

					<div class="space-y-4 text-left mb-8">
						<div class="flex items-start gap-3 p-4 bg-bg-muted rounded-lg">
							<Server class="w-5 h-5 text-primary mt-0.5" />
							<div>
								<p class="font-medium text-text">SMTP Configuration</p>
								<p class="text-sm text-text-muted">Connect your email server to start sending</p>
							</div>
						</div>
					</div>

					<div class="flex gap-3">
						<Button type="primary" size="lg" onclick={nextStep} class="flex-1 justify-center">
							Get Started
							<ChevronRight class="w-4 h-4 ml-2" />
						</Button>
					</div>

					<button
						type="button"
						onclick={skipForNow}
						class="mt-4 text-sm text-text-muted hover:text-text transition-colors"
					>
						Skip for now
					</button>
				</div>
			{:else if currentStep === 'smtp'}
				<!-- SMTP Configuration Step -->
				<div>
					<div class="mb-6">
						<h2 class="text-xl font-bold text-text">Email Configuration</h2>
						<p class="text-text-muted mt-1">Enter your SMTP server details</p>
					</div>

					{#if error}
						<Alert type="error" title="Error" class="mb-6">
							<p>{error}</p>
						</Alert>
					{/if}

					<form class="space-y-5" onsubmit={(e) => { e.preventDefault(); nextStep(); }}>
						<div class="grid grid-cols-2 gap-4">
							<div class="col-span-2 sm:col-span-1">
								<label for="smtp-host" class="form-label">SMTP Host</label>
								<Input
									id="smtp-host"
									type="text"
									placeholder="smtp.example.com"
									bind:value={smtpHost}
									required
								/>
							</div>
							<div class="col-span-2 sm:col-span-1">
								<label for="smtp-port" class="form-label">SMTP Port</label>
								<Select
									id="smtp-port"
									bind:value={smtpPort}
									options={portOptions}
								/>
							</div>
						</div>

						<div>
							<label for="smtp-user" class="form-label">SMTP Username</label>
							<Input
								id="smtp-user"
								type="text"
								placeholder="your-username"
								bind:value={smtpUser}
								required
							/>
						</div>

						<div>
							<label for="smtp-password" class="form-label">SMTP Password</label>
							<Input
								id="smtp-password"
								type="password"
								placeholder="your-password"
								bind:value={smtpPassword}
								required
							/>
						</div>

						<hr class="border-border" />

						<div>
							<label for="from-email" class="form-label">From Email</label>
							<Input
								id="from-email"
								type="email"
								placeholder="noreply@example.com"
								bind:value={fromEmail}
								required
							/>
							<p class="mt-1 text-xs text-text-muted">Default sender email address</p>
						</div>

						<div>
							<label for="from-name" class="form-label">From Name</label>
							<Input
								id="from-name"
								type="text"
								placeholder="Your Company"
								bind:value={fromName}
							/>
						</div>

						<div>
							<label for="reply-to" class="form-label">Reply-To Email (optional)</label>
							<Input
								id="reply-to"
								type="email"
								placeholder="support@example.com"
								bind:value={replyTo}
							/>
						</div>

						<div class="flex gap-3 pt-4">
							<Button type="secondary" onclick={prevStep} disabled={loading}>
								<ChevronLeft class="w-4 h-4 mr-2" />
								Back
							</Button>
							<Button htmlType="submit" type="primary" disabled={loading || !smtpHost || !smtpUser || !smtpPassword || !fromEmail} class="flex-1 justify-center">
								{#if loading}
									<Loader2 class="w-4 h-4 mr-2 animate-spin" />
									Saving...
								{:else}
									Save & Complete
									<ChevronRight class="w-4 h-4 ml-2" />
								{/if}
							</Button>
						</div>
					</form>

					<button
						type="button"
						onclick={skipForNow}
						class="mt-4 text-sm text-text-muted hover:text-text transition-colors w-full text-center"
					>
						Skip for now
					</button>
				</div>
			{:else if currentStep === 'complete'}
				<!-- Complete Step -->
				<div class="text-center">
					<div class="mx-auto w-16 h-16 rounded-full bg-green-500/10 flex items-center justify-center mb-6">
						<Check class="w-8 h-8 text-green-500" />
					</div>
					<h2 class="text-2xl font-bold text-text mb-2">You're all set!</h2>
					<p class="text-text-muted mb-8">
						Your email configuration has been saved. You're ready to start sending emails.
					</p>

					<Button type="primary" size="lg" onclick={goToDashboard} class="w-full justify-center">
						Go to Dashboard
						<ChevronRight class="w-4 h-4 ml-2" />
					</Button>
				</div>
			{/if}
		</Card>
	</div>

	<!-- Help Sidebar -->
	<div class="lg:col-span-2">
		<div class="sticky top-8 space-y-6">
			{#if currentStep === 'welcome'}
				<!-- Welcome Help -->
				<Card>
					<div class="flex items-center gap-2 mb-4">
						<HelpCircle class="w-5 h-5 text-primary" />
						<h3 class="font-semibold text-text">What you'll need</h3>
					</div>
					<ul class="space-y-3 text-sm text-text-muted">
						<li class="flex items-start gap-2">
							<span class="text-primary mt-0.5">1.</span>
							<span>SMTP server credentials from your email provider</span>
						</li>
						<li class="flex items-start gap-2">
							<span class="text-primary mt-0.5">2.</span>
							<span>A verified email address to send from</span>
						</li>
						<li class="flex items-start gap-2">
							<span class="text-primary mt-0.5">3.</span>
							<span>About 2 minutes of your time</span>
						</li>
					</ul>
				</Card>

				<Card>
					<div class="flex items-center gap-2 mb-4">
						<Zap class="w-5 h-5 text-amber-500" />
						<h3 class="font-semibold text-text">Popular providers</h3>
					</div>
					<ul class="space-y-2 text-sm text-text-muted">
						<li>Amazon SES</li>
						<li>SendGrid</li>
						<li>Mailgun</li>
						<li>Postmark</li>
						<li>Any SMTP server</li>
					</ul>
				</Card>
			{:else if currentStep === 'smtp'}
				<!-- SMTP Help -->
				<Card>
					<div class="flex items-center gap-2 mb-4">
						<HelpCircle class="w-5 h-5 text-primary" />
						<h3 class="font-semibold text-text">SMTP Settings Help</h3>
					</div>
					<div class="space-y-4 text-sm">
						<div>
							<p class="font-medium text-text">SMTP Host</p>
							<p class="text-text-muted">The hostname of your SMTP server (e.g., smtp.gmail.com, email-smtp.us-east-1.amazonaws.com)</p>
						</div>
						<div>
							<p class="font-medium text-text">SMTP Port</p>
							<p class="text-text-muted"><strong>587</strong> is recommended for most providers. Use <strong>465</strong> for SSL/TLS or <strong>25</strong> for unencrypted.</p>
						</div>
						<div>
							<p class="font-medium text-text">Username & Password</p>
							<p class="text-text-muted">Your SMTP credentials. For AWS SES, these are your IAM SMTP credentials (not your AWS access keys).</p>
						</div>
					</div>
				</Card>

				<Card>
					<div class="flex items-center gap-2 mb-4">
						<Lightbulb class="w-5 h-5 text-amber-500" />
						<h3 class="font-semibold text-text">Quick setup guides</h3>
					</div>
					<div class="space-y-3 text-sm">
						<div class="p-3 bg-bg-muted rounded-lg">
							<p class="font-medium text-text">Amazon SES</p>
							<p class="text-text-muted text-xs mt-1">Host: email-smtp.[region].amazonaws.com<br/>Port: 587</p>
						</div>
						<div class="p-3 bg-bg-muted rounded-lg">
							<p class="font-medium text-text">SendGrid</p>
							<p class="text-text-muted text-xs mt-1">Host: smtp.sendgrid.net<br/>Port: 587<br/>User: apikey</p>
						</div>
						<div class="p-3 bg-bg-muted rounded-lg">
							<p class="font-medium text-text">Mailgun</p>
							<p class="text-text-muted text-xs mt-1">Host: smtp.mailgun.org<br/>Port: 587</p>
						</div>
					</div>
				</Card>

				<Card>
					<div class="flex items-center gap-2 mb-4">
						<Shield class="w-5 h-5 text-green-500" />
						<h3 class="font-semibold text-text">Security</h3>
					</div>
					<p class="text-sm text-text-muted">
						Your SMTP password is encrypted before being stored. We never log or expose your credentials.
					</p>
				</Card>
			{:else if currentStep === 'complete'}
				<!-- Complete Help -->
				<Card>
					<div class="flex items-center gap-2 mb-4">
						<Lightbulb class="w-5 h-5 text-amber-500" />
						<h3 class="font-semibold text-text">What's next?</h3>
					</div>
					<ul class="space-y-3 text-sm text-text-muted">
						<li class="flex items-start gap-2">
							<Check class="w-4 h-4 text-green-500 mt-0.5" />
							<span>Create your first email list</span>
						</li>
						<li class="flex items-start gap-2">
							<Check class="w-4 h-4 text-green-500 mt-0.5" />
							<span>Import or add subscribers</span>
						</li>
						<li class="flex items-start gap-2">
							<Check class="w-4 h-4 text-green-500 mt-0.5" />
							<span>Design your first campaign</span>
						</li>
						<li class="flex items-start gap-2">
							<Check class="w-4 h-4 text-green-500 mt-0.5" />
							<span>Set up automation sequences</span>
						</li>
					</ul>
				</Card>
			{/if}
		</div>
	</div>
</div>
