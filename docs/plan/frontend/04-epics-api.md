# Frontend Implementation Plan: Epics API

**Priority:** MEDIUM
**Estimated Time:** 20 minutes
**Backend Reference:** `docs/plan/backend/02-epics-handler.md`

## Current State

### What Exists
- `frontend/src/lib/api/real/epics.ts` - Implementation exists
- `frontend/src/lib/api/types.ts` - IEpicsApi interface

### What to Verify
- Endpoints use nested project routes
- Update/delete use direct /epics/{id}

## Backend API Contract

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/projects/{projectId}/epics` | - | `Epic[]` |
| POST | `/projects/{projectId}/epics` | `{name, description}` | `Epic` |
| PUT | `/epics/{id}` | `{name, description}` | `Epic` |
| DELETE | `/epics/{id}` | - | `void` |

## Implementation Details

### Backend Response Format

```typescript
interface BackendEpic {
	id: string;
	project_id: string;
	name: string;
	description: string;
	start_date?: string;
	end_date?: string;
	created_at: string;
	updated_at: string;
}
```

### Frontend Type Mapping

```typescript
// From $lib/types/api
interface Epic {
	id: string;
	projectId: string;
	name: string;
	description: string;
	startDate?: string;
	endDate?: string;
	createdAt: string;
	updatedAt: string;
}
```

### Transformer Function

```typescript
function toFrontendEpic(epic: BackendEpic): Epic {
	return {
		id: epic.id,
		projectId: epic.project_id,
		name: epic.name,
		description: epic.description,
		startDate: epic.start_date,
		endDate: epic.end_date,
		createdAt: epic.created_at,
		updatedAt: epic.updated_at
	};
}
```

### API Implementation

```typescript
import type { ApiResponse, IEpicsApi, CreateEpicInput, UpdateEpicInput } from '../types';
import type { Epic } from '$lib/types/api';
import { httpClient } from './client';

interface BackendEpic {
	id: string;
	project_id: string;
	name: string;
	description: string;
	start_date?: string;
	end_date?: string;
	created_at: string;
	updated_at: string;
}

function toFrontendEpic(epic: BackendEpic): Epic {
	return {
		id: epic.id,
		projectId: epic.project_id,
		name: epic.name,
		description: epic.description,
		startDate: epic.start_date,
		endDate: epic.end_date,
		createdAt: epic.created_at,
		updatedAt: epic.updated_at
	};
}

function createResponse<T>(data: T): ApiResponse<T> {
	return { data, error: null, status: 200 };
}

function createErrorResponse<T>(error: string, status: number): ApiResponse<T> {
	return { data: null, error, status };
}

async function listEpics(projectId: string): Promise<ApiResponse<Epic[]>> {
	try {
		const epics = await httpClient.get<BackendEpic[]>(`/projects/${projectId}/epics`);
		return createResponse(epics.map(toFrontendEpic));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load epics';
		return createErrorResponse(message, 500);
	}
}

async function createEpic(input: CreateEpicInput): Promise<ApiResponse<Epic>> {
	try {
		const epic = await httpClient.post<BackendEpic>(`/projects/${input.projectId}/epics`, {
			name: input.name,
			description: input.description
		});
		return createResponse(toFrontendEpic(epic));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to create epic';
		return createErrorResponse(message, 500);
	}
}

async function updateEpic(id: string, input: UpdateEpicInput): Promise<ApiResponse<Epic>> {
	try {
		const body: Record<string, unknown> = {};
		if (input.name) body.name = input.name;
		if (input.description) body.description = input.description;

		const epic = await httpClient.put<BackendEpic>(`/epics/${id}`, body);
		return createResponse(toFrontendEpic(epic));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to update epic';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function deleteEpic(id: string): Promise<ApiResponse<void>> {
	try {
		await httpClient.delete(`/epics/${id}`);
		return createResponse(undefined);
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to delete epic';
		return createErrorResponse(message, 500);
	}
}

export const realEpicsApi: IEpicsApi = {
	listEpics,
	createEpic,
	updateEpic,
	deleteEpic
};
```

## Files to Verify

| File | Status |
|------|--------|
| `frontend/src/lib/api/real/epics.ts` | EXISTS - Verify alignment |
| `frontend/src/lib/api/types.ts` | EXISTS - Verify IEpicsApi |
