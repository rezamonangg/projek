import { getContext, setContext } from 'svelte';
import type { Component } from 'svelte';

type ToastType = 'success' | 'error' | 'warning' | 'info';

interface ToastApi {
	success: (message: string, duration?: number) => string;
	error: (message: string, duration?: number) => string;
	warning: (message: string, duration?: number) => string;
	info: (message: string, duration?: number) => string;
	addToast: (message: string, type: ToastType, duration?: number) => string;
	removeToast: (id: string) => void;
}

const TOAST_KEY = Symbol('toast');

export function setToastContext(api: ToastApi) {
	setContext(TOAST_KEY, api);
}

export function getToastContext(): ToastApi | undefined {
	return getContext(TOAST_KEY);
}

export type { ToastApi, ToastType };
