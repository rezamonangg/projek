# Phase 0: Project Setup

Setup project structure, tooling, and initial configuration.

## Task 0.1: Initialize Git Repository

**Commits:**
- `chore: initialize git repository`
- `chore: add .gitignore for Go, Node, and IDE files`
- `chore: add LICENSE file (MIT)`

## Task 0.2: Create Project Directory Structure

**Commits:**
- `chore: create backend directory structure`
- `chore: create frontend directory structure`
- `chore: create docs directory structure`

**Details:**
```
projek/
├── backend/
│   ├── cmd/api/
│   ├── internal/
│   ├── migrations/
│   ├── test/
│   └── pkg/
├── frontend/
│   ├── src/
│   ├── tests/
│   └── static/
└── docs/
```

## Task 0.3: Backend Go Module Initialization

**Commits:**
- `chore: initialize go module with go 1.26`
- `chore: add Makefile with common commands`
- `chore: add backend .gitignore`

**Details:**
- Module path: `github.com/monachy/projek`
- Create go.mod with Go 1.26

## Task 0.4: Frontend SvelteKit Initialization

**Commits:**
- `chore: initialize sveltekit project with TypeScript`
- `chore: add tailwindcss configuration`
- `chore: add vitest and playwright configurations`

**Details:**
- Use create-svelte with TypeScript template
- Install Tailwind CSS
- Setup Vitest for unit tests
- Setup Playwright for E2E tests

## Task 0.5: Docker Setup

**Commits:**
- `chore: add backend Dockerfile`
- `chore: add frontend Dockerfile`
- `chore: add docker-compose.yml for development`
- `chore: add .env.example file`

**Details:**
- Multi-stage Dockerfile for backend
- Node image for frontend dev
- PostgreSQL and Redis services
- Volume mounts for hot reload

## Task 0.6: Development Tooling

**Commits:**
- `chore: add air configuration for hot reload`
- `chore: add golangci-lint configuration`
- `chore: add prettier configuration for frontend`

## Task 0.7: Documentation Structure

**Commits:**
- `chore: add README with setup instructions`
- `chore: add CONTRIBUTING.md`
- `chore: add ARCHITECTURE.md`

---

## Total Commits: 17
**Estimated Time:** 1-2 days
