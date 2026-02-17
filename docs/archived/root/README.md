# Projek - Complete Implementation Plan

## Overview
Self-hosted project management suite (ClickUp/Lark alternative)

## Three-Phase Development

### Phase 1: Backend (BE)
Go backend with all APIs
- **Commits:** 182
- **Duration:** 5-6 weeks
- **Output:** Complete REST API

### Phase 2: Frontend with Mock Data (FE)
SvelteKit frontend using mock APIs
- **Commits:** 139
- **Duration:** 4-5 weeks
- **Output:** Working UI with mock data
- **Parallel:** Can start after BE Phase 3

### Phase 3: Integration
Connect FE to BE, remove mocks
- **Commits:** 66
- **Duration:** 1-2 weeks
- **Output:** Full stack application

## Total
- **Total Commits:** 387
- **Total Duration:** 10-13 weeks
- **Team:** 1-2 developers (1 BE, 1 FE)

## Directory Structure
```
docs/plan/
├── README.md                 # This file
├── backend/                  # BE implementation plan
│   ├── README.md
│   ├── phase-00-setup.md
│   ├── phase-01-infrastructure.md
│   ├── phase-02-authentication.md
│   ├── phase-03-projects.md
│   ├── phase-04-wiki-files.md
│   ├── phase-05-observability.md
│   ├── phase-06-testing.md
│   └── phase-07-deployment.md
├── frontend/                 # FE implementation plan
│   ├── README.md
│   ├── phase-00-setup.md
│   ├── phase-01-mock-api.md
│   ├── phase-02-auth.md
│   ├── phase-03-dashboard.md
│   ├── phase-04-members.md
│   ├── phase-05-projects.md
│   ├── phase-06-epics.md
│   ├── phase-07-kanban.md
│   ├── phase-08-wiki.md
│   ├── phase-09-admin.md
│   └── phase-10-testing.md
└── integration/              # Integration plan
    ├── README.md
    ├── phase-00-setup.md
    ├── phase-01-auth.md
    ├── phase-02-members.md
    ├── phase-03-projects.md
    ├── phase-04-tasks.md
    ├── phase-05-wiki.md
    ├── phase-06-admin.md
    └── phase-07-cleanup.md
```

## Development Workflow

### Parallel Development
```
Week 1-2:  BE Phase 0-1 (Setup, Infrastructure)
Week 3-4:  BE Phase 2 (Auth) + FE starts Phase 0
Week 5-6:  BE Phase 3 (Projects) + FE Phase 1-2
Week 7-8:  BE Phase 4 (Wiki) + FE Phase 3-5
Week 9-10: BE Phase 5-6 + FE Phase 6-8
Week 11-12: BE Phase 7 + FE Phase 9-10
Week 13:   Integration
```

### Commit Convention
- `chore:` - Setup, configuration, tooling
- `feat:` - New features
- `task:` - Implementation tasks, refactoring
- `fix:` - Bug fixes
- `docs:` - Documentation
- `test:` - Tests only
- `refactor:` - Code restructuring

## Tech Stack

### Backend
- **Language:** Go 1.26
- **Router:** Chi
- **Database:** PostgreSQL 15+ with pgx
- **Sessions:** Redis
- **Migrations:** dbmate
- **Logging:** Zerolog
- **Testing:** testify, uber/mock, testcontainers-go

### Frontend
- **Framework:** SvelteKit
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **State:** Svelte Stores
- **Editor:** TipTap
- **Drag & Drop:** @dnd-kit
- **Testing:** Vitest + Playwright

## Getting Started

1. Review all plan files in backend/, frontend/, integration/
2. Start with `backend/phase-00-setup.md`
3. Follow commits in order
4. Test each phase before next

## Success Criteria
- [ ] All backend endpoints working
- [ ] Frontend with mock data fully functional
- [ ] Integration complete, all real APIs connected
- [ ] All tests passing
- [ ] Docker Compose deployment working
