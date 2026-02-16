import type { ApiResponse } from '../types';
import type { User } from '$lib/types/api';
import { mockResponse, mockError, generateId } from './utils';
import { generateMockUser } from './generators';

const mockUsers: User[] = [
	generateMockUser({
		id: 'user-1',
		email: 'admin@example.com',
		firstName: 'Admin',
		lastName: 'User',
		role: 'admin',
		communityId: 'community-1'
	}),
	generateMockUser({
		id: 'user-2',
		email: 'member@example.com',
		firstName: 'John',
		lastName: 'Member',
		role: 'member',
		communityId: 'community-1'
	})
];

export async function login(email: string, password: string): Promise<ApiResponse<{ user: User; token: string }>> {
	const user = mockUsers.find((u) => u.email === email);

	if (!user || password !== 'password123') {
		return mockError('Invalid email or password', 401);
	}

	const token = generateId();
	return mockResponse({ user, token });
}

export async function logout(): Promise<ApiResponse<void>> {
	return mockResponse(undefined);
}

export async function getCurrentUser(): Promise<ApiResponse<User>> {
	return mockResponse(mockUsers[0]);
}
