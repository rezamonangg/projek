# Backend Implementation Order & Priority

This document outlines the recommended order for implementing missing backend modules.

## Priority Matrix

| Module | Priority | Blockers | Dependencies | Est. Time |
|--------|----------|----------|--------------|-----------|
| Projects | **P0** | Blocks project creation | None | 1-2h |
| Boards | **P1** | Blocks task viewing | Projects | 1h |
| Tasks | **P1** | Blocks task management | Boards | 1.5h |
| Epics | **P2** | Feature | Projects | 1h |
| Labels | **P2** | Feature | Tasks | 45m |
| Members | **P2** | Feature | None | 1h |
| Wiki | **P3** | Feature | None | 2h |
| Files | **P3** | Feature | None | 2-3h |
| Admin | **P3** | Feature | Members, Projects, Tasks | 1.5h |

## Recommended Implementation Order

### Phase 1: Core Flow (Must Have)

```
1. Projects Handler (01-projects-handler.md)
   └─ Enables: Project creation, listing
   
2. Boards Handler (03-boards-handler.md)
   └─ Enables: Board listing, task container
   
3. Tasks Handler (04-tasks-handler.md)
   └─ Enables: Full task management
```

### Phase 2: Enhanced Features (Should Have)

```
4. Epics Handler (02-epics-handler.md)
   └─ Enables: Epic management per project
   
5. Labels Handler (05-labels-handler.md)
   └─ Enables: Label creation and assignment
   
6. Members Handler (06-members-handler.md)
   └─ Enables: Member management, invitations
```

### Phase 3: Additional Features (Nice to Have)

```
7. Wiki Module (07-wiki-module.md)
   └─ Enables: Project documentation
   
8. Files Module (08-files-module.md)
   └─ Enables: File uploads and attachments
   
9. Admin Module (09-admin-module.md)
   └─ Enables: Admin dashboard and settings
```

## Quick Start (Phase 1 Only)

To unblock the project creation issue immediately, implement in this order:

```bash
# 1. Create project handler
touch backend/internal/project/handler.go

# 2. Modify router to mount /projects
# Edit backend/internal/router/router.go

# 3. Test
cd backend && go build ./... && go test ./...
```

## Route Registration Summary

After all implementations, router.go should include:

```go
func NewRouter(cfg *common.Config, logger zerolog.Logger, db *pgxpool.Pool, redis *redis.Client) *chi.Mux {
    r := chi.NewRouter()
    
    // ... middleware ...
    
    // Health endpoints
    r.Get("/health", healthHandler)
    r.Get("/ready", readinessHandler(db, redis))
    
    // Auth module (EXISTS)
    r.Mount("/auth", authHandler.Routes())
    
    // Initialize repositories
    memberRepo := member.NewRepository(db)
    projectRepo := project.NewRepository(db)
    projectEpicRepo := project.NewEpicRepository(db)
    projectBoardRepo := project.NewBoardRepository(db)
    taskRepo := task.NewRepository(db)
    taskLabelRepo := task.NewLabelRepository(db)
    
    // Initialize services
    memberService := member.NewService(memberRepo)
    projectService := project.NewService(projectRepo, projectEpicRepo, projectBoardRepo, taskRepo, taskLabelRepo)
    taskService := task.NewService(taskRepo, taskLabelRepo)
    labelService := task.NewLabelService(taskLabelRepo, taskRepo)
    wikiService := wiki.NewService(wiki.NewRepository(db))
    fileService := file.NewService(file.NewRepository(db), file.NewLocalStorage(cfg.UploadPath))
    adminService := admin.NewService(admin.NewRepository(db), admin.NewDBStatsProvider(db))
    
    // Mount routes
    r.Mount("/members", member.NewHandler(memberService).Routes())
    r.Mount("/projects", project.NewHandler(projectService).Routes())
    r.Mount("/epics", project.NewEpicHandler(projectService).Routes())
    r.Mount("/boards", project.NewBoardHandler(projectService).Routes())
    r.Mount("/tasks", task.NewHandler(taskService).Routes())
    r.Mount("/wiki", wiki.NewHandler(wikiService).Routes())
    r.Mount("/files", file.NewHandler(fileService).Routes())
    r.Mount("/admin", admin.NewHandler(adminService).Routes())
    
    // Nested routes
    r.Route("/projects/{projectId}", func(r chi.Router) {
        r.Mount("/epics", project.NewEpicHandler(projectService).ProjectRoutes())
        r.Mount("/boards", project.NewBoardHandler(projectService).ProjectRoutes())
        r.Mount("/labels", task.NewLabelHandler(labelService).ProjectRoutes())
        r.Mount("/wiki", wiki.NewHandler(wikiService).ProjectRoutes())
    })
    
    r.Route("/boards/{boardId}", func(r chi.Router) {
        r.Mount("/tasks", task.NewHandler(taskService).BoardRoutes())
    })
    
    r.Route("/tasks/{taskId}", func(r chi.Router) {
        r.Mount("/labels", task.NewLabelHandler(labelService).TaskRoutes())
    })
    
    return r
}
```

## Phase 4: E2E Testing (Quality Assurance)

```
10. E2E Shared Setup (10-e2e-shared-setup.md)
   └─ Testcontainers infrastructure
   └─ Isolated database per test suite
   └─ Test client with session handling
   └─ Migration runner

11. E2E Auth Flow (11-e2e-auth-flow.md)
   └─ Register community
   └─ Login as admin
   └─ Logout
   └─ Protected endpoints

12. E2E Admin Flow (12-e2e-admin-flow.md)
   └─ Complete 7-phase workflow:
      - Auth (register, login)
      - Project creation
      - Task creation (backlog, todo)
      - Member invitation
      - Dashboard stats
      - Label management
      - Epic management
   └─ Data integrity verification

13. E2E Project Flow (13-e2e-project-flow.md)
   └─ Project CRUD
   └─ Task lifecycle (backlog → done)
   └─ Epic-task relationships
   └─ Board management
```

## Total Estimated Time

- **Phase 1:** 4-4.5 hours (unblocks core functionality)
- **Phase 2:** 2.75 hours (enhanced features)
- **Phase 3:** 5.5 hours (additional features)
- **Phase 4:** 4-6 hours (E2E testing)
  - Setup: 1 hour
  - Auth flow: 1 hour
  - Admin flow: 1.5 hours
  - Project flow: 1 hour
- **Total:** ~16-19 hours for complete backend with tests
