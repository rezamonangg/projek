import type { ApiResponse, IAuthApi, RegisterInput } from '../types';
import type { User } from '$lib/types/api';
import { httpClient } from './client';
import { toCamelCase } from './transform';

interface LoginRequest {
	email: string;
	password: string;
}

interface LoginResponse {
	session_id: string;
	member: BackendMember;
}

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

async function login(email: string, password: string): Promise<ApiResponse<{ user: User; token: string }>> {
	try {
		const response = await httpClient.post<LoginResponse>('/auth/login', {
			email,
			password
		} as LoginRequest);

		const user = toFrontendUser(response.member);
		return createResponse({ user, token: response.session_id });
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Login failed';
		return createErrorResponse(message, 401);
	}
}

async function logout(): Promise<ApiResponse<void>> {
	try {
		await httpClient.post('/auth/logout');
		return createResponse(undefined);
	} catch {
		return createResponse(undefined);
	}
}

async function getCurrentUser(): Promise<ApiResponse<User>> {
	try {
		const member = await httpClient.get<BackendMember>('/auth/me');
		const user = toFrontendUser(member);
		return createResponse(user);
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Not authenticated';
		return createErrorResponse(message, 401);
	}
}

async function register(
	_input: RegisterInput
): Promise<ApiResponse<{ user: User; community: never; token: string }>> {
	return createErrorResponse('Registration is not implemented for real API', 501);
}

export const realAuthApi: IAuthApi = {
	login,
	logout,
	getCurrentUser,
	register
};
