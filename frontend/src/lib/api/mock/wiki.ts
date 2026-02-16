import type { ApiResponse } from '../types';
import type { WikiPage } from '$lib/types/api';
import { mockResponse, generateId } from './utils';
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

export async function listWikiPages(projectId: string): Promise<ApiResponse<WikiPage[]>> {
	const pages = mockWikiPages.filter((p) => p.projectId === projectId);
	return mockResponse(pages);
}

export async function getWikiPage(id: string): Promise<ApiResponse<WikiPage>> {
	const page = mockWikiPages.find((p) => p.id === id);
	if (!page) {
		return mockResponse(null as unknown as WikiPage);
	}
	return mockResponse(page);
}

interface CreateWikiPageInput {
	title: string;
	slug: string;
	content: Record<string, unknown>;
	projectId: string;
	parentId?: string;
}

export async function createWikiPage(input: CreateWikiPageInput): Promise<ApiResponse<WikiPage>> {
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
