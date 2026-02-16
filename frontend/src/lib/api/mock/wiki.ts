import type { ApiResponse, IWikiApi, CreateWikiPageInput, UpdateWikiPageInput } from '../types';
import type { WikiPage } from '$lib/types/api';
import { mockResponse, mockError, generateId } from './utils';
import { generateMockWikiPage } from './generators';

const mockWikiPages: WikiPage[] = [
	generateMockWikiPage({
		id: 'wiki-1',
		title: 'Getting Started',
		slug: 'getting-started',
		content: { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: 'Welcome to Projek!' }] }] },
		projectId: 'project-1'
	}),
	generateMockWikiPage({
		id: 'wiki-2',
		title: 'Architecture',
		slug: 'architecture',
		content: { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: 'System architecture overview' }] }] },
		projectId: 'project-1',
		parentId: 'wiki-1'
	}),
	generateMockWikiPage({
		id: 'wiki-3',
		title: 'API Reference',
		slug: 'api-reference',
		content: { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: 'API documentation' }] }] },
		projectId: 'project-1'
	})
];

async function listWikiPages(projectId: string): Promise<ApiResponse<WikiPage[]>> {
	const pages = mockWikiPages.filter((p) => p.projectId === projectId);
	return mockResponse(pages);
}

async function getWikiPage(id: string): Promise<ApiResponse<WikiPage>> {
	const page = mockWikiPages.find((p) => p.id === id);
	if (!page) {
		return mockError('Wiki page not found', 404);
	}
	return mockResponse(page);
}

async function createWikiPage(input: CreateWikiPageInput): Promise<ApiResponse<WikiPage>> {
	const newPage = generateMockWikiPage({
		id: generateId(),
		title: input.title,
		slug: input.slug,
		content: input.content,
		projectId: input.projectId,
		parentId: input.parentId
	});
	mockWikiPages.push(newPage);
	return mockResponse(newPage);
}

async function updateWikiPage(
	id: string,
	input: UpdateWikiPageInput
): Promise<ApiResponse<WikiPage>> {
	const index = mockWikiPages.findIndex((p) => p.id === id);
	if (index === -1) {
		return mockError('Wiki page not found', 404);
	}

	mockWikiPages[index] = {
		...mockWikiPages[index],
		...input,
		updatedAt: new Date().toISOString()
	};

	return mockResponse(mockWikiPages[index]);
}

async function deleteWikiPage(id: string): Promise<ApiResponse<void>> {
	const index = mockWikiPages.findIndex((p) => p.id === id);
	if (index !== -1) {
		mockWikiPages.splice(index, 1);
	}
	return mockResponse(undefined);
}

export const mockWikiApi: IWikiApi = {
	listWikiPages,
	getWikiPage,
	createWikiPage,
	updateWikiPage,
	deleteWikiPage
};
