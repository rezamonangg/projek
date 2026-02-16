# Integration Phase 4: Task & Label APIs

Replace mock task APIs with real backend endpoints.

## Task 4.1: Real Task API Implementation

**Commits:**
- `feat: add real task api client`
- `feat: integrate task endpoints`

**Endpoints:**
- GET /api/boards/:id/tasks
- POST /api/tasks
- GET /api/tasks/:id
- PUT /api/tasks/:id
- DELETE /api/tasks/:id
- PATCH /api/tasks/:id/status

## Task 4.2: Real Label API Implementation

**Commits:**
- `feat: add real label api client`
- `feat: integrate label endpoints`

**Endpoints:**
- GET /api/communities/:id/labels
- POST /api/communities/:id/labels
- POST /api/tasks/:id/labels
- DELETE /api/tasks/:id/labels/:labelId

## Task 4.3: Update Drag and Drop

**Commits:**
- `feat: update kanban drag-drop with real api`
- `feat: add optimistic updates`
- `feat: add error rollback`

## Task 4.4: Switch to Real Tasks

**Commits:**
- `feat: switch task api to real implementation`
- `feat: switch label api to real implementation`
- `test: verify kanban board with real data`
- `test: verify task CRUD operations`

## Task 4.5: Remove Mock Tasks

**Commits:**
- `chore: remove mock task and label apis`

---

## Total Integration Commits (Phase 4): 9
