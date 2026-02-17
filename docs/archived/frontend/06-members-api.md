# Frontend Implementation Plan: Members API

**Priority:** MEDIUM
**Estimated Time:** 30 minutes
**Backend Reference:** `docs/plan/backend/06-members-handler.md`

## Current State

### What Exists
- `frontend/src/lib/api/real/members.ts` - Implementation exists
- `frontend/src/lib/api/types.ts` - IMembersApi interface

### What to Verify
- Pagination response format
- Invite endpoint uses /members/invite
- Update uses /members/{id}

## Backend API Contract

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/members` | - | `{items: Member[], total, page, page_size, total_pages}` |
| POST | `/members/invite` | `{email, role}` | `{id: invitationId}` |
| PUT | `/members/{id}` | `{first_name, last_name}` | `Member` |

## Implementation Details

### Backend Response Format

```typescript
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
```

### Frontend Type Mapping

```typescript
// From $lib/types/api (User type)
interface User {
	id: string;
	communityId: string;
	email: string;
	firstName: string;
	lastName: string;
	role: 'admin' | 'member';
	isActive: boolean;
	createdAt: string;
	updatedAt: string;
}
```

### Transformer Function

```typescript
function toFrontendMember(member: BackendMember): User {
	return {
		id: member.id,
		communityId: member.community_id,
		email: member.email,
		firstName: member.first_name,
		lastName: member.last_name,
		role: member.role,
		isActive: member.is_active,
		createdAt: member.created_at,
		updatedAt: member.updated_at
	};
}
```

### API Implementation

```typescript
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

function toFrontendMember(member: BackendMember): User {
	return {
		id: member.id,
		communityId: member.community_id,
		email: member.email,
		firstName: member.first_name,
		lastName: member.last_name,
		role: member.role,
		isActive: member.is_active,
		createdAt: member.created_at,
		updatedAt: member.updated_at
	};
}

function createResponse<T>(data: T): ApiResponse<T> {
	return { data, error: null, status: 200 };
}

function createErrorResponse<T>(error: string, status: number): ApiResponse<T> {
	return { data: null, error, status };
}

async function listMembers(): Promise<ApiResponse<PaginatedResponse<User>>> {
	try {
		const response = await httpClient.get<BackendPaginatedMembers>('/members');
		return createResponse({
			items: response.items.map(toFrontendMember),
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
		const response = await httpClient.post<{ id: string }>('/members/invite', {
			email: input.email,
			role: input.role
		});
		return createResponse({ invitationId: response.id });
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to invite member';
		const status = message.includes('already exists') ? 409 : 500;
		return createErrorResponse(message, status);
	}
}

async function updateProfile(userId: string, input: UpdateProfileInput): Promise<ApiResponse<User>> {
	try {
		const member = await httpClient.put<BackendMember>(`/members/${userId}`, {
			first_name: input.firstName,
			last_name: input.lastName
		});
		return createResponse(toFrontendMember(member));
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
```

## Error Handling

| Status Code | Scenario | Action |
|-------------|----------|--------|
| 400 | Validation error | Show validation message |
| 401 | Unauthorized | Redirect to login |
| 403 | Forbidden (non-admin invite) | Show permission error |
| 404 | Member not found | Show not found message |
| 409 | Email already exists | Show conflict message |
| 500 | Server error | Show generic error |

## Files to Verify

| File | Status |
|------|--------|
| `frontend/src/lib/api/real/members.ts` | EXISTS - Verify alignment |
| `frontend/src/lib/api/types.ts` | EXISTS - Verify IMembersApi |
