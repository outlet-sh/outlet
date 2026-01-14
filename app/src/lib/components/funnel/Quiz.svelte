<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { trackEvent } from '$lib/analytics';
	import { trackFunnelStep } from '$lib/funnel';
	import { getQuiz, submitQuiz, type QuizInfo, type QuizQuestion } from '$lib/api';

	interface Props {
		quizSlug: string;
		collectEmail?: boolean;
		collectName?: boolean;
		onComplete?: (segments: string[], redirectUrl: string) => void;
		class?: string;
	}

	let {
		quizSlug,
		collectEmail = true,
		collectName = true,
		onComplete,
		class: className = ''
	}: Props = $props();

	let loading = $state(true);
	let quiz = $state<QuizInfo | null>(null);
	let loadError = $state('');

	let currentStep = $state(0); // 0 = questions, 1 = contact info (if collecting), 2 = results
	let currentQuestionIndex = $state(0);
	let answers = $state<Record<string, string>>({});
	let isSubmitting = $state(false);
	let submitError = $state('');

	let contactInfo = $state({
		name: '',
		email: ''
	});

	let results = $state<{ segments: string[]; redirectUrl: string } | null>(null);

	onMount(async () => {
		await loadQuiz();
		trackFunnelStep('quiz_started', `/quiz/${quizSlug}`);
	});

	async function loadQuiz() {
		loading = true;
		loadError = '';

		try {
			quiz = await getQuiz({}, quizSlug);
			if (!quiz) {
				loadError = 'Quiz not found';
			}
		} catch (err: any) {
			console.error('Failed to load quiz:', err);
			loadError = err.message || 'Failed to load quiz';
		} finally {
			loading = false;
		}
	}

	function selectAnswer(questionId: string, optionId: string) {
		answers = { ...answers, [questionId]: optionId };

		trackEvent('quiz_answer', {
			quiz_slug: quizSlug,
			question_id: questionId,
			option_id: optionId
		});

		// Auto-advance after short delay
		setTimeout(() => {
			if (quiz && currentQuestionIndex < quiz.questions.length - 1) {
				currentQuestionIndex++;
			} else {
				// All questions answered
				if (collectEmail || collectName) {
					currentStep = 1; // Go to contact info
				} else {
					handleSubmit();
				}
			}
		}, 300);
	}

	function goBack() {
		if (currentStep === 1) {
			currentStep = 0;
			currentQuestionIndex = (quiz?.questions.length ?? 1) - 1;
		} else if (currentQuestionIndex > 0) {
			currentQuestionIndex--;
		}
	}

	async function handleSubmit() {
		if (!quiz) return;

		isSubmitting = true;
		submitError = '';

		try {
			const response = await submitQuiz({}, {
				quiz_slug: quizSlug,
				answers,
				name: contactInfo.name || undefined,
				email: contactInfo.email || undefined
			}, quizSlug);

			if (!response.success) {
				throw new Error('Failed to submit quiz');
			}

			results = {
				segments: response.segments,
				redirectUrl: response.redirect_url || '/apply'
			};

			trackEvent('quiz_completed', {
				quiz_slug: quizSlug,
				segments: response.segments
			});

			trackFunnelStep('quiz_completed', `/quiz/${quizSlug}`);

			currentStep = 2; // Show results

			// Call onComplete callback if provided
			if (onComplete) {
				onComplete(response.segments, response.redirect_url || '/apply');
			}
		} catch (err: any) {
			submitError = err.message || 'Something went wrong. Please try again.';
		} finally {
			isSubmitting = false;
		}
	}

	function handleRedirect() {
		if (results?.redirectUrl && browser) {
			goto(results.redirectUrl);
		}
	}

	$effect(() => {
		if (quiz) {
			// Calculate progress
		}
	});

	// Progress calculation
	let progress = $derived(
		quiz ? ((currentQuestionIndex + 1) / quiz.questions.length) * 100 : 0
	);

	let currentQuestion = $derived(
		quiz?.questions[currentQuestionIndex] ?? null
	);
</script>

<div class="quiz {className}">
	{#if loading}
		<div class="flex justify-center items-center py-12">
			<div class="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-600"></div>
		</div>
	{:else if loadError}
		<div class="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
			<p class="text-red-800">{loadError}</p>
		</div>
	{:else if quiz}
		<div class="max-w-2xl mx-auto">
			<!-- Header -->
			<div class="text-center mb-8">
				<h2 class="text-2xl md:text-3xl font-bold text-slate-900">{quiz.title}</h2>
				{#if quiz.description && currentStep === 0}
					<p class="mt-2 text-slate-600">{quiz.description}</p>
				{/if}
			</div>

			<!-- Progress bar (only during questions) -->
			{#if currentStep === 0}
				<div class="mb-8">
					<div class="flex justify-between text-sm text-slate-500 mb-2">
						<span>Question {currentQuestionIndex + 1} of {quiz.questions.length}</span>
						<span>{Math.round(progress)}% complete</span>
					</div>
					<div class="h-2 bg-slate-200 rounded-full overflow-hidden">
						<div
							class="h-full bg-gradient-to-r from-blue-500 to-violet-500 transition-all duration-300"
							style="width: {progress}%"
						></div>
					</div>
				</div>
			{/if}

			<!-- Questions -->
			{#if currentStep === 0 && currentQuestion}
				<div class="bg-white rounded-2xl border border-slate-200 shadow-lg p-8">
					<h3 class="text-xl font-semibold text-slate-900 mb-6">{currentQuestion.question}</h3>

					<div class="space-y-3">
						{#each currentQuestion.options || [] as option}
							{@const isSelected = answers[currentQuestion.id] === option.id}
							<button
								type="button"
								onclick={() => selectAnswer(currentQuestion.id, option.id)}
								class="w-full text-left p-4 rounded-xl border-2 transition-all {isSelected
									? 'border-blue-500 bg-blue-50'
									: 'border-slate-200 hover:border-slate-300 hover:bg-slate-50'}"
							>
								<div class="flex items-center gap-3">
									<div class="w-5 h-5 rounded-full border-2 flex items-center justify-center flex-shrink-0 {isSelected ? 'border-blue-500' : 'border-slate-300'}">
										{#if isSelected}
											<div class="w-3 h-3 rounded-full bg-blue-500"></div>
										{/if}
									</div>
									<span class="text-slate-700">{option.label}</span>
								</div>
							</button>
						{/each}
					</div>

					{#if currentQuestionIndex > 0}
						<button
							type="button"
							onclick={goBack}
							class="mt-6 text-slate-500 hover:text-slate-700 text-sm flex items-center gap-1"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
							</svg>
							Previous question
						</button>
					{/if}
				</div>
			{/if}

			<!-- Contact Info Collection -->
			{#if currentStep === 1}
				<div class="bg-white rounded-2xl border border-slate-200 shadow-lg p-8">
					<h3 class="text-xl font-semibold text-slate-900 mb-2">Almost done!</h3>
					<p class="text-slate-600 mb-6">Enter your details to get your personalized results.</p>

					<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
						{#if collectName}
							<div>
								<label for="quiz-name" class="block text-sm font-semibold text-slate-700 mb-1">Name</label>
								<input
									type="text"
									id="quiz-name"
									bind:value={contactInfo.name}
									class="w-full px-4 py-3 rounded-lg border border-slate-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-colors"
									placeholder="Your name"
								/>
							</div>
						{/if}

						{#if collectEmail}
							<div>
								<label for="quiz-email" class="block text-sm font-semibold text-slate-700 mb-1">Email *</label>
								<input
									type="email"
									id="quiz-email"
									bind:value={contactInfo.email}
									required
									class="w-full px-4 py-3 rounded-lg border border-slate-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-colors"
									placeholder="you@example.com"
								/>
							</div>
						{/if}

						{#if submitError}
							<div class="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
								{submitError}
							</div>
						{/if}

						<div class="flex gap-3 pt-2">
							<button
								type="button"
								onclick={goBack}
								class="px-6 py-3 border border-slate-300 text-slate-700 font-semibold rounded-lg hover:bg-slate-50 transition-colors"
							>
								Back
							</button>
							<button
								type="submit"
								disabled={isSubmitting || (collectEmail && !contactInfo.email)}
								class="flex-1 py-3 px-6 bg-blue-600 text-white font-bold rounded-lg hover:bg-blue-700 disabled:bg-slate-400 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
							>
								{#if isSubmitting}
									<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
										<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
										<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
									</svg>
									Getting Results...
								{:else}
									Get My Results
								{/if}
							</button>
						</div>
					</form>
				</div>
			{/if}

			<!-- Results -->
			{#if currentStep === 2 && results}
				<div class="bg-white rounded-2xl border border-slate-200 shadow-lg p-8 text-center">
					<div class="w-16 h-16 mx-auto mb-4 rounded-full bg-green-100 flex items-center justify-center">
						<svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
						</svg>
					</div>

					<h3 class="text-2xl font-bold text-slate-900 mb-2">Your Results Are Ready!</h3>
					<p class="text-slate-600 mb-6">
						Based on your answers, we've identified the best next step for you.
					</p>

					{#if results.segments.length > 0}
						<div class="mb-6 p-4 bg-blue-50 rounded-lg">
							<p class="text-sm font-semibold text-blue-900 mb-2">Your Profile:</p>
							<div class="flex flex-wrap gap-2 justify-center">
								{#each results.segments as segment}
									<span class="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm font-medium">
										{segment.replace(/_/g, ' ')}
									</span>
								{/each}
							</div>
						</div>
					{/if}

					<button
						type="button"
						onclick={handleRedirect}
						class="w-full py-4 px-6 bg-gradient-to-r from-blue-500 to-violet-500 text-white font-bold rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all"
					>
						See Your Personalized Recommendations
					</button>
				</div>
			{/if}
		</div>
	{/if}
</div>
