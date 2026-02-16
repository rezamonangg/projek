import type { ApiResponse, ITasksApi, CreateTaskInput, UpdateTaskInput, MoveTaskInput } from '../types';
import type { Task } from '$lib/types/api';
import { httpClient } from './client';

type BackendStatus = 'backlog' | 'todo' | 'inprogress' | 'codereview' | 'intest' | 'needdeploy' | 'done';

interface BackendTask {
	id: string;
	board_id: string;
	epic_id?: string;
	title: string;
	description: string;
	status: BackendStatus;
	position: number;
	story_points?: number;
	due_date?: string;
	assignee_id?: string;
	reporter_id: string;
	created_at: string;
	updated_at: string;
	labels?: string[];
}

function toFrontendStatus(status: BackendStatus): Task['status'] {
	switch (status) {
		case 'done':
			return 'done';
		case 'inprogress':
		case 'codereview':
		case 'intest':
		case 'needdeploy':
			return 'inprogress';
		case 'todo':
			return 'todo';
		default:
			return 'backlog';
	}
}

function toBackendStatus(status: Task['status']): BackendStatus {
	return status as BackendStatus;
}

function toFrontendTask(task: BackendTask): Task {
	return {
		id: task.id,
		title: task.title,
		description: task.description,
		status: toFrontendStatus(task.status),
		boardId: task.board_id,
		epicId: task.epic_id,
		assigneeId: task.assignee_id,
		labels: task.labels || [],
		position: task.position,
		createdAt: task.created_at,
		updatedAt: task.updated_at
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

async function listTasks(boardId: string): Promise<ApiResponse<Task[]>> {
	try {
		const tasks = await httpClient.get<BackendTask[]>(`/boards/${boardId}/tasks`);
		return createResponse(tasks.map(toFrontendTask));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load tasks';
		return createErrorResponse(message, 500);
	}
}

async function createTask(input: CreateTaskInput): Promise<ApiResponse<Task>> {
	try {
		const task = await httpClient.post<BackendTask>(`/boards/${input.boardId}/tasks`, {
			title: input.title,
			description: input.description,
			epic_id: input.epicId,
			assignee_id: input.assigneeId,
			status: 'backlog',
			reporter_id: input.assigneeId || '00000000-0000-0000-0000-000000000000'
		});
		return createResponse(toFrontendTask(task));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to create task';
		return createErrorResponse(message, 500);
	}
}

async function updateTask(id: string, input: UpdateTaskInput): Promise<ApiResponse<Task>> {
	try {
		const body: Record<string, unknown> = {};
		if (input.title) body.title = input.title;
		if (input.description) body.description = input.description;
		if (input.status) body.status = toBackendStatus(input.status);
		if (input.epicId) body.epic_id = input.epicId;
		if (input.assigneeId) body.assignee_id = input.assigneeId;

		const task = await httpClient.put<BackendTask>(`/tasks/${id}`, body);
		return createResponse(toFrontendTask(task));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to update task';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function moveTask(id: string, input: MoveTaskInput): Promise<ApiResponse<Task>> {
	try {
		const task = await httpClient.patch<BackendTask>(`/tasks/${id}/move`, {
			status: toBackendStatus(input.status),
			position: input.position
		});
		return createResponse(toFrontendTask(task));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to move task';
		const status = message.includes('not found') ? 404 : 500;
		return createErrorResponse(message, status);
	}
}

async function deleteTask(id: string): Promise<ApiResponse<void>> {
	try {
		await httpClient.delete(`/tasks/${id}`);
		return createResponse(undefined);
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to delete task';
		return createErrorResponse(message, 500);
	}
}

export const realTasksApi: ITasksApi = {
	listTasks,
	createTask,
	updateTask,
	moveTask,
	deleteTask
};
