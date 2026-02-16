import type { ApiResponse, PaginatedResponse } from '../types';
import type { User } from '$lib/types/api';
import { mockResponse, mockError, generateId } from './utils';
import { generateMockUser } from './generators';

const mockMembers: User[] = [
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
		email: 'john@example.com',
		firstName: 'John',
		lastName: 'Doe',
		role: 'member',
		communityId: 'community-1'
	}),
	generateMockUser({
		id: 'user-3',
		email: 'jane@example.com',
		firstName: 'Jane',
		lastName: 'Smith',
		role: 'member',
		communityId: 'community-1'
	}),
	generateMockUser({
		id: 'user-4',
		email: 'bob@example.com',
		firstName: 'Bob',
		lastName: 'Wilson',
		role: 'member',
		communityId: 'community-1'
	})
];

export async function listMembers(): Promise<ApiResponse<PaginatedResponse<User>>> {
	return mockResponse({
		items: mockMembers,
		total: mockMembers.length,
		page: 1,
		pageSize: 10,
		totalPages: 1
	});
}

interface InviteMemberInput {
	email: string;
	role: 'admin' | 'member';
}

export async function inviteMember(
	input: InviteMemberInput
): Promise<ApiResponse<{ invitationId: string }>> {
	const existingMember = mockMembers.find((m) => m.email === input.email);
	if (existingMember) {
		return mockError('A member with this email already exists', 409);
	}

	const invitationId = generateId();
	return mockResponse({ invitationId });
}

interface UpdateProfileInput {
	firstName?: string;
	lastName?: string;
}

export async function updateProfile(
	userId: string,
	input: UpdateProfileInput
): Promise<ApiResponse<User>> {
	const memberIndex = mockMembers.findIndex((m) => m.id === userId);
	if (memberIndex === -1) {
		return mockError('User not found', 404);
	}

	mockMembers[memberIndex] = {
		...mockMembers[memberIndex],
		...input
	};

	return mockResponse(mockMembers[memberIndex]);
}
