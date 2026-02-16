import type { ApiResponse, IFilesApi, UploadFileOptions, ListAttachmentsOptions } from '../types';
import type { Attachment } from '$lib/types/api';
import { mockResponse, mockError, generateId } from './utils';
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

async function uploadFile(
	file: File,
	options: UploadFileOptions
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

async function downloadFile(id: string): Promise<ApiResponse<Blob>> {
	const attachment = mockAttachments.find((a) => a.id === id);
	if (!attachment) {
		return mockError('File not found', 404);
	}

	const blob = new Blob(['Mock file content'], { type: attachment.mimeType });
	return mockResponse(blob);
}

async function listAttachments(options: ListAttachmentsOptions): Promise<ApiResponse<Attachment[]>> {
	let attachments = mockAttachments;
	if (options.taskId) {
		attachments = attachments.filter((a) => a.taskId === options.taskId);
	}
	if (options.wikiPageId) {
		attachments = attachments.filter((a) => a.wikiPageId === options.wikiPageId);
	}
	return mockResponse(attachments);
}

export const mockFilesApi: IFilesApi = {
	uploadFile,
	downloadFile,
	listAttachments
};
