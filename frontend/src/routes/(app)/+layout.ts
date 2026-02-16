import { browser } from '$app/environment';
import { redirect } from '@sveltejs/kit';
import { getStoredUser } from '$lib/utils/storage';
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
		throw redirect(302, '/login');
	}

	return { user };
};
