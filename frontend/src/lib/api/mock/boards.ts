import type { ApiResponse } from '../types';
import type { Board } from '$lib/types/api';
import { mockResponse } from './utils';
import { generateMockBoard } from './generators';

const mockBoards: Board[] = [
	generateMockBoard({
		id: 'board-1',
		name: 'Main Board',
		description: 'Primary kanban board for the project',
		projectId: 'project-1'
	}),
	generateMockBoard({
		id: 'board-2',
		name: 'Sprint Board',
		description: 'Sprint planning board',
		projectId: 'project-1'
	}),
	generateMockBoard({
		id: 'board-3',
		name: 'Backlog',
		description: 'Product backlog',
		projectId: 'project-2'
	})
];

export async function listBoards(projectId: string): Promise<ApiResponse<Board[]>> {
	const boards = mockBoards.filter((b) => b.projectId === projectId);
	return mockResponse(boards);
}
