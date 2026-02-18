# Frontend Testing: Kanban Board

## Overview

E2E tests for kanban board page with real backend API integration for boards and tasks.

## Task 1: Board Loading

**Test Scenarios:**
- TC-01: Load kanban board successfully
  - Login via API
  - Create test project and board via API
  - Create test tasks via API
  - Navigate to `/projects/{id}/boards/{boardId}`
  - Verify board loads
  - Verify columns displayed: Backlog, To Do, In Progress, Done
  - Verify tasks displayed in correct columns
  - Verify task cards show title and info

- TC-02: Empty board
  - Login via API
  - Create empty board via API
  - Navigate to board
  - Verify empty columns shown
  - Verify "No tasks" message or placeholder

- TC-03: Board loading error
  - Login via API
  - Navigate to non-existent board
  - Verify error state
  - Verify 404 handling or error message

**API Coverage:**
- `GET /boards/{id}`
- `GET /boards/{boardId}/tasks`
- Column rendering
- Task card rendering

## Task 2: Task Cards Display

**Test Scenarios:**
- TC-01: Task card details
  - Login via API
  - Create task with labels, assignee via API
  - Navigate to board
  - Verify task card shows:
    - Title
    - Labels (with colors)
    - Assignee avatar
    - Due date (if set)

- TC-02: Multiple tasks in column
  - Login via API
  - Create multiple tasks in same column
  - Verify all tasks displayed
  - Verify scrollable if many tasks

- TC-03: Task with long title
  - Create task with very long title
  - Verify truncation or wrapping

**API Coverage:**
- Task data structure
- Label data
- Member/assignee data

## Task 3: Task Drag and Drop (If Implemented)

**Test Scenarios:**
- TC-01: Move task between columns
  - Login via API
  - Create task in "To Do" column
  - Drag task to "In Progress"
  - Verify visual feedback during drag
  - Verify task stays in new column after drop
  - Verify API called to update task status

- TC-02: Reorder tasks in same column
  - Login via API
  - Create multiple tasks
  - Reorder by dragging
  - Verify new order persisted

**API Coverage:**
- `PATCH /tasks/{id}/move`
- Drag and drop events

## Task 4: Create Task from Board

**Test Scenarios:**
- TC-01: Create new task via board
  - Login via API
  - Navigate to board
  - Click "Add Task" in column
  - Fill task form
  - Submit
  - Verify task appears in column
  - Verify API called correctly

- TC-02: Quick add task
  - Test quick add functionality if exists
  - Verify inline editing

**API Coverage:**
- `POST /boards/{boardId}/tasks`
- Task creation form

## Task 5: Board Navigation

**Test Scenarios:**
- TC-01: Breadcrumb navigation
  - Login via API
  - Navigate to board
  - Verify breadcrumb shows: Projects > {Project Name} > Boards > {Board Name}
  - Verify breadcrumb links work

- TC-02: Back to project
  - Click back link or breadcrumb
  - Verify navigation to project page

**API Coverage:**
- Project data for breadcrumb
- Navigation routing

## File Changes

```
frontend/tests/e2e/
└── kanban.spec.ts (UPDATE - Replace existing tests)
```

## Test Setup Requirements

```typescript
// Helper to create test board with tasks
async function createTestBoardWithTasks(page, projectId: string) {
  // Create board via API
  // Create tasks via API
  // Return board ID and task IDs
}

// Helper to cleanup
async function cleanupTestBoard(page, boardId: string) {
  // Delete board and tasks via API
}
```

## Expected Test Count

- Board loading: 3
- Task cards: 3
- Drag and drop: 2
- Create task: 2
- Navigation: 2
- **Total: 12 tests**

## Notes

- Current `kanban.spec.ts` only tests redirect when not authenticated
- Need actual board and task API integration
- Drag and drop tests may be complex, mark as optional if needed
