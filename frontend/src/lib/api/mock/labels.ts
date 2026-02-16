import type { ApiResponse } from '../types';
import type { Label, Task } from '$lib/types/api';
import { mockResponse, generateId } from './utils';
import { generateMockLabel } from './generators';
import { mockTasks } from './tasks';

const mockLabels: Label[] = [
	generateMockLabel({
		id: 'label-1',
		name: 'Bug',
		color: '#EF4444',
		communityId: 'community-1'
	}),
	generateMockLabel({
		id: 'label-2',
		name: 'Feature',
		color: '#3B82F6',
		communityId: 'community-1'
	}),
	generateMockLabel({
		id: 'label-3',
		name: 'Enhancement',
		color: '#10B981',
		communityId: 'community-1'
	}),
	generateMockLabel({
		id: 'label-4',
		name: 'Documentation',
		color: '#F59E0B',
		communityId: 'community-1'
	}),
	generateMockLabel({
		id: 'label-5',
		name: 'Design',
		color: '#8B5CF6',
		communityId: 'community-1'
	})
];

export async function listLabels(communityId: string): Promise<ApiResponse<Label[]>> {
	const labels = mockLabels.filter((l) => l.communityId === communityId);
	return mockResponse(labels);
}

interface CreateLabelInput {
	name: string;
	color: string;
	communityId: string;
}

export async function createLabel(input: CreateLabelInput): Promise<ApiResponse<Label>> {
	const newLabel = generateMockLabel({
		id: generateId(),
		name: input.name,
		color: input.color,
		communityId: input.communityId
	});
	mockLabels.push(newLabel);
	return mockResponse(newLabel);
}

export async function assignLabelToTask(
	taskId: string,
	labelId: string
): Promise<ApiResponse<Task>> {
	const taskIndex = mockTasks.findIndex((t) => t.id === taskId);
	if (taskIndex === -1) {
		return mockResponse(null as unknown as Task);
	}

	const label = mockLabels.find((l) => l.id === labelId);
	if (!label) {
		return mockResponse(null as unknown as Task);
	}

	if (!mockTasks[taskIndex].labels.includes(labelId)) {
		mockTasks[taskIndex].labels.push(labelId);
	}

	return mockResponse(mockTasks[taskIndex]);
}
