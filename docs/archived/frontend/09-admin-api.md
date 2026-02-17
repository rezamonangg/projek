# Frontend Implementation Plan: Admin API

**Priority:** LOW
**Estimated Time:** 30 minutes
**Backend Reference:** `docs/plan/backend/09-admin-module.md`

## Current State

### What Exists
- `frontend/src/lib/api/real/admin.ts` - Implementation exists
- `frontend/src/lib/api/types.ts` - IAdminApi interface

### What to Verify
- Settings endpoint structure
- Dashboard stats structure
- Admin-only access handling (403)

## Backend API Contract

| Method | Endpoint | Request Body | Response |
|--------|----------|--------------|----------|
| GET | `/admin/settings` | - | `CommunitySettings` |
| PUT | `/admin/settings` | `{allow_member_registration?, require_email_verification?}` | `CommunitySettings` |
| GET | `/admin/dashboard` | - | `DashboardStats` |

## Implementation Details

### Backend Response Format

```typescript
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
```

### Frontend Type Mapping

```typescript
// From $lib/types/api
interface CommunitySettings {
	id: string;
	communityId: string;
	allowMemberRegistration: boolean;
	requireEmailVerification: boolean;
	createdAt: string;
	updatedAt: string;
	emailConfig?: {
		smtpHost: string;
		smtpPort: number;
		fromEmail: string;
	};
	storageConfig?: {
		provider: 'local' | 's3';
		maxFileSize: number;
	};
	metricsEnabled?: boolean;
}

interface DashboardStats {
	totalMembers: number;
	totalProjects: number;
	totalTasks: number;
	activeTasks: number;
	completedTasks: number;
}
```

### Transformer Function

```typescript
function toFrontendSettings(settings: BackendCommunitySettings): CommunitySettings {
	return {
		id: settings.id,
		communityId: settings.community_id,
		allowMemberRegistration: settings.allow_member_registration,
		requireEmailVerification: settings.require_email_verification,
		createdAt: settings.created_at,
		updatedAt: settings.updated_at
	};
}

function toFrontendStats(stats: BackendDashboardStats): DashboardStats {
	return {
		totalMembers: stats.total_members,
		totalProjects: stats.total_projects,
		totalTasks: stats.total_tasks,
		activeTasks: stats.active_tasks,
		completedTasks: stats.completed_tasks
	};
}
```

### API Implementation

```typescript
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
		communityId: settings.community_id,
		allowMemberRegistration: settings.allow_member_registration,
		requireEmailVerification: settings.require_email_verification,
		createdAt: settings.created_at,
		updatedAt: settings.updated_at
	};
}

function toFrontendStats(stats: BackendDashboardStats): DashboardStats {
	return {
		totalMembers: stats.total_members,
		totalProjects: stats.total_projects,
		totalTasks: stats.total_tasks,
		activeTasks: stats.active_tasks,
		completedTasks: stats.completed_tasks
	};
}

function createResponse<T>(data: T): ApiResponse<T> {
	return { data, error: null, status: 200 };
}

function createErrorResponse<T>(error: string, status: number): ApiResponse<T> {
	return { data: null, error, status };
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

async function updateCommunitySettings(input: UpdateCommunitySettingsInput): Promise<ApiResponse<CommunitySettings>> {
	try {
		const body: Record<string, unknown> = {};
		if (input.emailConfig) body.email_config = input.emailConfig;
		if (input.storageConfig) body.storage_config = input.storageConfig;
		if (input.metricsEnabled !== undefined) body.metrics_enabled = input.metricsEnabled;

		const settings = await httpClient.put<BackendCommunitySettings>('/admin/settings', body);
		return createResponse(toFrontendSettings(settings));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to update settings';
		const status = message.includes('forbidden') ? 403 : 500;
		return createErrorResponse(message, status);
	}
}

async function getDashboardStats(): Promise<ApiResponse<DashboardStats>> {
	try {
		const stats = await httpClient.get<BackendDashboardStats>('/admin/dashboard');
		return createResponse(toFrontendStats(stats));
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Failed to load dashboard';
		return createErrorResponse(message, 500);
	}
}

export const realAdminApi: IAdminApi = {
	getCommunitySettings,
	updateCommunitySettings,
	getDashboardStats
};
```

## Error Handling

| Status Code | Scenario | Action |
|-------------|----------|--------|
| 400 | Validation error | Show validation message |
| 401 | Unauthorized | Redirect to login |
| 403 | Forbidden (non-admin) | Show permission error |
| 500 | Server error | Show generic error |

## Files to Verify

| File | Status |
|------|--------|
| `frontend/src/lib/api/real/admin.ts` | EXISTS - Verify alignment |
| `frontend/src/lib/api/types.ts` | EXISTS - Verify IAdminApi |
