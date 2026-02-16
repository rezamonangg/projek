<script lang="ts">
	import { Button } from '$lib/components/common';
	import { onMount } from 'svelte';

	interface Props {
		children: import('svelte').Snippet;
	}

	let { children }: Props = $props();

	let error = $state<Error | null>(null);

	function reset() {
		error = null;
	}

	onMount(() => {
		const handleError = (event: ErrorEvent) => {
			error = event.error;
			event.preventDefault();
		};

		window.addEventListener('error', handleError);
		return () => window.removeEventListener('error', handleError);
	});
</script>

{#if error}
	<div class="min-h-[400px] flex items-center justify-center p-8">
		<div class="text-center max-w-md">
			<div class="w-16 h-16 mx-auto mb-4 bg-red-100 rounded-full flex items-center justify-center">
				<svg class="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
				</svg>
			</div>
			<h2 class="text-xl font-semibold text-gray-900 mb-2">Something went wrong</h2>
			<p class="text-gray-500 mb-4">{error.message || 'An unexpected error occurred'}</p>
			<Button onclick={reset}>Try Again</Button>
		</div>
	</div>
{:else}
	{@render children()}
{/if}
