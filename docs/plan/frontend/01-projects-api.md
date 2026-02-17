# Frontend Implementation Plan: Projects API

**Priority:** HIGH (blocking user flow)
**Estimated Time:** 30 minutes
**Backend Reference:** `docs/plan/backend/01-projects-handler.md`

## Current State

### What Exists
- `frontend/src/lib/api/real/projects.ts` - Full implementation
- `frontend/src/lib/api/types.ts` - IProjectsApi interface

### What to Verify/Update
- Ensure transformer matches backend response structure
- Verify pagination response format
- Check error handling for all status codes

## Backend API Contract

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/projects` | - | `{items: Project[], total, page, page_size, total_pages}` |
| POST | `/projects` | `{name, description, key}` | `Project` |
| GET | `/projects/{id}` | - | `Project` |
| PUT | `/projects/{id}` | `{name?, description?, is_archived?}` | `Project` |
| DELETE | `/projects/{id}` | - | `void` |

## Implementation Details

### Backend Response Format

```typescript
interface BackendProject {
	id: string;
	community_id: string;
	name: string;
	description: string;
	key: string;
	is_archived: boolean;
	created_at: string;
	updated_at: string;
}

interface BackendPaginatedProjects {
	items: BackendProject[];
	total: number;
	page: number;
	page_size: number;
	total_pages: number;
}
```

### Frontend Type Mapping

```typescript
// From $lib/types/api
interface Project {
	id: string;
	name: string;
	description: string;
	status: 'active' | 'archived';
	communityId: string;
	createdAt: string;
	updatedAt: string;
}
```

### Transformer Function

```typescript
function toFrontendProject(project: BackendProject): Project {
	return {
		id: project.id,
		name: project.name,
		description: project.description,
		status: project.is_archived ? 'archived' : 'active',
		communityId: project.community_id,
		createdAt: project.created_at,
		updatedAt: project.updated_at
	};
}
```

### API Implementation

```typescript
import type { ApiResponse, PaginatedResponse, IProjectsApi, CreateProjectInput, UpdateProjectInput } from '../types';
import type { Project } from '$lib/types/api';
import { httpClient } from './client';

// Backend types (snake_case)
interface BackendProject {
	id: string;
	community_id: string;
	name: string;
	description: string;
	key: string;
	is_archived: boolean;
	created_at: string;
	updated_at: string;
}

interface BackendPaginatedProjects {
	items: BackendProject[];
	total: number;
	page: number;
	page_size: number;
	total_pages: number;
}

// Transformer
function toFrontendProject(project: BackendProject): Project {
	return {
		id: project.id,
		name: project.name,
		description: project.description,
		status: project.is_archived ? 'archived' : 'active',
		communityId: project.community_id,
		createdAt: project.created_at,
		updatedAt: project.updated_at
	};
}

// Helper functions
function createResponse<T>(data: T): ApiResponse<T> {
	return { data, error: null, status: 200 };
}

function createErrorResponse<T>(error: string, status: number): ApiResponse<T> {
	return { data: null, error, status };
}

// API methods
async function listProjects(): Promise<ApiResponse<PaginatedResponse<Project>>> {
	try {
		const response = await httpClient.get<BackendPaginatedProjects>('/projects');
		return createResponse({
			items: response.items.map(toFrontendProject),
			total: response.total,
			page: response.page,
			pageSize: response.page_size,
			totalPages: response.total_pages
		});
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load projects';
		return createErrorResponse(message, 500);
	}
}

async function getProject(id: string): Promise<ApiResponse<Project>> {
	try {
		const project = await httpClient.get<BackendProject>(`/projects/${id}`);
		return createResponse(toFrontendProject(project));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load project';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function createProject(input: CreateProjectInput): Promise<ApiResponse<Project>> {
	try {
		const key = input.name.substring(0, 3).toUpperCase();
		const project = await httpClient.post<BackendProject>('/projects', {
			name: input.name,
			description: input.description,
			key
		});
		return createResponse(toFrontendProject(project));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to create project';
		return createErrorResponse(message, 500);
	}
}

async function updateProject(id: string, input: UpdateProjectInput): Promise<ApiResponse<Project>> {
	try {
		const body: Record<string, unknown> = {};
		if (input.name) body.name = input.name;
		if (input.description) body.description = input.description;
		if (input.status) body.is_archived = input.status === 'archived';

		const project = await httpClient.put<BackendProject>(`/projects/${id}`, body);
		return createResponse(toFrontendProject(project));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to update project';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function deleteProject(id: string): Promise<ApiResponse<void>> {
	try {
		await httpClient.delete(`/projects/${id}`);
		return createResponse(undefined);
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to delete project';
		return createErrorResponse(message, 500);
	}
}

export const realProjectsApi: IProjectsApi = {
	listProjects,
	getProject,
	createProject,
	updateProject,
	deleteProject
};
```

## Error Handling

| Status Code | Scenario | Action |
|-------------|----------|--------|
| 400 | Validation error | Show validation message |
| 401 | Unauthorized | Redirect to login |
| 403 | Forbidden | Show permission error |
| 404 | Not found | Show not found message |
| 409 | Conflict (duplicate key) | Show conflict message |
| 500 | Server error | Show generic error |

## Test Cases

| Test | Description |
|------|-------------|
| `listProjects_returnsPaginatedResponse` | List returns correct pagination |
| `getProject_byId_returnsProject` | Get by ID works |
| `getProject_notFound_returns404` | 404 for missing project |
| `createProject_validInput_returnsProject` | Create works |
| `createProject_validationError_returns400` | Validation handled |
| `updateProject_partialUpdate_returnsProject` | Partial update works |
| `deleteProject_existingId_returnsVoid` | Delete works |

## Verification

```bash
# Type check
cd frontend && npm run check

# Lint
cd frontend && npm run lint

# Test
cd frontend && npm test -- projects
```

## Files to Verify

| File | Status |
|------|--------|
| `frontend/src/lib/api/real/projects.ts` | EXISTS - Verify alignment |
| `frontend/src/lib/api/types.ts` | EXISTS - Verify IProjectsApi |
