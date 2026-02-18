# Frontend Testing: Projects Page

## Overview

E2E tests for projects list page with real backend API integration for CRUD operations.

## Task 1: Projects List Loading

**Test Scenarios:**
- TC-01: Load projects list
  - Login via API
  - Create test projects via API
  - Navigate to `/projects`
  - Verify loading state (skeleton cards) shown initially
  - Wait for projects to load
  - Verify all projects displayed as cards
  - Verify project details: name, description, status, created date
  - Verify status badges (active/archived)

- TC-02: Empty projects list
  - Login via API with fresh user
  - Navigate to `/projects`
  - Verify empty state message
  - Verify "Create your first project" button

- TC-03: Projects list API error
  - Login via API
  - Simulate API failure
  - Navigate to `/projects`
  - Verify error handling
  - Verify retry mechanism if exists

**API Coverage:**
- `GET /projects`
- Loading states (SkeletonCard component)
- Error states

## Task 2: Create Project

**Test Scenarios:**
- TC-01: Successfully create project
  - Login via API
  - Navigate to `/projects`
  - Click "Create Project" button
  - Fill in project name
  - Fill in project description
  - Click "Create"
  - Verify modal closes
  - Verify new project appears in list
  - Verify API called correctly

- TC-02: Create project validation - empty name
  - Login via API
  - Open create project modal
  - Submit without name
  - Verify validation error "Project name is required"
  - Verify API not called

- TC-03: Create project - cancel
  - Login via API
  - Open create project modal
  - Click Cancel
  - Verify modal closes
  - Verify no API call made

- TC-04: Create project - long name/description
  - Test with maximum length inputs
  - Verify handled correctly

**API Coverage:**
- `POST /projects`
- Form validation
- Modal interactions

## Task 3: Project Card Interactions

**Test Scenarios:**
- TC-01: Click project card
  - Login via API
  - Navigate to `/projects`
  - Click on a project card
  - Verify navigation to `/projects/{id}`
  - Verify project details page loads

- TC-02: Project card hover states
  - Verify hover effects on cards
  - Verify cursor changes to pointer

- TC-03: Status badge display
  - Verify active projects have green badge
  - Verify archived projects have gray badge

**API Coverage:**
- Project data structure
- Navigation routing

## Task 4: Page Layout

**Test Scenarios:**
- TC-01: Projects page layout
  - Login via API
  - Navigate to `/projects`
  - Verify breadcrumb shows "Projects"
  - Verify page title correct
  - Verify "Create Project" button visible
  - Verify responsive grid layout

**API Coverage:**
- Page authentication check

## File Changes

```
frontend/tests/e2e/
└── projects.spec.ts (UPDATE - Replace existing tests)
```

## Test Setup Requirements

```typescript
// Helper to create test projects
async function createTestProject(page, data: { name: string, description: string }) {
  // Create via API
}

// Helper to cleanup test projects
async function cleanupTestProjects(page) {
  // Delete created projects via API
}
```

## Expected Test Count

- List loading: 3
- Create project: 4
- Card interactions: 3
- Page layout: 1
- **Total: 11 tests**

## Notes

- Current `projects.spec.ts` has misleading tests that only check redirects
- Replace all existing tests with these comprehensive tests
- Ensure proper cleanup of test data after each test
