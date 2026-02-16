import type { ApiResponse, IEpicsApi, CreateEpicInput, UpdateEpicInput } from '../types';
import type { Epic } from '$lib/types/api';
import { mockResponse, mockError, generateId } from './utils';
import { generateMockEpic } from './generators';

const mockEpics: Epic[] = [
	generateMockEpic({
		id: 'epic-1',
		name: 'Homepage',
		description: 'Homepage redesign and implementation',
		color: '#3B82F6',
		projectId: 'project-1'
	}),
	generateMockEpic({
		id: 'epic-2',
		name: 'User Dashboard',
		description: 'User dashboard features',
		color: '#10B981',
		projectId: 'project-1'
	}),
	generateMockEpic({
		id: 'epic-3',
		name: 'Authentication',
		description: 'Authentication and authorization',
		color: '#F59E0B',
		projectId: 'project-2'
	})
];

async function listEpics(projectId: string): Promise<ApiResponse<Epic[]>> {
	const epics = mockEpics.filter((e) => e.projectId === projectId);
	return mockResponse(epics);
}

async function createEpic(input: CreateEpicInput): Promise<ApiResponse<Epic>> {
	const newEpic = generateMockEpic({
		id: generateId(),
		name: input.name,
		description: input.description,
		color: input.color,
		projectId: input.projectId
	});
	mockEpics.push(newEpic);
	return mockResponse(newEpic);
}

async function updateEpic(id: string, input: UpdateEpicInput): Promise<ApiResponse<Epic>> {
	const index = mockEpics.findIndex((e) => e.id === id);
	if (index === -1) {
		return mockError('Epic not found', 404);
	}

	mockEpics[index] = {
		...mockEpics[index],
		...input
	};

	return mockResponse(mockEpics[index]);
}

async function deleteEpic(id: string): Promise<ApiResponse<void>> {
	const index = mockEpics.findIndex((e) => e.id === id);
	if (index !== -1) {
		mockEpics.splice(index, 1);
	}
	return mockResponse(undefined);
}

export const mockEpicsApi: IEpicsApi = {
	listEpics,
	createEpic,
	updateEpic,
	deleteEpic
};
