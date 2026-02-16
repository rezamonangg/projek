import { browser } from '$app/environment';

const TOKEN_KEY = 'projek_token';
const USER_KEY = 'projek_user';

export function getToken(): string | null {
	if (!browser) return null;
	return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
	if (!browser) return;
	localStorage.setItem(TOKEN_KEY, token);
}

export function removeToken(): void {
	if (!browser) return;
	localStorage.removeItem(TOKEN_KEY);
}

export function getStoredUser(): string | null {
	if (!browser) return null;
	return localStorage.getItem(USER_KEY);
}

export function setStoredUser(user: string): void {
	if (!browser) return;
	localStorage.setItem(USER_KEY, user);
}

export function removeStoredUser(): void {
	if (!browser) return;
	localStorage.removeItem(USER_KEY);
}
