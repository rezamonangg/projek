import type { ApiResponse, IAdminApi, UpdateCommunitySettingsInput } from '../types';
import type { CommunitySettings, DashboardStats } from '$lib/types/api';
import { mockResponse } from './utils';
import { generateMockCommunitySettings, generateMockDashboardStats } from './generators';

let mockSettings = generateMockCommunitySettings({
	id: 'community-1',
	name: 'Acme Corp',
	slug: 'acme-corp',
	emailConfig: {
		smtpHost: 'smtp.example.com',
		smtpPort: 587,
		smtpUser: 'notifications@acme.com',
		fromEmail: 'notifications@acme.com'
	},
	storageConfig: {
		provider: 'local'
	},
	metricsEnabled: true
});

async function getCommunitySettings(): Promise<ApiResponse<CommunitySettings>> {
	return mockResponse(mockSettings);
}

async function updateCommunitySettings(
	input: UpdateCommunitySettingsInput
): Promise<ApiResponse<CommunitySettings>> {
	mockSettings = {
		...mockSettings,
		...input
	};
	return mockResponse(mockSettings);
}

async function getDashboardStats(): Promise<ApiResponse<DashboardStats>> {
	return mockResponse(generateMockDashboardStats());
}

export const mockAdminApi: IAdminApi = {
	getCommunitySettings,
	updateCommunitySettings,
	getDashboardStats
};
