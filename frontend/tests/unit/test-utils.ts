import { render } from '@testing-library/svelte';
import { writable } from 'svelte/store';
import type { Component } from 'svelte';

export function renderWithStores<T extends Component>(
	component: T,
	options: {
		props?: Record<string, unknown>;
		stores?: Record<string, unknown>;
	} = {}
) {
	const { props = {}, stores = {} } = options;

	const defaultStores = {
		page: writable({ url: { pathname: '/', searchParams: new URLSearchParams() } }),
		...stores
	};

	return render(component, {
		props,
		context: new Map(Object.entries(defaultStores))
	});
}

export function createMockUser(overrides = {}) {
	return {
		id: 'test-user-id',
		email: 'test@example.com',
		firstName: 'Test',
		lastName: 'User',
		role: 'member',
		communityId: 'test-community-id',
		createdAt: new Date().toISOString(),
		...overrides
	};
}

export function createMockProject(overrides = {}) {
	return {
		id: 'test-project-id',
		name: 'Test Project',
		description: 'Test description',
		status: 'active',
		communityId: 'test-community-id',
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString(),
		...overrides
	};
}

export function createMockTask(overrides = {}) {
	return {
		id: 'test-task-id',
		title: 'Test Task',
		description: 'Test description',
		status: 'todo',
		boardId: 'test-board-id',
		labels: [],
		position: 0,
		createdAt: new Date().toISOString(),
		updatedAt: new Date().toISOString(),
		...overrides
	};
}
