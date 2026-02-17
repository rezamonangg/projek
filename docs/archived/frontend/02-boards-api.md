# Frontend Implementation Plan: Boards API

**Priority:** HIGH
**Estimated Time:** 20 minutes
**Backend Reference:** `docs/plan/backend/03-boards-handler.md`

## Current State

### What Exists
- `frontend/src/lib/api/real/boards.ts` - Implementation exists
- `frontend/src/lib/api/types.ts` - IBoardsApi interface

### What to Verify
- Endpoints use correct URL patterns
- Backend response format matches

## Backend API Contract

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/projects/{projectId}/boards` | - | `Board[]` |
| POST | `/projects/{projectId}/boards` | `{name}` | `Board` |
| GET | `/boards/{id}` | - | `Board` |
| DELETE | `/boards/{id}` | - | `void` |

## Implementation Details

### Backend Response Format

```typescript
interface BackendBoard {
	id: string;
	project_id: string;
	name: string;
	created_at: string;
	updated_at: string;
}
```

### Frontend Type Mapping

```typescript
// From $lib/types/api
interface Board {
	id: string;
	projectId: string;
	name: string;
	createdAt: string;
	updatedAt: string;
}
```

### Transformer Function

```typescript
function toFrontendBoard(board: BackendBoard): Board {
	return {
		id: board.id,
		projectId: board.project_id,
		name: board.name,
		createdAt: board.created_at,
		updatedAt: board.updated_at
	};
}
```

### API Implementation

```typescript
import type { ApiResponse, IBoardsApi, CreateBoardInput } from '../types';
import type { Board } from '$lib/types/api';
import { httpClient } from './client';

interface BackendBoard {
	id: string;
	project_id: string;
	name: string;
	created_at: string;
	updated_at: string;
}

function toFrontendBoard(board: BackendBoard): Board {
	return {
		id: board.id,
		projectId: board.project_id,
		name: board.name,
		createdAt: board.created_at,
		updatedAt: board.updated_at
	};
}

function createResponse<T>(data: T): ApiResponse<T> {
	return { data, error: null, status: 200 };
}

function createErrorResponse<T>(error: string, status: number): ApiResponse<T> {
	return { data: null, error, status };
}

async function listBoards(projectId: string): Promise<ApiResponse<Board[]>> {
	try {
		const boards = await httpClient.get<BackendBoard[]>(`/projects/${projectId}/boards`);
		return createResponse(boards.map(toFrontendBoard));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load boards';
		return createErrorResponse(message, 500);
	}
}

async function getBoard(id: string): Promise<ApiResponse<Board>> {
	try {
		const board = await httpClient.get<BackendBoard>(`/boards/${id}`);
		return createResponse(toFrontendBoard(board));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load board';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function createBoard(input: CreateBoardInput): Promise<ApiResponse<Board>> {
	try {
		const board = await httpClient.post<BackendBoard>(`/projects/${input.projectId}/boards`, {
			name: input.name,
			description: input.description
		});
		return createResponse(toFrontendBoard(board));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to create board';
		return createErrorResponse(message, 500);
	}
}

async function deleteBoard(id: string): Promise<ApiResponse<void>> {
	try {
		await httpClient.delete(`/boards/${id}`);
		return createResponse(undefined);
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to delete board';
		return createErrorResponse(message, 500);
	}
}

export const realBoardsApi: IBoardsApi = {
	listBoards,
	getBoard,
	createBoard,
	deleteBoard
};
```

## Error Handling

| Status Code | Scenario | Action |
|-------------|----------|--------|
| 400 | Validation error | Show validation message |
| 401 | Unauthorized | Redirect to login |
| 404 | Board not found | Show not found message |
| 500 | Server error | Show generic error |

## Files to Verify

| File | Status |
|------|--------|
| `frontend/src/lib/api/real/boards.ts` | EXISTS - Verify alignment |
| `frontend/src/lib/api/types.ts` | EXISTS - Verify IBoardsApi |
