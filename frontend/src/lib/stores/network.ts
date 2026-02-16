import { browser } from '$app/environment';
import { writable } from 'svelte/store';

const offlineStore = writable(false);

if (browser) {
	offlineStore.set(!navigator.onLine);

	window.addEventListener('online', () => offlineStore.set(false));
	window.addEventListener('offline', () => offlineStore.set(true));
}

export const offline = offlineStore;
