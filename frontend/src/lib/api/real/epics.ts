import type { ApiResponse, IEpicsApi, CreateEpicInput, UpdateEpicInput } from '../types';
import type { Epic } from '$lib/types/api';
import { httpClient } from './client';

interface BackendEpic {
	id: string;
	project_id: string;
	name: string;
	description: string;
	start_date?: string;
	end_date?: string;
	created_at: string;
	updated_at: string;
}

function toFrontendEpic(epic: BackendEpic): Epic {
	const hash = epic.id.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
	const colors = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#EC4899'];
	const color = colors[hash % colors.length];

	return {
		id: epic.id,
		name: epic.name,
		description: epic.description,
		color,
		projectId: epic.project_id,
		createdAt: epic.created_at
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

async function listEpics(projectId: string): Promise<ApiResponse<Epic[]>> {
	try {
		const epics = await httpClient.get<BackendEpic[]>(`/projects/${projectId}/epics`);
		return createResponse(epics.map(toFrontendEpic));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load epics';
		return createErrorResponse(message, 500);
	}
}

async function createEpic(input: CreateEpicInput): Promise<ApiResponse<Epic>> {
	try {
		const epic = await httpClient.post<BackendEpic>(`/projects/${input.projectId}/epics`, {
			name: input.name,
			description: input.description
		});
		return createResponse(toFrontendEpic(epic));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to create epic';
		return createErrorResponse(message, 500);
	}
}

async function updateEpic(id: string, input: UpdateEpicInput): Promise<ApiResponse<Epic>> {
	try {
		const epic = await httpClient.put<BackendEpic>(`/epics/${id}`, {
			name: input.name,
			description: input.description
		});
		return createResponse(toFrontendEpic(epic));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to update epic';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function deleteEpic(id: string): Promise<ApiResponse<void>> {
	try {
		await httpClient.delete(`/epics/${id}`);
		return createResponse(undefined);
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to delete epic';
		return createErrorResponse(message, 500);
	}
}

export const realEpicsApi: IEpicsApi = {
	listEpics,
	createEpic,
	updateEpic,
	deleteEpic
};
