# Phase 3: Projects & Kanban Core

Projects, epics, kanban boards, tasks, labels, and workflow management.

## Task 3.1: Project Domain Models

**Commits:**
- `feat: add project entity model`
- `feat: add epic entity model`
- `feat: add kanban board entity model`

**Models:**
```go
type Project struct {
    ID          uuid.UUID
    CommunityID uuid.UUID
    Name        string
    Description string
    Status      ProjectStatus // active, archived
    CreatedBy   uuid.UUID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Epic struct {
    ID          uuid.UUID
    ProjectID   uuid.UUID
    Name        string
    Description string
    Status      EpicStatus // active, completed
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Board struct {
    ID          uuid.UUID
    ProjectID   uuid.UUID
    Name        string
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## Task 3.2: Task Domain Models

**Commits:**
- `feat: add task entity model`
- `feat: add task status enum and workflow`
- `feat: add task assignment model`

**Models:**
```go
type Task struct {
    ID          uuid.UUID
    BoardID     uuid.UUID
    EpicID      *uuid.UUID
    Title       string
    Description string
    Status      TaskStatus // backlog, todo, inprogress, codereview, intest, needdeploy, done
    StartDate   *time.Time
    EndDate     *time.Time
    CreatedBy   uuid.UUID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type TaskStatus string
const (
    StatusBacklog    TaskStatus = "backlog"
    StatusTodo       TaskStatus = "todo"
    StatusInProgress TaskStatus = "inprogress"
    StatusCodeReview TaskStatus = "codereview"
    StatusInTest     TaskStatus = "intest"
    StatusNeedDeploy TaskStatus = "needdeploy"
    StatusDone       TaskStatus = "done"
)

type TaskAssignee struct {
    TaskID   uuid.UUID
    MemberID uuid.UUID
    AssignedAt time.Time
}
```

## Task 3.3: Label System

**Commits:**
- `feat: add label entity model`
- `feat: add task-label association model`
- `feat: add label colors and metadata`

**Models:**
```go
type Label struct {
    ID          uuid.UUID
    CommunityID uuid.UUID
    Name        string
    Color       string // hex color
    CreatedAt   time.Time
}

type TaskLabel struct {
    TaskID  uuid.UUID
    LabelID uuid.UUID
}
```

## Task 3.4: Project Repository

**Commits:**
- `feat: add project repository interface`
- `feat: implement project CRUD operations`
- `feat: add project filtering and pagination`

## Task 3.5: Project Service

**Commits:**
- `feat: add project service with business logic`
- `feat: add project permission checks`
- `feat: add project archiving functionality`

## Task 3.6: Project Handlers

**Commits:**
- `feat: add create project handler`
- `feat: add list projects handler`
- `feat: add get project handler`
- `feat: add update project handler`
- `feat: add delete/archive project handler`

## Task 3.7: Epic Repository

**Commits:**
- `feat: add epic repository interface`
- `feat: implement epic CRUD operations`
- `feat: add epic listing by project`

## Task 3.8: Epic Service & Handlers

**Commits:**
- `feat: add epic service`
- `feat: add epic handlers (CRUD)`

## Task 3.9: Board Repository

**Commits:**
- `feat: add board repository interface`
- `feat: implement board CRUD operations`
- `feat: add board listing by project`

## Task 3.10: Board Service & Handlers

**Commits:**
- `feat: add board service`
- `feat: add board handlers (CRUD)`

## Task 3.11: Task Repository

**Commits:**
- `feat: add task repository interface`
- `feat: implement task CRUD operations`
- `feat: add task filtering by status, assignee, epic`
- `feat: add task search functionality`

## Task 3.12: Task Service

**Commits:**
- `feat: add task service with business logic`
- `feat: add task workflow transitions`
- `feat: add task assignment logic`
- `feat: add task date validation`

## Task 3.13: Task Handlers

**Commits:**
- `feat: add create task handler`
- `feat: add list tasks handler with filters`
- `feat: add get task handler`
- `feat: add update task handler`
- `feat: add move task status handler`
- `feat: add delete task handler`

## Task 3.14: Label Repository

**Commits:**
- `feat: add label repository interface`
- `feat: implement label CRUD operations`

## Task 3.15: Label Service & Handlers

**Commits:**
- `feat: add label service`
- `feat: add label handlers`
- `feat: add assign/remove label from task`

## Task 3.16: Project Permissions

**Commits:**
- `feat: add project permission model`
- `feat: add project member repository`
- `feat: add permission middleware`

**Permission model:**
- Admin: full access to all projects
- Member: access only to invited projects

## Task 3.17: Database Migrations for Projects

**Commits:**
- `task: add projects table migration`
- `task: add epics table migration`
- `task: add boards table migration`
- `task: add tasks table migration`
- `task: add labels and task_labels tables`
- `task: add project_members table for permissions`
- `task: add indexes for performance`

---

## Total Commits: 42
**Estimated Time:** 4-5 days
