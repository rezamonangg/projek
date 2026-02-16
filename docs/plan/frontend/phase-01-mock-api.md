# Frontend Phase 1: Mock API Layer

Create mock API layer for frontend development without backend.

## Task 1.1: Mock API Client Setup

**Commits:**
- `feat: add mock api client interface`
- `feat: add mock data generators`
- `feat: add mock delay simulation`

## Task 1.2: Auth Mock APIs

**Commits:**
- `feat: add mock login api`
- `feat: add mock logout api`
- `feat: add mock current user api`
- `feat: add mock token storage`

**Mock Data:**
```typescript
const mockUser = {
  id: 'uuid',
  email: 'admin@example.com',
  firstName: 'Admin',
  lastName: 'User',
  role: 'admin',
  communityId: 'uuid'
};
```

## Task 1.3: Member Mock APIs

**Commits:**
- `feat: add mock list members api`
- `feat: add mock invite member api`
- `feat: add mock update profile api`

**Mock Data:**
```typescript
const mockMembers = [
  { id: '1', email: 'user1@example.com', firstName: 'John', lastName: 'Doe', role: 'admin' },
  { id: '2', email: 'user2@example.com', firstName: 'Jane', lastName: 'Smith', role: 'member' }
];
```

## Task 1.4: Project Mock APIs

**Commits:**
- `feat: add mock list projects api`
- `feat: add mock create project api`
- `feat: add mock update project api`
- `feat: add mock delete project api`

**Mock Data:**
```typescript
const mockProjects = [
  { 
    id: '1', 
    name: 'Website Redesign', 
    description: 'Redesign company website',
    status: 'active',
    createdAt: '2024-01-01'
  }
];
```

## Task 1.5: Epic Mock APIs

**Commits:**
- `feat: add mock list epics api`
- `feat: add mock create epic api`
- `feat: add mock update epic api`

## Task 1.6: Board Mock APIs

**Commits:**
- `feat: add mock list boards api`
- `feat: add mock create board api`
- `feat: add mock board details api`

## Task 1.7: Task Mock APIs

**Commits:**
- `feat: add mock list tasks api`
- `feat: add mock create task api`
- `feat: add mock update task api`
- `feat: add mock move task api`
- `feat: add mock delete task api`

**Mock Data:**
```typescript
const mockTasks = [
  {
    id: '1',
    title: 'Design homepage',
    description: 'Create Figma designs',
    status: 'inprogress',
    assigneeId: '1',
    boardId: '1',
    epicId: '1',
    labels: ['design', 'high-priority']
  }
];
```

## Task 1.8: Label Mock APIs

**Commits:**
- `feat: add mock list labels api`
- `feat: add mock create label api`
- `feat: add mock assign label to task api`

## Task 1.9: Wiki Mock APIs

**Commits:**
- `feat: add mock list wiki pages api`
- `feat: add mock get wiki page api`
- `feat: add mock create wiki page api`
- `feat: add mock update wiki page api`

**Mock Data:**
```typescript
const mockWikiPages = [
  {
    id: '1',
    title: 'Getting Started',
    slug: 'getting-started',
    content: { /* TipTap JSON */ },
    projectId: '1'
  }
];
```

## Task 1.10: File Mock APIs

**Commits:**
- `feat: add mock upload file api`
- `feat: add mock download file api`
- `feat: add mock list attachments api`

## Task 1.11: Admin Mock APIs

**Commits:**
- `feat: add mock get community settings api`
- `feat: add mock update community settings api`
- `feat: add mock get dashboard stats api`

**Mock Data:**
```typescript
const mockDashboardStats = {
  totalMembers: 10,
  totalProjects: 5,
  totalTasks: 150,
  tasksByStatus: {
    backlog: 20,
    todo: 30,
    inprogress: 40,
    done: 60
  }
};
```

---

## Total Frontend Commits (Phase 1): 24
