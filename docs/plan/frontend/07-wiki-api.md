# Frontend Implementation Plan: Wiki API

**Priority:** LOW
**Estimated Time:** 30 minutes
**Backend Reference:** `docs/plan/backend/07-wiki-module.md`

## Current State

### What Exists
- `frontend/src/lib/api/real/wiki.ts` - Implementation exists
- `frontend/src/lib/api/types.ts` - IWikiApi interface

### What to Verify
- Content is stored as JSON string
- Nested pages via parent_id

## Backend API Contract

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/projects/{projectId}/wiki` | - | `WikiPage[]` |
| POST | `/projects/{projectId}/wiki` | `{title, content, parent_id?}` | `WikiPage` |
| GET | `/wiki/{id}` | - | `WikiPage` |
| PUT | `/wiki/{id}` | `{title?, content?, parent_id?}` | `WikiPage` |
| DELETE | `/wiki/{id}` | - | `void` |

## Implementation Details

### Backend Response Format

```typescript
interface BackendWikiPage {
	id: string;
	project_id: string;
	parent_id?: string;
	title: string;
	slug: string;
	content: string; // JSON string
	created_at: string;
	updated_at: string;
}
```

### Frontend Type Mapping

```typescript
// From $lib/types/api
interface WikiPage {
	id: string;
	projectId: string;
	parentId?: string;
	title: string;
	slug: string;
	content: Record<string, unknown>; // Parsed JSON
	createdAt: string;
	updatedAt: string;
}
```

### Transformer Function

```typescript
function toFrontendWikiPage(page: BackendWikiPage): WikiPage {
	return {
		id: page.id,
		projectId: page.project_id,
		parentId: page.parent_id,
		title: page.title,
		slug: page.slug,
		content: typeof page.content === 'string' ? JSON.parse(page.content) : page.content,
		createdAt: page.created_at,
		updatedAt: page.updated_at
	};
}

function toBackendContent(content: Record<string, unknown>): string {
	return JSON.stringify(content);
}
```

### API Implementation

```typescript
import type { ApiResponse, IWikiApi, CreateWikiPageInput, UpdateWikiPageInput } from '../types';
import type { WikiPage } from '$lib/types/api';
import { httpClient } from './client';

interface BackendWikiPage {
	id: string;
	project_id: string;
	parent_id?: string;
	title: string;
	slug: string;
	content: string;
	created_at: string;
	updated_at: string;
}

function toFrontendWikiPage(page: BackendWikiPage): WikiPage {
	let parsedContent: Record<string, unknown> = {};
	try {
		parsedContent = typeof page.content === 'string' ? JSON.parse(page.content) : page.content;
	} catch {
		parsedContent = {};
	}
	
	return {
		id: page.id,
		projectId: page.project_id,
		parentId: page.parent_id,
		title: page.title,
		slug: page.slug,
		content: parsedContent,
		createdAt: page.created_at,
		updatedAt: page.updated_at
	};
}

function createResponse<T>(data: T): ApiResponse<T> {
	return { data, error: null, status: 200 };
}

function createErrorResponse<T>(error: string, status: number): ApiResponse<T> {
	return { data: null, error, status };
}

async function listWikiPages(projectId: string): Promise<ApiResponse<WikiPage[]>> {
	try {
		const pages = await httpClient.get<BackendWikiPage[]>(`/projects/${projectId}/wiki`);
		return createResponse(pages.map(toFrontendWikiPage));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load wiki pages';
		return createErrorResponse(message, 500);
	}
}

async function getWikiPage(id: string): Promise<ApiResponse<WikiPage>> {
	try {
		const page = await httpClient.get<BackendWikiPage>(`/wiki/${id}`);
		return createResponse(toFrontendWikiPage(page));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load wiki page';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function createWikiPage(input: CreateWikiPageInput): Promise<ApiResponse<WikiPage>> {
	try {
		const page = await httpClient.post<BackendWikiPage>(`/projects/${input.projectId}/wiki`, {
			title: input.title,
			content: JSON.stringify(input.content),
			parent_id: input.parentId
		});
		return createResponse(toFrontendWikiPage(page));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to create wiki page';
		return createErrorResponse(message, 500);
	}
}

async function updateWikiPage(id: string, input: UpdateWikiPageInput): Promise<ApiResponse<WikiPage>> {
	try {
		const body: Record<string, unknown> = {};
		if (input.title) body.title = input.title;
		if (input.slug) body.slug = input.slug;
		if (input.content) body.content = JSON.stringify(input.content);
		if (input.parentId !== undefined) body.parent_id = input.parentId;

		const page = await httpClient.put<BackendWikiPage>(`/wiki/${id}`, body);
		return createResponse(toFrontendWikiPage(page));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to update wiki page';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function deleteWikiPage(id: string): Promise<ApiResponse<void>> {
	try {
		await httpClient.delete(`/wiki/${id}`);
		return createResponse(undefined);
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to delete wiki page';
		return createErrorResponse(message, 500);
	}
}

export const realWikiApi: IWikiApi = {
	listWikiPages,
	getWikiPage,
	createWikiPage,
	updateWikiPage,
	deleteWikiPage
};
```

## Files to Verify

| File | Status |
|------|--------|
| `frontend/src/lib/api/real/wiki.ts` | EXISTS - Verify alignment |
| `frontend/src/lib/api/types.ts` | EXISTS - Verify IWikiApi |
