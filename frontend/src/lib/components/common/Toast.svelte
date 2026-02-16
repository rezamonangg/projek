<script lang="ts">
	import { onMount } from 'svelte';

	type ToastType = 'success' | 'error' | 'warning' | 'info';

	interface Toast {
		id: string;
		message: string;
		type: ToastType;
		duration?: number;
	}

	interface Props {
		toast: Toast;
		onclose: (id: string) => void;
	}

	let { toast, onclose }: Props = $props();

	const iconMap: Record<ToastType, string> = {
		success: '✓',
		error: '✕',
		warning: '⚠',
		info: 'ℹ'
	};

	const colorMap: Record<ToastType, string> = {
		success: 'bg-green-50 border-green-200 text-green-800',
		error: 'bg-red-50 border-red-200 text-red-800',
		warning: 'bg-yellow-50 border-yellow-200 text-yellow-800',
		info: 'bg-blue-50 border-blue-200 text-blue-800'
	};

	onMount(() => {
		if (toast.duration !== 0) {
			const timer = setTimeout(() => {
				onclose(toast.id);
			}, toast.duration || 5000);
			return () => clearTimeout(timer);
		}
	});
</script>

<div
	class="flex items-center gap-3 px-4 py-3 rounded-lg border shadow-lg {colorMap[toast.type]}"
	role="alert"
>
	<span class="text-lg" aria-hidden="true">{iconMap[toast.type]}</span>
	<p class="flex-1 text-sm font-medium">{toast.message}</p>
	<button
		type="button"
		class="p-1 hover:bg-black/10 rounded transition-colors"
		onclick={() => onclose(toast.id)}
		aria-label="Dismiss notification"
	>
		<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
		</svg>
	</button>
</div>
