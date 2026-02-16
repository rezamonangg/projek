import type { ApiResponse, IAdminApi, UpdateCommunitySettingsInput } from '../types';
import type { CommunitySettings, DashboardStats } from '$lib/types/api';
import { httpClient } from './client';

interface BackendCommunitySettings {
	id: string;
	community_id: string;
	allow_member_registration: boolean;
	require_email_verification: boolean;
	created_at: string;
	updated_at: string;
}

interface BackendDashboardStats {
	total_members: number;
	total_projects: number;
	total_tasks: number;
	active_tasks: number;
	completed_tasks: number;
}

function toFrontendSettings(settings: BackendCommunitySettings): CommunitySettings {
	return {
		id: settings.id,
		name: '',
		slug: '',
		metricsEnabled: true
	};
}

function toFrontendStats(stats: BackendDashboardStats): DashboardStats {
	return {
		totalMembers: stats.total_members,
		totalProjects: stats.total_projects,
		totalTasks: stats.total_tasks,
		tasksByStatus: {
			backlog: 0,
			todo: 0,
			inprogress: stats.active_tasks,
			done: stats.completed_tasks
		}
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

async function getCommunitySettings(): Promise<ApiResponse<CommunitySettings>> {
	try {
		const settings = await httpClient.get<BackendCommunitySettings>('/admin/settings');
		return createResponse(toFrontendSettings(settings));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load settings';
		return createErrorResponse(message, 500);
	}
}

async function updateCommunitySettings(_input: UpdateCommunitySettingsInput): Promise<ApiResponse<CommunitySettings>> {
	try {
		const settings = await httpClient.put<BackendCommunitySettings>('/admin/settings', {});
		return createResponse(toFrontendSettings(settings));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to update settings';
		return createErrorResponse(message, 500);
	}
}

async function getDashboardStats(): Promise<ApiResponse<DashboardStats>> {
	try {
		const stats = await httpClient.get<BackendDashboardStats>('/admin/dashboard');
		return createResponse(toFrontendStats(stats));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load dashboard stats';
		return createErrorResponse(message, 500);
	}
}

export const realAdminApi: IAdminApi = {
	getCommunitySettings,
	updateCommunitySettings,
	getDashboardStats
};
