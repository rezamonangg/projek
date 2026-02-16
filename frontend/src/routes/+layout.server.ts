import type { LayoutServerLoad } from './$types';
import { authStore } from '$lib/stores/auth';
import { getStoredUser } from '$lib/api/mock/storage';

export const load: LayoutServerLoad = async () => {
	const storedUser = getStoredUser();
	let user = null;

	if (storedUser) {
		try {
			user = JSON.parse(storedUser);
		} catch {
			user = null;
		}
	}

	return {
		user
	};
};
