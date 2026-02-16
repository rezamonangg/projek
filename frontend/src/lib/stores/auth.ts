import { writable } from 'svelte/store';
import { goto } from '$app/navigation';
import { authApi } from '$lib/api';
import { setToken, removeToken, setStoredUser, removeStoredUser } from '$lib/api/mock/storage';
import type { User } from '$lib/types/api';

interface AuthState {
	user: User | null;
	loading: boolean;
	error: string | null;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		loading: false,
		error: null
	});

	return {
		subscribe,
		async login(email: string, password: string) {
			update((s) => ({ ...s, loading: true, error: null }));

			const response = await authApi.login(email, password);

			if (response.error) {
				update((s) => ({ ...s, loading: false, error: response.error }));
				return false;
			}

			if (response.data) {
				setToken(response.data.token);
				setStoredUser(JSON.stringify(response.data.user));
				set({ user: response.data.user, loading: false, error: null });
				goto('/dashboard');
				return true;
			}

			return false;
		},
		async logout() {
			await authApi.logout();
			removeToken();
			removeStoredUser();
			set({ user: null, loading: false, error: null });
			goto('/login');
		},
		async register(input: {
			communityName: string;
			slug: string;
			adminEmail: string;
			adminPassword: string;
		}) {
			update((s) => ({ ...s, loading: true, error: null }));

			const response = await authApi.register(input);

			if (response.error) {
				update((s) => ({ ...s, loading: false, error: response.error }));
				return false;
			}

			if (response.data) {
				setToken(response.data.token);
				setStoredUser(JSON.stringify(response.data.user));
				set({ user: response.data.user, loading: false, error: null });
				goto('/dashboard');
				return true;
			}

			return false;
		},
		async loadUser() {
			update((s) => ({ ...s, loading: true }));

			const response = await authApi.getCurrentUser();

			if (response.data) {
				set({ user: response.data, loading: false, error: null });
			} else {
				set({ user: null, loading: false, error: null });
			}
		},
		setUser(user: User | null) {
			update((s) => ({ ...s, user }));
		}
	};
}

export const authStore = createAuthStore();
