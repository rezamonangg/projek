import type { ApiResponse } from '../types';
import type { Epic } from '$lib/types/api';
import { mockResponse } from './utils';
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

export async function listEpics(projectId: string): Promise<ApiResponse<Epic[]>> {
	const epics = mockEpics.filter((e) => e.projectId === projectId);
	return mockResponse(epics);
}
