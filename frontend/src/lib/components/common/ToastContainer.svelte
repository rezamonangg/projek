<script lang="ts">
	import Toast from './Toast.svelte';
	import { setToastContext } from './toast-context';

	type ToastType = 'success' | 'error' | 'warning' | 'info';

	interface ToastItem {
		id: string;
		message: string;
		type: ToastType;
		duration?: number;
	}

	let toasts = $state<ToastItem[]>([]);

	function addToast(message: string, type: ToastType = 'info', duration: number = 5000) {
		const id = crypto.randomUUID();
		toasts = [...toasts, { id, message, type, duration }];
		return id;
	}

	function removeToast(id: string) {
		toasts = toasts.filter((t) => t.id !== id);
	}

	function success(message: string, duration?: number) {
		return addToast(message, 'success', duration);
	}

	function error(message: string, duration?: number) {
		return addToast(message, 'error', duration);
	}

	function warning(message: string, duration?: number) {
		return addToast(message, 'warning', duration);
	}

	function info(message: string, duration?: number) {
		return addToast(message, 'info', duration);
	}

	setToastContext({ success, error, warning, info, addToast, removeToast });
</script>

<div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
	{#each toasts as toast (toast.id)}
		<Toast {toast} onclose={removeToast} />
	{/each}
</div>
