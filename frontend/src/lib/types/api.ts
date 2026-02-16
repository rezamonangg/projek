export interface User {
	id: string;
	email: string;
	firstName: string;
	lastName: string;
	role: 'admin' | 'member';
	communityId: string;
	avatarUrl?: string;
	createdAt: string;
}

export interface Community {
	id: string;
	name: string;
	slug: string;
	createdAt: string;
}

export interface Project {
	id: string;
	name: string;
	description: string;
	status: 'active' | 'archived';
	communityId: string;
	createdAt: string;
	updatedAt: string;
}

export interface Epic {
	id: string;
	name: string;
	description: string;
	color: string;
	projectId: string;
	createdAt: string;
}

export interface Board {
	id: string;
	name: string;
	description?: string;
	projectId: string;
	createdAt: string;
}

export interface Task {
	id: string;
	title: string;
	description: string;
	status: 'backlog' | 'todo' | 'inprogress' | 'done';
	assigneeId?: string;
	boardId: string;
	epicId?: string;
	labels: string[];
	position: number;
	createdAt: string;
	updatedAt: string;
}

export interface Label {
	id: string;
	name: string;
	color: string;
	communityId: string;
}

export interface WikiPage {
	id: string;
	title: string;
	slug: string;
	content: Record<string, unknown>;
	projectId: string;
	parentId?: string;
	createdAt: string;
	updatedAt: string;
}

export interface Attachment {
	id: string;
	filename: string;
	url: string;
	mimeType: string;
	size: number;
	taskId?: string;
	wikiPageId?: string;
	uploadedBy: string;
	createdAt: string;
}

export interface DashboardStats {
	totalMembers: number;
	totalProjects: number;
	totalTasks: number;
	tasksByStatus: {
		backlog: number;
		todo: number;
		inprogress: number;
		done: number;
	};
}

export interface CommunitySettings {
	id: string;
	name: string;
	slug: string;
	emailConfig?: {
		smtpHost: string;
		smtpPort: number;
		smtpUser: string;
		fromEmail: string;
	};
	storageConfig?: {
		provider: 'local' | 's3';
		bucket?: string;
	};
	metricsEnabled: boolean;
}
