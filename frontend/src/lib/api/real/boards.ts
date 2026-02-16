import type { ApiResponse, IBoardsApi, CreateBoardInput } from '../types';
import type { Board } from '$lib/types/api';
import { httpClient } from './client';

interface BackendBoard {
	id: string;
	project_id: string;
	name: string;
	created_at: string;
	updated_at: string;
}

function toFrontendBoard(board: BackendBoard): Board {
	return {
		id: board.id,
		name: board.name,
		projectId: board.project_id,
		createdAt: board.created_at
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

async function listBoards(projectId: string): Promise<ApiResponse<Board[]>> {
	try {
		const boards = await httpClient.get<BackendBoard[]>(`/projects/${projectId}/boards`);
		return createResponse(boards.map(toFrontendBoard));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load boards';
		return createErrorResponse(message, 500);
	}
}

async function getBoard(id: string): Promise<ApiResponse<Board>> {
	try {
		const board = await httpClient.get<BackendBoard>(`/boards/${id}`);
		return createResponse(toFrontendBoard(board));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load board';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function createBoard(input: CreateBoardInput): Promise<ApiResponse<Board>> {
	try {
		const board = await httpClient.post<BackendBoard>(`/projects/${input.projectId}/boards`, {
			name: input.name
		});
		return createResponse(toFrontendBoard(board));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to create board';
		return createErrorResponse(message, 500);
	}
}

export const realBoardsApi: IBoardsApi = {
	listBoards,
	getBoard,
	createBoard
};
