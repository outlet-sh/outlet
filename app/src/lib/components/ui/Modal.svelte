<!--
  Modal Component
  Consistent modal dialogs with header/body/footer structure
-->

<script lang="ts">
	interface Props {
		show: boolean;
		title: string;
		size?: 'sm' | 'md' | 'lg' | 'xl';
		closeOnBackdrop?: boolean;
		closeOnEscape?: boolean;
		showCloseButton?: boolean;
		onclose?: () => void;
		children: any;
		footer?: any;
	}

	let {
		show = $bindable(false),
		title,
		size = 'md',
		closeOnBackdrop = false,
		closeOnEscape = true,
		showCloseButton = true,
		onclose,
		children,
		footer
	}: Props = $props();

	const sizeClasses = {
		sm: 'max-w-md',
		md: 'max-w-lg',
		lg: 'max-w-2xl',
		xl: 'max-w-4xl'
	};

	const heightClasses = {
		sm: 'h-auto max-h-[60vh]',
		md: 'h-auto max-h-[70vh]',
		lg: 'h-auto max-h-[80vh]',
		xl: 'h-[85vh]'
	};

	function closeModal() {
		show = false;
		onclose?.();
	}

	function handleBackdropClick() {
		if (closeOnBackdrop) {
			closeModal();
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && closeOnEscape) {
			closeModal();
		}
	}
</script>

{#if show}
	<div
		class="modal-backdrop"
		onclick={handleBackdropClick}
		onkeydown={handleKeydown}
		role="button"
		tabindex="0"
		aria-label="Close modal"
	>
		<div
			class="modal-dialog modal-{size}"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.key === 'Enter' && e.stopPropagation()}
			role="dialog"
			aria-modal="true"
			tabindex="0"
		>
			<!-- Header -->
			<div class="modal-header">
				<div class="modal-header-content">
					<h3 class="modal-title">{title}</h3>
					{#if showCloseButton}
						<button
							onclick={closeModal}
							class="modal-close"
							aria-label="Close"
						>
							<svg class="modal-close-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M6 18L18 6M6 6l12 12"
								/>
							</svg>
						</button>
					{/if}
				</div>
			</div>

			<!-- Body (scrollable) -->
			<div class="modal-body">
				{@render children()}
			</div>

			<!-- Footer (if provided) -->
			{#if footer}
				<div class="modal-footer">
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	@reference "$src/app.css";
	@layer components.modal {
		.modal-backdrop {
			@apply fixed inset-0 z-50 flex items-center justify-center p-4;
			background: color-mix(in srgb, black 50%, transparent);
			backdrop-filter: blur(4px);
		}

		.modal-dialog {
			@apply w-full rounded-xl flex flex-col shadow-2xl bg-bg;
			backdrop-filter: blur(4px);
			border: 1px solid var(--color-border);
		}

		.modal-sm {
			@apply max-w-md h-auto;
			max-height: 60vh;
		}

		.modal-md {
			@apply max-w-lg h-auto;
			max-height: 70vh;
		}

		.modal-lg {
			@apply max-w-2xl h-auto;
			max-height: 80vh;
		}

		.modal-xl {
			@apply max-w-4xl;
			height: 85vh;
		}

		.modal-header {
			@apply flex-shrink-0 px-6 py-4;
			border-bottom: 1px solid var(--color-border);
		}

		.modal-header-content {
			@apply flex items-start justify-between;
		}

		.modal-title {
			@apply text-xl font-bold text-text;
		}

		.modal-close {
			@apply ml-auto rounded-lg p-1.5 transition-colors text-text-muted;

			&:hover {
				@apply bg-bg-secondary text-text;
			}
		}

		.modal-close-icon {
			@apply h-5 w-5;
		}

		.modal-body {
			@apply flex-1 overflow-y-auto p-6;
		}

		.modal-footer {
			@apply flex-shrink-0 px-6 py-4 bg-bg-secondary;
			border-top: 1px solid var(--color-border);
		}
	}
</style>
