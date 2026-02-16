import type { ApiResponse, PaginatedResponse, IMembersApi, InviteMemberInput, UpdateProfileInput } from '../types';
import type { User } from '$lib/types/api';
import { httpClient } from './client';

interface BackendMember {
	id: string;
	community_id: string;
	email: string;
	first_name: string;
	last_name: string;
	role: 'admin' | 'member';
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

interface BackendPaginatedMembers {
	items: BackendMember[];
	total: number;
	page: number;
	page_size: number;
	total_pages: number;
}

interface InvitationResponse {
	id: string;
}

function toFrontendUser(member: BackendMember): User {
	return {
		id: member.id,
		email: member.email,
		firstName: member.first_name,
		lastName: member.last_name,
		role: member.role,
		communityId: member.community_id,
		createdAt: member.created_at
	};
}

function createResponse<T>(data: T): ApiResponse<T> {
	return {
		data,
		error: null,
		status: 200
	};
}

function createErrorResponse<T>(error: string, status: number): ApiResponse<T> {
	return {
		data: null,
		error,
		status
	};
}

async function listMembers(): Promise<ApiResponse<PaginatedResponse<User>>> {
	try {
		const response = await httpClient.get<BackendPaginatedMembers>('/members');
		const items = response.items.map(toFrontendUser);
		return createResponse({
			items,
			total: response.total,
			page: response.page,
			pageSize: response.page_size,
			totalPages: response.total_pages
		});
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load members';
		return createErrorResponse(message, 500);
	}
}

async function inviteMember(input: InviteMemberInput): Promise<ApiResponse<{ invitationId: string }>> {
	try {
		const response = await httpClient.post<InvitationResponse>('/members/invite', {
			email: input.email,
			role: input.role
		});
		return createResponse({ invitationId: response.id });
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to invite member';
		const status = message.includes('already') ? 409 : 500;
		return createErrorResponse(message, status);
	}
}

async function updateProfile(userId: string, input: UpdateProfileInput): Promise<ApiResponse<User>> {
	try {
		const response = await httpClient.put<BackendMember>(`/members/${userId}`, {
			first_name: input.firstName,
			last_name: input.lastName
		});
		const user = toFrontendUser(response);
		return createResponse(user);
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to update profile';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

export const realMembersApi: IMembersApi = {
	listMembers,
	inviteMember,
	updateProfile
};
