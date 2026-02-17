# Backend Phase 3: Projects & Kanban

Projects, epics, kanban boards, tasks, labels, and workflow management.

## Task 3.1: Project Domain Models

**Commits:**
- `feat: add project entity model`
- `feat: add epic entity model`
- `feat: add kanban board entity model`

**Package:** `internal/project/model.go`

## Task 3.2: Task Domain Models

**Commits:**
- `feat: add task entity model`
- `feat: add task status enum and workflow`
- `feat: add task assignment model`

**Package:** `internal/task/model.go`

**TaskStatus values:**
- backlog, todo, inprogress, codereview, intest, needdeploy, done

## Task 3.3: Label System

**Commits:**
- `feat: add label entity model`
- `feat: add task-label association model`

**Package:** `internal/task/label.go`

## Task 3.4: Project Repository

**Commits:**
- `feat: add project repository interface`
- `feat: implement project CRUD operations`
- `feat: add project filtering and pagination`

**Package:** `internal/project/repository.go`

## Task 3.5: Project Service

**Commits:**
- `feat: add project service with business logic`
- `feat: add project permission checks`
- `feat: add project archiving functionality`

**Package:** `internal/project/service.go`

## Task 3.6: Project Handlers

**Commits:**
- `feat: add create project handler`
- `feat: add list projects handler`
- `feat: add get project handler`
- `feat: add update project handler`
- `feat: add delete/archive project handler`

**Package:** `internal/project/handler.go`

## Task 3.7: Epic Repository

**Commits:**
- `feat: add epic repository interface`
- `feat: implement epic CRUD operations`

**Package:** `internal/project/epic_repository.go`

## Task 3.8: Epic Service & Handlers

**Commits:**
- `feat: add epic service`
- `feat: add epic handlers (CRUD)`

## Task 3.9: Board Repository

**Commits:**
- `feat: add board repository interface`
- `feat: implement board CRUD operations`

**Package:** `internal/project/board_repository.go`

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

**Package:** `internal/task/repository.go`

## Task 3.12: Task Service

**Commits:**
- `feat: add task service with business logic`
- `feat: add task workflow transitions`
- `feat: add task assignment logic`

**Package:** `internal/task/service.go`

## Task 3.13: Task Handlers

**Commits:**
- `feat: add create task handler`
- `feat: add list tasks handler with filters`
- `feat: add get task handler`
- `feat: add update task handler`
- `feat: add move task status handler`
- `feat: add delete task handler`

**Package:** `internal/task/handler.go`

## Task 3.14: Label Repository

**Commits:**
- `feat: add label repository interface`
- `feat: implement label CRUD operations`

## Task 3.15: Label Service & Handlers

**Commits:**
- `feat: add label service`
- `feat: add label handlers`

## Task 3.16: Project Permissions

**Commits:**
- `feat: add project permission model`
- `feat: add project member repository`
- `feat: add permission middleware`

**Package:** `internal/project/permission.go`

## Task 3.17: Database Migrations

**Commits:**
- `task: add projects table migration`
- `task: add epics table migration`
- `task: add boards table migration`
- `task: add tasks table migration`
- `task: add labels and task_labels tables`
- `task: add project_members table`
- `task: add indexes for performance`

---

## Total Backend Commits (Phase 3): 42
