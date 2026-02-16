import type { ApiResponse } from '../types';
import type { Task } from '$lib/types/api';
import { mockResponse, generateId } from './utils';
import { generateMockTask } from './generators';

export const mockTasks: Task[] = [
	generateMockTask({
		id: 'task-1',
		title: 'Design homepage hero section',
		description: 'Create the hero section design with carousel',
		status: 'inprogress',
		boardId: 'board-1',
		epicId: 'epic-1',
		assigneeId: 'user-1',
		labels: ['design'],
		position: 0
	}),
	generateMockTask({
		id: 'task-2',
		title: 'Implement navigation component',
		description: 'Build responsive navigation with dropdown menus',
		status: 'todo',
		boardId: 'board-1',
		epicId: 'epic-1',
		assigneeId: 'user-2',
		labels: ['frontend'],
		position: 0
	}),
	generateMockTask({
		id: 'task-3',
		title: 'Set up authentication flow',
		description: 'Implement login, logout, and registration',
		status: 'done',
		boardId: 'board-1',
		epicId: 'epic-3',
		assigneeId: 'user-1',
		labels: ['backend', 'security'],
		position: 0
	}),
	generateMockTask({
		id: 'task-4',
		title: 'Create API documentation',
		description: 'Document all REST endpoints',
		status: 'backlog',
		boardId: 'board-1',
		labels: ['documentation'],
		position: 0
	})
];

export async function listTasks(boardId: string): Promise<ApiResponse<Task[]>> {
	const tasks = mockTasks.filter((t) => t.boardId === boardId);
	return mockResponse(tasks);
}

interface CreateTaskInput {
	title: string;
	description: string;
	boardId: string;
	epicId?: string;
	assigneeId?: string;
	labels?: string[];
}

export async function createTask(input: CreateTaskInput): Promise<ApiResponse<Task>> {
	const newTask = generateMockTask({
		id: generateId(),
		title: input.title,
		description: input.description,
		boardId: input.boardId,
		epicId: input.epicId,
		assigneeId: input.assigneeId,
		labels: input.labels || [],
		status: 'backlog',
		position: mockTasks.filter((t) => t.boardId === input.boardId).length
	});
	mockTasks.push(newTask);
	return mockResponse(newTask);
}

interface UpdateTaskInput {
	title?: string;
	description?: string;
	status?: 'backlog' | 'todo' | 'inprogress' | 'done';
	assigneeId?: string;
	epicId?: string;
	labels?: string[];
}

export async function updateTask(id: string, input: UpdateTaskInput): Promise<ApiResponse<Task>> {
	const index = mockTasks.findIndex((t) => t.id === id);
	if (index === -1) {
		return mockResponse(null as unknown as Task);
	}

	mockTasks[index] = {
		...mockTasks[index],
		...input,
		updatedAt: new Date().toISOString()
	};

	return mockResponse(mockTasks[index]);
}

interface MoveTaskInput {
	status: 'backlog' | 'todo' | 'inprogress' | 'done';
	position: number;
}

export async function moveTask(id: string, input: MoveTaskInput): Promise<ApiResponse<Task>> {
	const index = mockTasks.findIndex((t) => t.id === id);
	if (index === -1) {
		return mockResponse(null as unknown as Task);
	}

	mockTasks[index] = {
		...mockTasks[index],
		status: input.status,
		position: input.position,
		updatedAt: new Date().toISOString()
	};

	return mockResponse(mockTasks[index]);
}

export async function deleteTask(id: string): Promise<ApiResponse<void>> {
	const index = mockTasks.findIndex((t) => t.id === id);
	if (index !== -1) {
		mockTasks.splice(index, 1);
	}
	return mockResponse(undefined);
}
