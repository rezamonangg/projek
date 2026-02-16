import type { ApiResponse } from '../types';
import type { User, Community } from '$lib/types/api';
import { mockResponse, mockError, generateId } from './utils';
import { generateMockUser, generateMockCommunity } from './generators';

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

interface RegisterInput {
	communityName: string;
	slug: string;
	adminEmail: string;
	adminPassword: string;
}

export async function register(
	input: RegisterInput
): Promise<ApiResponse<{ user: User; community: Community; token: string }>> {
	const existingUser = mockUsers.find((u) => u.email === input.adminEmail);
	if (existingUser) {
		return mockError('An account with this email already exists', 409);
	}

	const community = generateMockCommunity({
		id: generateId(),
		name: input.communityName,
		slug: input.slug
	});

	const user = generateMockUser({
		id: generateId(),
		email: input.adminEmail,
		firstName: 'Admin',
		lastName: 'User',
		role: 'admin',
		communityId: community.id
	});

	mockUsers.push(user);

	const token = generateId();
	return mockResponse({ user, community, token });
}
