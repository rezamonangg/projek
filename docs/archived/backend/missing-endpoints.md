# Backend: Missing API Endpoints

Endpoints needed for frontend dashboard functionality.

## Task: Projects API

**Commits:**
- `feat: add projects repository interface and implementation`
- `feat: add projects service`
- `feat: add projects handlers (list, get, create, update, delete)`
- `feat: mount projects routes in router`

**Endpoints:**
- `GET /projects` - List all projects
- `GET /projects/:id` - Get project by ID
- `POST /projects` - Create project
- `PUT /projects/:id` - Update project
- `DELETE /projects/:id` - Delete project

**Package:** `internal/project/repository.go`, `internal/project/service.go`, `internal/project/handler.go`

## Task: Admin Dashboard Stats API

**Commits:**
- `feat: add admin service with dashboard stats`
- `feat: add admin handlers`
- `feat: mount admin routes in router`

**Endpoints:**
- `GET /admin/dashboard` - Get dashboard statistics (project count, task count, member count, etc.)

**Package:** `internal/admin/service.go`, `internal/admin/handler.go`

**Response:**
```json
{
  "success": true,
  "data": {
    "total_projects": 5,
    "total_tasks": 42,
    "total_members": 10,
    "tasks_by_status": {
      "backlog": 10,
      "todo": 15,
      "inprogress": 12,
      "done": 5
    }
  }
}
```

## Task: Members API

**Commits:**
- `feat: add members list handler`
- `feat: add member profile handlers`

**Endpoints:**
- `GET /members` - List all members
- `GET /members/:id` - Get member by ID
- `PUT /members/:id` - Update member profile

## Total Commits: 10
