import type { ApiResponse, PaginatedResponse } from '../types';
import type { User } from '$lib/types/api';
import { mockResponse } from './utils';
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
