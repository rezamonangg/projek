import type { ApiResponse, IFilesApi, UploadFileOptions, ListAttachmentsOptions } from '../types';
import type { Attachment } from '$lib/types/api';
import { PUBLIC_API_URL } from '$env/static/public';

interface BackendFileAttachment {
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
	created_at: number;
}

function toFrontendAttachment(file: BackendFileAttachment): Attachment {
	return {
		id: file.id,
		filename: file.original_filename,
		url: `/files/${file.id}`,
		mimeType: file.content_type,
		size: file.size,
		taskId: file.task_id,
		uploadedBy: file.uploader_id,
		createdAt: new Date(file.created_at * 1000).toISOString()
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

const MAX_FILE_SIZE = 10 * 1024 * 1024;

async function uploadFile(file: File, options: UploadFileOptions): Promise<ApiResponse<Attachment>> {
	try {
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
			const errorData = await response.json().catch(() => ({}));
			throw new Error(errorData.error?.message || `Upload failed: ${response.status}`);
		}

		const data = await response.json();
		return createResponse(toFrontendAttachment(data));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to upload file';
		return createErrorResponse(message, 500);
	}
}

async function downloadFile(id: string): Promise<ApiResponse<Blob>> {
	try {
		const response = await fetch(`${PUBLIC_API_URL}/files/${id}`, {
			credentials: 'include'
		});

		if (!response.ok) {
			throw new Error(`Download failed: ${response.status}`);
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
		const endpoint = queryString ? `/files?${queryString}` : '/files';

		const response = await fetch(`${PUBLIC_API_URL}${endpoint}`, {
			credentials: 'include'
		});

		if (!response.ok) {
			throw new Error(`Failed to list files: ${response.status}`);
		}

		const json = await response.json();
		const data = json.data || json;
		const files = Array.isArray(data) ? data : data.items || [];
		return createResponse(files.map(toFrontendAttachment));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to list files';
		return createErrorResponse(message, 500);
	}
}

export const realFilesApi: IFilesApi = {
	uploadFile,
	downloadFile,
	listAttachments
};
