# Frontend Implementation Plan: Files API

**Priority:** LOW
**Estimated Time:** 45 minutes
**Backend Reference:** `docs/plan/backend/08-files-module.md`

## Current State

### What Exists
- `frontend/src/lib/api/real/files.ts` - Implementation exists
- `frontend/src/lib/api/types.ts` - IFilesApi interface

### What to Verify
- Upload uses FormData (not JSON)
- Download returns Blob
- File size limit handling (10MB)

## Backend API Contract

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| POST | `/files/upload` | FormData: file, uploader_id, project_id, task_id? | `Attachment` |
| GET | `/files/{id}` | - | File blob |
| GET | `/files?task_id=...` | - | `Attachment[]` |

## Implementation Details

### Backend Response Format

```typescript
interface BackendAttachment {
	id: string;
	project_id: string;
	task_id?: string;
	uploader_id: string;
	filename: string;
	original_filename: string;
	content_type: string;
	size: number;
	storage_path: string;
	storage_type: string;
	created_at: string;
}
```

### Frontend Type Mapping

```typescript
// From $lib/types/api
interface Attachment {
	id: string;
	projectId: string;
	taskId?: string;
	uploaderId: string;
	filename: string;
	originalFilename: string;
	contentType: string;
	size: number;
	storagePath: string;
	storageType: string;
	createdAt: string;
}
```

### Transformer Function

```typescript
function toFrontendAttachment(attachment: BackendAttachment): Attachment {
	return {
		id: attachment.id,
		projectId: attachment.project_id,
		taskId: attachment.task_id,
		uploaderId: attachment.uploader_id,
		filename: attachment.filename,
		originalFilename: attachment.original_filename,
		contentType: attachment.content_type,
		size: attachment.size,
		storagePath: attachment.storage_path,
		storageType: attachment.storage_type,
		createdAt: attachment.created_at
	};
}
```

### API Implementation

```typescript
import type { ApiResponse, IFilesApi, UploadFileOptions, ListAttachmentsOptions } from '../types';
import type { Attachment } from '$lib/types/api';
import { PUBLIC_API_URL } from '$env/static/public';

interface BackendAttachment {
	id: string;
	project_id: string;
	task_id?: string;
	uploader_id: string;
	filename: string;
	original_filename: string;
	content_type: string;
	size: number;
	storage_path: string;
	storage_type: string;
	created_at: string;
}

function toFrontendAttachment(attachment: BackendAttachment): Attachment {
	return {
		id: attachment.id,
		projectId: attachment.project_id,
		taskId: attachment.task_id,
		uploaderId: attachment.uploader_id,
		filename: attachment.filename,
		originalFilename: attachment.original_filename,
		contentType: attachment.content_type,
		size: attachment.size,
		storagePath: attachment.storage_path,
		storageType: attachment.storage_type,
		createdAt: attachment.created_at
	};
}

function createResponse<T>(data: T): ApiResponse<T> {
	return { data, error: null, status: 200 };
}

function createErrorResponse<T>(error: string, status: number): ApiResponse<T> {
	return { data: null, error, status };
}

const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB

async function uploadFile(file: File, options: UploadFileOptions): Promise<ApiResponse<Attachment>> {
	try {
		// Validate file size
		if (file.size > MAX_FILE_SIZE) {
			return createErrorResponse('File size exceeds 10MB limit', 400);
		}

		const formData = new FormData();
		formData.append('file', file);
		formData.append('uploader_id', options.uploadedBy);
		if (options.taskId) {
			formData.append('task_id', options.taskId);
		}

		const response = await fetch(`${PUBLIC_API_URL}/files/upload`, {
			method: 'POST',
			body: formData,
			credentials: 'include'
		});

		if (!response.ok) {
			let errorMessage = `HTTP ${response.status}`;
			try {
				const errorData = await response.json();
				if (errorData.error?.message) {
					errorMessage = errorData.error.message;
				}
			} catch {
				// Ignore JSON parse errors
			}
			return createErrorResponse(errorMessage, response.status);
		}

		const attachment: BackendAttachment = await response.json();
		return createResponse(toFrontendAttachment(attachment));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to upload file';
		return createErrorResponse(message, 500);
	}
}

async function downloadFile(id: string): Promise<ApiResponse<Blob>> {
	try {
		const response = await fetch(`${PUBLIC_API_URL}/files/${id}`, {
			method: 'GET',
			credentials: 'include'
		});

		if (!response.ok) {
			const message = response.status === 404 ? 'File not found' : `HTTP ${response.status}`;
			return createErrorResponse(message, response.status);
		}

		const blob = await response.blob();
		return createResponse(blob);
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to download file';
		return createErrorResponse(message, 500);
	}
}

async function listAttachments(options: ListAttachmentsOptions): Promise<ApiResponse<Attachment[]>> {
	try {
		const params = new URLSearchParams();
		if (options.taskId) {
			params.append('task_id', options.taskId);
		}

		const queryString = params.toString();
		const url = queryString ? `/files?${queryString}` : '/files';
		
		const attachments = await fetch(`${PUBLIC_API_URL}${url}`, {
			method: 'GET',
			credentials: 'include'
		}).then(res => res.json() as Promise<BackendAttachment[]>);

		return createResponse(attachments.map(toFrontendAttachment));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to list attachments';
		return createErrorResponse(message, 500);
	}
}

export const realFilesApi: IFilesApi = {
	uploadFile,
	downloadFile,
	listAttachments
};
```

## Key Implementation Notes

1. **FormData for upload**: Unlike other APIs, file upload uses `FormData` instead of JSON
2. **No Content-Type header**: Let the browser set it for FormData
3. **Blob for download**: Returns raw file data, not JSON
4. **File size validation**: Client-side check before upload (10MB limit)

## Error Handling

| Status Code | Scenario | Action |
|-------------|----------|--------|
| 400 | File too large / invalid | Show size/type error |
| 401 | Unauthorized | Redirect to login |
| 404 | File not found | Show not found message |
| 413 | Payload too large | Show size limit error |
| 500 | Server error | Show generic error |

## Files to Verify

| File | Status |
|------|--------|
| `frontend/src/lib/api/real/files.ts` | EXISTS - Verify alignment |
| `frontend/src/lib/api/types.ts` | EXISTS - Verify IFilesApi |
