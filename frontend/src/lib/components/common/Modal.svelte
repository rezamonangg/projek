<script lang="ts">
	import { onMount } from 'svelte';

	interface Props {
		open: boolean;
		title?: string;
		size?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
		onclose: () => void;
		children?: import('svelte').Snippet;
		footer?: import('svelte').Snippet;
	}

	let { open, title, size = 'md', onclose, children, footer }: Props = $props();

	const sizeClasses: Record<string, string> = {
		sm: 'max-w-sm',
		md: 'max-w-md',
		lg: 'max-w-lg',
		xl: 'max-w-xl',
		full: 'max-w-4xl'
	};

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && open) {
			onclose();
		}
	}

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			onclose();
		}
	}

	onMount(() => {
		document.addEventListener('keydown', handleKeydown);
		return () => document.removeEventListener('keydown', handleKeydown);
	});

	$effect(() => {
		if (open) {
			document.body.style.overflow = 'hidden';
		} else {
			document.body.style.overflow = '';
		}
		return () => {
			document.body.style.overflow = '';
		};
	});
</script>

{#if open}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
		onclick={handleBackdropClick}
		role="dialog"
		tabindex="-1"
		aria-modal="true"
		aria-labelledby={title ? 'modal-title' : undefined}
	>
		<div
			class="w-full {sizeClasses[size]} bg-white rounded-xl shadow-2xl transform transition-all"
		>
			{#if title}
				<div class="flex items-center justify-between px-6 py-4 border-b border-gray-200">
					<h2 id="modal-title" class="text-lg font-semibold text-gray-900">
						{title}
					</h2>
					<button
						type="button"
						class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
						onclick={onclose}
						aria-label="Close modal"
					>
						<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M6 18L18 6M6 6l12 12"
							/>
						</svg>
					</button>
				</div>
			{/if}

			<div class="px-6 py-4">
				{#if children}
					{@render children()}
				{/if}
			</div>

			{#if footer}
				<div class="px-6 py-4 border-t border-gray-200 bg-gray-50 rounded-b-xl">
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
{/if}
