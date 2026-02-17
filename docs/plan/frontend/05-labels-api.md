# Frontend Implementation Plan: Labels API

**Priority:** MEDIUM
**Estimated Time:** 20 minutes
**Backend Reference:** `docs/plan/backend/05-labels-handler.md`

## Current State

### What Exists
- `frontend/src/lib/api/real/labels.ts` - Implementation exists
- `frontend/src/lib/api/types.ts` - ILabelsApi interface

### What to Verify
- Labels are project-scoped (not community-scoped)
- Assignment endpoint uses /tasks/{taskId}/labels

## Backend API Contract

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/projects/{projectId}/labels` | - | `Label[]` |
| POST | `/projects/{projectId}/labels` | `{name, color}` | `Label` |
| POST | `/tasks/{taskId}/labels` | `{label_id}` | `Task` |

## Implementation Details

### Backend Response Format

```typescript
interface BackendLabel {
	id: string;
	project_id: string;
	name: string;
	color: string;
	created_at: string;
}
```

### Frontend Type Mapping

```typescript
// From $lib/types/api
interface Label {
	id: string;
	projectId: string;
	name: string;
	color: string;
	createdAt: string;
}
```

### Transformer Function

```typescript
function toFrontendLabel(label: BackendLabel): Label {
	return {
		id: label.id,
		projectId: label.project_id,
		name: label.name,
		color: label.color,
		createdAt: label.created_at
	};
}
```

### API Implementation

```typescript
import type { ApiResponse, ILabelsApi, CreateLabelInput } from '../types';
import type { Label, Task } from '$lib/types/api';
import { httpClient } from './client';

interface BackendLabel {
	id: string;
	project_id: string;
	name: string;
	color: string;
	created_at: string;
}

interface BackendTask {
	id: string;
	board_id: string;
	epic_id?: string;
	title: string;
	description: string;
	status: string;
	position: number;
	assignee_id?: string;
	reporter_id: string;
	created_at: string;
	updated_at: string;
	labels?: string[];
}

function toFrontendLabel(label: BackendLabel): Label {
	return {
		id: label.id,
		projectId: label.project_id,
		name: label.name,
		color: label.color,
		createdAt: label.created_at
	};
}

function toFrontendTask(task: BackendTask): Task {
	return {
		id: task.id,
		title: task.title,
		description: task.description,
		status: task.status as Task['status'],
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

async function listLabels(projectId: string): Promise<ApiResponse<Label[]>> {
	try {
		const labels = await httpClient.get<BackendLabel[]>(`/projects/${projectId}/labels`);
		return createResponse(labels.map(toFrontendLabel));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load labels';
		return createErrorResponse(message, 500);
	}
}

async function createLabel(input: CreateLabelInput): Promise<ApiResponse<Label>> {
	try {
		const label = await httpClient.post<BackendLabel>(`/projects/${input.communityId}/labels`, {
			name: input.name,
			color: input.color
		});
		return createResponse(toFrontendLabel(label));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to create label';
		return createErrorResponse(message, 500);
	}
}

async function assignLabelToTask(taskId: string, labelId: string): Promise<ApiResponse<Task>> {
	try {
		const task = await httpClient.post<BackendTask>(`/tasks/${taskId}/labels`, {
			label_id: labelId
		});
		return createResponse(toFrontendTask(task));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to assign label';
		return createErrorResponse(message, 500);
	}
}

export const realLabelsApi: ILabelsApi = {
	listLabels,
	createLabel,
	assignLabelToTask
};
```

## Note

The `CreateLabelInput.communityId` should be updated to `projectId` to match backend:
```typescript
export interface CreateLabelInput {
	name: string;
	color: string;
	projectId: string; // Changed from communityId
}
```

## Files to Verify

| File | Status |
|------|--------|
| `frontend/src/lib/api/real/labels.ts` | EXISTS - Verify alignment |
| `frontend/src/lib/api/types.ts` | EXISTS - Verify ILabelsApi |
