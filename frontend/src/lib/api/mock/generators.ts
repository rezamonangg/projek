import type {
	User,
	Community,
	Project,
	Epic,
	Board,
	Task,
	Label,
	WikiPage,
	Attachment,
	DashboardStats,
	CommunitySettings
} from '../types/api';
import { generateId } from './utils';

const now = () => new Date().toISOString();

export function generateMockUser(overrides: Partial<User> = {}): User {
	return {
		id: generateId(),
		email: 'user@example.com',
		firstName: 'John',
		lastName: 'Doe',
		role: 'member',
		communityId: generateId(),
		createdAt: now(),
		...overrides
	};
}

export function generateMockCommunity(overrides: Partial<Community> = {}): Community {
	return {
		id: generateId(),
		name: 'Acme Corp',
		slug: 'acme-corp',
		createdAt: now(),
		...overrides
	};
}

export function generateMockProject(overrides: Partial<Project> = {}): Project {
	return {
		id: generateId(),
		name: 'Website Redesign',
		description: 'Redesign the company website',
		status: 'active',
		communityId: generateId(),
		createdAt: now(),
		updatedAt: now(),
		...overrides
	};
}

export function generateMockEpic(overrides: Partial<Epic> = {}): Epic {
	return {
		id: generateId(),
		name: 'Homepage',
		description: 'Homepage redesign epic',
		color: '#3B82F6',
		projectId: generateId(),
		createdAt: now(),
		...overrides
	};
}

export function generateMockBoard(overrides: Partial<Board> = {}): Board {
	return {
		id: generateId(),
		name: 'Main Board',
		description: 'Primary kanban board',
		projectId: generateId(),
		createdAt: now(),
		...overrides
	};
}

export function generateMockTask(overrides: Partial<Task> = {}): Task {
	return {
		id: generateId(),
		title: 'Design homepage',
		description: 'Create Figma designs for the homepage',
		status: 'todo',
		boardId: generateId(),
		labels: [],
		position: 0,
		createdAt: now(),
		updatedAt: now(),
		...overrides
	};
}

export function generateMockLabel(overrides: Partial<Label> = {}): Label {
	return {
		id: generateId(),
		name: 'Bug',
		color: '#EF4444',
		communityId: generateId(),
		...overrides
	};
}

export function generateMockWikiPage(overrides: Partial<WikiPage> = {}): WikiPage {
	return {
		id: generateId(),
		title: 'Getting Started',
		slug: 'getting-started',
		content: { type: 'doc', content: [] },
		projectId: generateId(),
		createdAt: now(),
		updatedAt: now(),
		...overrides
	};
}

export function generateMockAttachment(overrides: Partial<Attachment> = {}): Attachment {
	return {
		id: generateId(),
		filename: 'document.pdf',
		url: '/uploads/document.pdf',
		mimeType: 'application/pdf',
		size: 1024,
		uploadedBy: generateId(),
		createdAt: now(),
		...overrides
	};
}

export function generateMockDashboardStats(): DashboardStats {
	return {
		totalMembers: 10,
		totalProjects: 5,
		totalTasks: 150,
		tasksByStatus: {
			backlog: 20,
			todo: 30,
			inprogress: 40,
			done: 60
		}
	};
}

export function generateMockCommunitySettings(
	overrides: Partial<CommunitySettings> = {}
): CommunitySettings {
	return {
		id: generateId(),
		name: 'Acme Corp',
		slug: 'acme-corp',
		metricsEnabled: true,
		...overrides
	};
}
