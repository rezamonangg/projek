<script lang="ts">
	import { browser } from '$app/environment';
	import { offline } from '$lib/stores/network';

	let showRestored = $state(false);
	let wasOffline = $state(false);

	$effect(() => {
		if ($offline) {
			wasOffline = true;
			showRestored = false;
		} else if (wasOffline && !$offline) {
			showRestored = true;
			setTimeout(() => {
				showRestored = false;
			}, 3000);
		}
	});
</script>

{#if browser}
	{#if $offline}
		<div class="fixed bottom-4 left-1/2 -translate-x-1/2 z-50 bg-red-600 text-white px-4 py-2 rounded-lg shadow-lg flex items-center gap-2">
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636a9 9 0 010 12.728m0 0l-2.829-2.829m2.829 2.829L21 21M15.536 8.464a5 5 0 010 7.072m0 0l-2.829-2.829m-4.243 2.829a4.978 4.978 0 01-1.414-2.83m-1.414 5.658a9 9 0 01-2.167-9.238m7.824 2.167a1 1 0 111.414 1.414m-1.414-1.414L3 3m8.293 8.293l1.414 1.414" />
			</svg>
			<span class="text-sm font-medium">You are offline</span>
		</div>
	{:else if showRestored}
		<div class="fixed bottom-4 left-1/2 -translate-x-1/2 z-50 bg-green-600 text-white px-4 py-2 rounded-lg shadow-lg flex items-center gap-2">
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
			</svg>
			<span class="text-sm font-medium">Back online</span>
		</div>
	{/if}
{/if}
