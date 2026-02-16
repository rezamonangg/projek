import type { ApiResponse } from '../types';
import type { Attachment } from '$lib/types/api';
import { mockResponse, generateId } from './utils';
import { generateMockAttachment } from './generators';

const mockAttachments: Attachment[] = [
	generateMockAttachment({
		id: 'attachment-1',
		filename: 'design-mockup.png',
		url: '/uploads/design-mockup.png',
		mimeType: 'image/png',
		size: 1024000,
		taskId: 'task-1',
		uploadedBy: 'user-1'
	}),
	generateMockAttachment({
		id: 'attachment-2',
		filename: 'requirements.pdf',
		url: '/uploads/requirements.pdf',
		mimeType: 'application/pdf',
		size: 512000,
		taskId: 'task-2',
		uploadedBy: 'user-2'
	})
];

export async function uploadFile(
	file: File,
	options: { taskId?: string; wikiPageId?: string; uploadedBy: string }
): Promise<ApiResponse<Attachment>> {
	const newAttachment = generateMockAttachment({
		id: generateId(),
		filename: file.name,
		url: `/uploads/${file.name}`,
		mimeType: file.type,
		size: file.size,
		taskId: options.taskId,
		wikiPageId: options.wikiPageId,
		uploadedBy: options.uploadedBy
	});
	mockAttachments.push(newAttachment);
	return mockResponse(newAttachment);
}

export async function downloadFile(id: string): Promise<ApiResponse<Blob>> {
	const attachment = mockAttachments.find((a) => a.id === id);
	if (!attachment) {
		return mockResponse(null as unknown as Blob);
	}

	const blob = new Blob(['Mock file content'], { type: attachment.mimeType });
	return mockResponse(blob);
}
