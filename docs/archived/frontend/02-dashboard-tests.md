# Frontend Testing: Dashboard Page

## Overview

E2E tests for dashboard page with real backend API integration for stats and projects.

## Task 1: Dashboard Stats Loading

**Test Scenarios:**
- TC-01: Load dashboard stats successfully
  - Login via API
  - Navigate to `/dashboard`
  - Verify loading state shown initially
  - Wait for stats to load
  - Verify all stat cards displayed:
    - Total Members count
    - Total Projects count
    - Total Tasks count
    - Completed Tasks count
  - Verify stats match API response

- TC-02: Dashboard stats API error handling
  - Login via API
  - Mock API to return 500 error (or use network offline)
  - Navigate to `/dashboard`
  - Verify error state or empty stats
  - Verify no crash

**API Coverage:**
- `GET /admin/dashboard/stats`
- Loading state handling
- Error state handling

## Task 2: Recent Projects Section

**Test Scenarios:**
- TC-01: Display recent projects
  - Login via API
  - Create test projects via API
  - Navigate to `/dashboard`
  - Verify recent projects section shown
  - Verify up to 5 projects displayed
  - Verify project names, descriptions, and status
  - Verify project links work

- TC-02: Empty projects state
  - Login via API with fresh user
  - Navigate to `/dashboard`
  - Verify "No projects yet" message
  - Verify "Create your first project" prompt shown

**API Coverage:**
- `GET /projects` (limited to 5 items)
- Project data structure validation

## Task 3: Task Distribution Chart

**Test Scenarios:**
- TC-01: Task distribution bars
  - Login via API
  - Navigate to `/dashboard`
  - Verify task distribution section
  - Verify bars for:
    - Backlog
    - To Do
    - In Progress
    - Done
  - Verify percentage calculations

- TC-02: Empty tasks state
  - Login via API
  - Navigate to `/dashboard`
  - Verify charts handle zero tasks gracefully

**API Coverage:**
- `GET /admin/dashboard/stats` (tasksByStatus)
- Chart rendering with real data

## Task 4: Page Layout and Navigation

**Test Scenarios:**
- TC-01: Dashboard layout elements
  - Login via API
  - Navigate to `/dashboard`
  - Verify breadcrumb shows "Dashboard"
  - Verify page title correct
  - Verify welcome message shown

- TC-02: Sidebar navigation
  - Login via API
  - Navigate to `/dashboard`
  - Verify sidebar visible
  - Verify navigation links work:
    - Dashboard
    - Projects
    - Members
    - Settings

**API Coverage:**
- Authentication check on page load
- Navigation between protected routes

## File Changes

```
frontend/tests/e2e/
└── dashboard.spec.ts (CREATE NEW)
```

## Test Setup Requirements

```typescript
// Create test data via API
async function createTestProjects(page, count: number) {
  // Create projects via API for testing
}

async function getDashboardStats() {
  // Helper to get expected stats from API
}
```

## Expected Test Count

- Stats loading: 2
- Projects section: 2
- Task distribution: 2
- Layout/navigation: 2
- **Total: 8 new tests**
