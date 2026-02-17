# Frontend Implementation Plan: Tasks API

**Priority:** HIGH
**Estimated Time:** 30 minutes
**Backend Reference:** `docs/plan/backend/04-tasks-handler.md`

## Current State

### What Exists
- `frontend/src/lib/api/real/tasks.ts` - Full implementation
- `frontend/src/lib/api/types.ts` - ITasksApi interface

### What to Verify
- Status enum mapping (backend has 7 statuses)
- Move endpoint uses PATCH
- All status transitions work

## Backend API Contract

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/boards/{boardId}/tasks` | - | `Task[]` |
| POST | `/boards/{boardId}/tasks` | `{title, description, epic_id?, assignee_id?, reporter_id, status}` | `Task` |
| PUT | `/tasks/{id}` | `{title?, description?, status?, epic_id?, assignee_id?}` | `Task` |
| PATCH | `/tasks/{id}/move` | `{status, position}` | `Task` |
| DELETE | `/tasks/{id}` | - | `void` |

## Status Mapping

Backend has 7 statuses:
```typescript
type BackendStatus = 'backlog' | 'todo' | 'inprogress' | 'codereview' | 'intest' | 'needdeploy' | 'done';
```

Frontend groups some statuses:
```typescript
type FrontendStatus = 'backlog' | 'todo' | 'inprogress' | 'done';
```

### Status Mapping Logic

```typescript
function toFrontendStatus(status: BackendStatus): FrontendStatus {
	switch (status) {
		case 'done':
			return 'done';
		case 'inprogress':
		case 'codereview':
		case 'intest':
		case 'needdeploy':
			return 'inprogress';
		case 'todo':
			return 'todo';
		default:
			return 'backlog';
	}
}

function toBackendStatus(status: FrontendStatus): BackendStatus {
	return status as BackendStatus;
}
```

## Implementation Details

### Backend Response Format

```typescript
type BackendStatus = 'backlog' | 'todo' | 'inprogress' | 'codereview' | 'intest' | 'needdeploy' | 'done';

interface BackendTask {
	id: string;
	board_id: string;
	epic_id?: string;
	title: string;
	description: string;
	status: BackendStatus;
	position: number;
	story_points?: number;
	due_date?: string;
	assignee_id?: string;
	reporter_id: string;
	created_at: string;
	updated_at: string;
	labels?: string[];
}
```

### Frontend Type Mapping

```typescript
// From $lib/types/api
interface Task {
	id: string;
	title: string;
	description: string;
	status: 'backlog' | 'todo' | 'inprogress' | 'done';
	boardId: string;
	epicId?: string;
	assigneeId?: string;
	labels: string[];
	position: number;
	createdAt: string;
	updatedAt: string;
}
```

### API Implementation

```typescript
import type { ApiResponse, ITasksApi, CreateTaskInput, UpdateTaskInput, MoveTaskInput } from '../types';
import type { Task } from '$lib/types/api';
import { httpClient } from './client';

type BackendStatus = 'backlog' | 'todo' | 'inprogress' | 'codereview' | 'intest' | 'needdeploy' | 'done';

interface BackendTask {
	id: string;
	board_id: string;
	epic_id?: string;
	title: string;
	description: string;
	status: BackendStatus;
	position: number;
	story_points?: number;
	due_date?: string;
	assignee_id?: string;
	reporter_id: string;
	created_at: string;
	updated_at: string;
	labels?: string[];
}

function toFrontendStatus(status: BackendStatus): Task['status'] {
	switch (status) {
		case 'done':
			return 'done';
		case 'inprogress':
		case 'codereview':
		case 'intest':
		case 'needdeploy':
			return 'inprogress';
		case 'todo':
			return 'todo';
		default:
			return 'backlog';
	}
}

function toBackendStatus(status: Task['status']): BackendStatus {
	return status as BackendStatus;
}

function toFrontendTask(task: BackendTask): Task {
	return {
		id: task.id,
		title: task.title,
		description: task.description,
		status: toFrontendStatus(task.status),
		boardId: task.board_id,
		epicId: task.epic_id,
		assigneeId: task.assignee_id,
		labels: task.labels || [],
		position: task.position,
		createdAt: task.created_at,
		updatedAt: task.updated_at
	};
}

function createResponse<T>(data: T): ApiResponse<T> {
	return { data, error: null, status: 200 };
}

function createErrorResponse<T>(error: string, status: number): ApiResponse<T> {
	return { data: null, error, status };
}

async function listTasks(boardId: string): Promise<ApiResponse<Task[]>> {
	try {
		const tasks = await httpClient.get<BackendTask[]>(`/boards/${boardId}/tasks`);
		return createResponse(tasks.map(toFrontendTask));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load tasks';
		return createErrorResponse(message, 500);
	}
}

async function createTask(input: CreateTaskInput): Promise<ApiResponse<Task>> {
	try {
		const task = await httpClient.post<BackendTask>(`/boards/${input.boardId}/tasks`, {
			title: input.title,
			description: input.description,
			epic_id: input.epicId,
			assignee_id: input.assigneeId,
			status: 'backlog',
			reporter_id: input.assigneeId || '00000000-0000-0000-0000-000000000000'
		});
		return createResponse(toFrontendTask(task));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to create task';
		return createErrorResponse(message, 500);
	}
}

async function updateTask(id: string, input: UpdateTaskInput): Promise<ApiResponse<Task>> {
	try {
		const body: Record<string, unknown> = {};
		if (input.title) body.title = input.title;
		if (input.description) body.description = input.description;
		if (input.status) body.status = toBackendStatus(input.status);
		if (input.epicId) body.epic_id = input.epicId;
		if (input.assigneeId) body.assignee_id = input.assigneeId;

		const task = await httpClient.put<BackendTask>(`/tasks/${id}`, body);
		return createResponse(toFrontendTask(task));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to update task';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function moveTask(id: string, input: MoveTaskInput): Promise<ApiResponse<Task>> {
	try {
		const task = await httpClient.patch<BackendTask>(`/tasks/${id}/move`, {
			status: toBackendStatus(input.status),
			position: input.position
		});
		return createResponse(toFrontendTask(task));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to move task';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function deleteTask(id: string): Promise<ApiResponse<void>> {
	try {
		await httpClient.delete(`/tasks/${id}`);
		return createResponse(undefined);
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to delete task';
		return createErrorResponse(message, 500);
	}
}

export const realTasksApi: ITasksApi = {
	listTasks,
	createTask,
	updateTask,
	moveTask,
	deleteTask
};
```

## Error Handling

| Status Code | Scenario | Action |
|-------------|----------|--------|
| 400 | Validation error | Show validation message |
| 401 | Unauthorized | Redirect to login |
| 404 | Task not found | Show not found message |
| 500 | Server error | Show generic error |

## Files to Verify

| File | Status |
|------|--------|
| `frontend/src/lib/api/real/tasks.ts` | EXISTS - Verify alignment |
| `frontend/src/lib/api/types.ts` | EXISTS - Verify ITasksApi |
