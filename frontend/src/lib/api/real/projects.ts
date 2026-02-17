import type { ApiResponse, PaginatedResponse, IProjectsApi, CreateProjectInput, UpdateProjectInput } from '../types';
import type { Project } from '$lib/types/api';
import { httpClient } from './client';

interface BackendProject {
	id: string;
	community_id: string;
	name: string;
	description: string;
	key: string;
	is_archived: boolean;
	created_at: string;
	updated_at: string;
}

interface BackendPaginatedProjects {
	items: BackendProject[];
	total: number;
	page: number;
	page_size: number;
	total_pages: number;
}

function toFrontendProject(project: BackendProject): Project {
	return {
		id: project.id,
		name: project.name,
		description: project.description,
		status: project.is_archived ? 'archived' : 'active',
		communityId: project.community_id,
		createdAt: project.created_at,
		updatedAt: project.updated_at
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

async function listProjects(): Promise<ApiResponse<PaginatedResponse<Project>>> {
	try {
		const response = await httpClient.get<BackendPaginatedProjects>('/projects');
		const items = response.items.map(toFrontendProject);
		return createResponse({
			items,
			total: response.total,
			page: response.page,
			pageSize: response.page_size,
			totalPages: response.total_pages
		});
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load projects';
		return createErrorResponse(message, 500);
	}
}

async function getProject(id: string): Promise<ApiResponse<Project>> {
	try {
		const project = await httpClient.get<BackendProject>(`/projects/${id}`);
		return createResponse(toFrontendProject(project));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load project';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function createProject(input: CreateProjectInput): Promise<ApiResponse<Project>> {
	try {
		const baseKey = input.name.replace(/[^a-zA-Z]/g, '').substring(0, 4).toUpperCase();
		const suffix = Math.random().toString(36).substring(2, 5).toUpperCase();
		const key = baseKey.padEnd(2, 'X').substring(0, 2) + suffix;
		const project = await httpClient.post<BackendProject>('/projects', {
			name: input.name,
			description: input.description,
			key
		});
		return createResponse(toFrontendProject(project));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to create project';
		return createErrorResponse(message, 500);
	}
}

async function updateProject(id: string, input: UpdateProjectInput): Promise<ApiResponse<Project>> {
	try {
		const body: Record<string, unknown> = {};
		if (input.name) body.name = input.name;
		if (input.description) body.description = input.description;
		if (input.status) body.is_archived = input.status === 'archived';

		const project = await httpClient.put<BackendProject>(`/projects/${id}`, body);
		return createResponse(toFrontendProject(project));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to update project';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function deleteProject(id: string): Promise<ApiResponse<void>> {
	try {
		await httpClient.delete(`/projects/${id}`);
		return createResponse(undefined);
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to delete project';
		return createErrorResponse(message, 500);
	}
}

export const realProjectsApi: IProjectsApi = {
	listProjects,
	getProject,
	createProject,
	updateProject,
	deleteProject
};
