# Frontend Implementation Order & Priority

This document outlines the recommended order for implementing frontend API modules based on the backend handlers.

## Priority Matrix

| Module | Priority | Backend Dependency | Est. Time |
|--------|----------|-------------------|-----------|
| Projects | **P0** | Projects Handler | 30m |
| Boards | **P1** | Boards Handler | 20m |
| Tasks | **P1** | Tasks Handler | 30m |
| Epics | **P2** | Epics Handler | 20m |
| Labels | **P2** | Labels Handler | 20m |
| Members | **P2** | Members Handler | 30m |
| Wiki | **P3** | Wiki Module | 30m |
| Files | **P3** | Files Module | 45m |
| Admin | **P3** | Admin Module | 30m |

## Recommended Implementation Order

### Phase 1: Core Flow (Must Have)

```
1. Projects API (01-projects-api.md)
   └─ Existing: real/projects.ts (UPDATE if needed)
   └─ Verify: All endpoints match backend contract
   
2. Boards API (02-boards-api.md)
   └─ Existing: real/boards.ts (UPDATE if needed)
   └─ Verify: Endpoints match /projects/{projectId}/boards
   
3. Tasks API (03-tasks-api.md)
   └─ Existing: real/tasks.ts (UPDATE if needed)
   └─ Verify: Endpoints match /boards/{boardId}/tasks
```

### Phase 2: Enhanced Features (Should Have)

```
4. Epics API (04-epics-api.md)
   └─ Existing: real/epics.ts (UPDATE if needed)
   └─ Verify: Endpoints match /projects/{projectId}/epics
   
5. Labels API (05-labels-api.md)
   └─ Existing: real/labels.ts (UPDATE if needed)
   └─ Verify: Endpoints match /projects/{projectId}/labels
   
6. Members API (06-members-api.md)
   └─ Existing: real/members.ts (UPDATE if needed)
   └─ Verify: Pagination matches backend response
```

### Phase 3: Additional Features (Nice to Have)

```
7. Wiki API (07-wiki-api.md)
   └─ Existing: real/wiki.ts (UPDATE if needed)
   └─ Verify: Endpoints match /projects/{projectId}/wiki
   
8. Files API (08-files-api.md)
   └─ Existing: real/files.ts (UPDATE if needed)
   └─ Verify: FormData upload matches backend
   
9. Admin API (09-admin-api.md)
   └─ Existing: real/admin.ts (UPDATE if needed)
   └─ Verify: Settings and dashboard endpoints
```

## API Endpoint Summary

| Frontend API | Backend Endpoint | Status |
|--------------|------------------|--------|
| `listProjects()` | `GET /projects` | EXISTS |
| `getProject(id)` | `GET /projects/{id}` | EXISTS |
| `createProject(input)` | `POST /projects` | EXISTS |
| `updateProject(id, input)` | `PUT /projects/{id}` | EXISTS |
| `deleteProject(id)` | `DELETE /projects/{id}` | EXISTS |
| `listBoards(projectId)` | `GET /projects/{projectId}/boards` | EXISTS |
| `getBoard(id)` | `GET /boards/{id}` | EXISTS |
| `createBoard(input)` | `POST /projects/{projectId}/boards` | EXISTS |
| `listTasks(boardId)` | `GET /boards/{boardId}/tasks` | EXISTS |
| `createTask(input)` | `POST /boards/{boardId}/tasks` | EXISTS |
| `updateTask(id, input)` | `PUT /tasks/{id}` | EXISTS |
| `moveTask(id, input)` | `PATCH /tasks/{id}/move` | EXISTS |
| `deleteTask(id)` | `DELETE /tasks/{id}` | EXISTS |
| `listEpics(projectId)` | `GET /projects/{projectId}/epics` | EXISTS |
| `createEpic(input)` | `POST /projects/{projectId}/epics` | EXISTS |
| `updateEpic(id, input)` | `PUT /epics/{id}` | EXISTS |
| `deleteEpic(id)` | `DELETE /epics/{id}` | EXISTS |
| `listLabels(projectId)` | `GET /projects/{projectId}/labels` | EXISTS |
| `createLabel(input)` | `POST /projects/{projectId}/labels` | EXISTS |
| `assignLabelToTask(taskId, labelId)` | `POST /tasks/{taskId}/labels` | EXISTS |
| `listMembers()` | `GET /members` | EXISTS |
| `inviteMember(input)` | `POST /members/invite` | EXISTS |
| `updateProfile(userId, input)` | `PUT /members/{id}` | EXISTS |
| `listWikiPages(projectId)` | `GET /projects/{projectId}/wiki` | EXISTS |
| `getWikiPage(id)` | `GET /wiki/{id}` | EXISTS |
| `createWikiPage(input)` | `POST /projects/{projectId}/wiki` | EXISTS |
| `updateWikiPage(id, input)` | `PUT /wiki/{id}` | EXISTS |
| `deleteWikiPage(id)` | `DELETE /wiki/{id}` | EXISTS |
| `uploadFile(file, options)` | `POST /files/upload` | EXISTS |
| `downloadFile(id)` | `GET /files/{id}` | EXISTS |
| `listAttachments(options)` | `GET /files?task_id=...` | EXISTS |
| `getCommunitySettings()` | `GET /admin/settings` | EXISTS |
| `updateCommunitySettings(input)` | `PUT /admin/settings` | EXISTS |
| `getDashboardStats()` | `GET /admin/dashboard` | EXISTS |

## File Structure

```
frontend/src/lib/api/
├── index.ts              # Export all APIs
├── types.ts              # Shared types and interfaces
├── loading.ts            # Loading state API (mock)
├── real/
│   ├── client.ts         # HTTP client with retry logic
│   ├── transform.ts      # Snake_case to camelCase transformers
│   ├── auth.ts           # Auth API implementation
│   ├── projects.ts       # Projects API implementation
│   ├── boards.ts         # Boards API implementation
│   ├── tasks.ts          # Tasks API implementation
│   ├── epics.ts          # Epics API implementation
│   ├── labels.ts         # Labels API implementation
│   ├── members.ts        # Members API implementation
│   ├── wiki.ts           # Wiki API implementation
│   ├── files.ts          # Files API implementation
│   └── admin.ts          # Admin API implementation
```

## Implementation Checklist

Before marking any API as complete, verify:

- [ ] Backend endpoint exists and is tested
- [ ] TypeScript types match backend response
- [ ] Transformer functions handle snake_case → camelCase
- [ ] Error handling covers all HTTP status codes
- [ ] Loading states are properly managed
- [ ] API is exported in `index.ts`

## Total Estimated Time

- **Phase 1:** 1.5 hours (core functionality)
- **Phase 2:** 1 hour (enhanced features)
- **Phase 3:** 1.5 hours (additional features)
- **Total:** ~4 hours for complete frontend API alignment
