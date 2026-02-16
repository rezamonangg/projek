import { writable } from 'svelte/store';

interface LoadingState {
	global: boolean;
	operations: Map<string, boolean>;
}

function createLoadingStore() {
	const { subscribe, update } = writable<LoadingState>({
		global: false,
		operations: new Map()
	});

	return {
		subscribe,
		start(operation?: string) {
			update((state) => {
				const newOperations = new Map(state.operations);
				if (operation) {
					newOperations.set(operation, true);
				}
				return {
					global: newOperations.size > 0,
					operations: newOperations
				};
			});
		},
		stop(operation?: string) {
			update((state) => {
				const newOperations = new Map(state.operations);
				if (operation) {
					newOperations.delete(operation);
				}
				return {
					global: newOperations.size > 0,
					operations: newOperations
				};
			});
		},
		isLoading(operation?: string): boolean {
			let loading = false;
			subscribe((state) => {
				if (operation) {
					loading = state.operations.get(operation) ?? false;
				} else {
					loading = state.global;
				}
			})();
			return loading;
		}
	};
}

export const loadingStore = createLoadingStore();

export async function withLoading<T>(operation: string, fn: () => Promise<T>): Promise<T> {
	loadingStore.start(operation);
	try {
		return await fn();
	} finally {
		loadingStore.stop(operation);
	}
}
