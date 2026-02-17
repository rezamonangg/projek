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
	let content: Record<string, unknown> = { type: 'doc', content: [] };
	try {
		if (page.content) {
			content = JSON.parse(page.content);
		}
	} catch {
		content = { type: 'doc', content: [{ type: 'paragraph', content: [{ type: 'text', text: page.content }] }] };
	}

	return {
		id: page.id,
		title: page.title,
		slug: page.slug,
		content,
		projectId: page.project_id,
		parentId: page.parent_id,
		createdAt: page.created_at,
		updatedAt: page.updated_at
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
