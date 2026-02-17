import { browser } from '$app/environment';
import { redirect } from '@sveltejs/kit';
import { authApi } from '$lib/api';
import { getStoredUser, setStoredUser, removeStoredUser } from '$lib/utils/storage';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ url }) => {
	if (!browser) {
		return { user: null };
	}

	const storedUser = getStoredUser();
	let user = null;

	if (storedUser) {
		try {
			user = JSON.parse(storedUser);
		} catch {
			user = null;
		}
	}

	if (!user) {
		const response = await authApi.getCurrentUser();
		if (response.data) {
			user = response.data;
			setStoredUser(JSON.stringify(user));
		}
	}

	if (!user) {
		removeStoredUser();
		throw redirect(302, '/login');
	}

	return { user };
};
