import type { ApiResponse, IBoardsApi, CreateBoardInput } from '../types';
import type { Board } from '$lib/types/api';
import { mockResponse, mockError, generateId } from './utils';
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

async function listBoards(projectId: string): Promise<ApiResponse<Board[]>> {
	const boards = mockBoards.filter((b) => b.projectId === projectId);
	return mockResponse(boards);
}

async function getBoard(id: string): Promise<ApiResponse<Board>> {
	const board = mockBoards.find((b) => b.id === id);
	if (!board) {
		return mockError('Board not found', 404);
	}
	return mockResponse(board);
}

async function createBoard(input: CreateBoardInput): Promise<ApiResponse<Board>> {
	const newBoard = generateMockBoard({
		id: generateId(),
		name: input.name,
		description: input.description,
		projectId: input.projectId
	});
	mockBoards.push(newBoard);
	return mockResponse(newBoard);
}

export const mockBoardsApi: IBoardsApi = {
	listBoards,
	getBoard,
	createBoard
};
