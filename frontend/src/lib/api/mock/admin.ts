import type { ApiResponse } from '../types';
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

export async function getCommunitySettings(): Promise<ApiResponse<CommunitySettings>> {
	return mockResponse(mockSettings);
}
