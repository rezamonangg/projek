import type { ApiResponse, PaginatedResponse } from '../types';
import type { Project } from '$lib/types/api';
import { mockResponse } from './utils';
import { generateMockProject } from './generators';

const mockProjects: Project[] = [
	generateMockProject({
		id: 'project-1',
		name: 'Website Redesign',
		description: 'Redesign the company website with modern UI/UX',
		status: 'active',
		communityId: 'community-1'
	}),
	generateMockProject({
		id: 'project-2',
		name: 'Mobile App',
		description: 'Native mobile application for iOS and Android',
		status: 'active',
		communityId: 'community-1'
	}),
	generateMockProject({
		id: 'project-3',
		name: 'API Platform',
		description: 'RESTful API platform for third-party integrations',
		status: 'active',
		communityId: 'community-1'
	}),
	generateMockProject({
		id: 'project-4',
		name: 'Legacy System',
		description: 'Maintenance of the legacy system',
		status: 'archived',
		communityId: 'community-1'
	})
];

export async function listProjects(): Promise<ApiResponse<PaginatedResponse<Project>>> {
	return mockResponse({
		items: mockProjects,
		total: mockProjects.length,
		page: 1,
		pageSize: 10,
		totalPages: 1
	});
}

export async function getProject(id: string): Promise<ApiResponse<Project>> {
	const project = mockProjects.find((p) => p.id === id);
	if (!project) {
		return mockResponse(null as unknown as Project);
	}
	return mockResponse(project);
}
