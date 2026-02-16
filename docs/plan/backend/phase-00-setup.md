# Backend Phase 0: Project Setup

Setup Go project structure, tooling, and initial configuration.

## Task 0.1: Initialize Git Repository

**Commits:**
- `chore: initialize git repository`
- `chore: add .gitignore for Go and IDE files`
- `chore: add LICENSE file (MIT)`

## Task 0.2: Create Directory Structure

**Commits:**
- `chore: create backend directory structure`
- `chore: create backend internal package structure`

**Structure:**
```
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── auth/
│   ├── common/
│   ├── community/
│   ├── member/
│   ├── project/
│   ├── task/
│   ├── wiki/
│   ├── file/
│   ├── email/
│   └── router/
├── migrations/
├── test/
└── pkg/
```

## Task 0.3: Go Module Initialization

**Commits:**
- `chore: initialize go module with go 1.26`
- `chore: add Makefile with common commands`
- `chore: add backend .gitignore`
- `chore: add golangci-lint configuration`

**Commands:**
```bash
cd backend
go mod init github.com/monachy/projek
go mod tidy
```

## Task 0.4: Makefile Setup

**Commits:**
- `chore: add Makefile with build, test, lint commands`
- `chore: add dbmate migration commands to Makefile`
- `chore: add development commands (run, watch)`

**Makefile targets:**
- `build` - Build binary
- `run` - Run development server
- `test` - Run unit tests
- `test-e2e` - Run E2E tests
- `lint` - Run golangci-lint
- `migrate-up` - Run dbmate up
- `migrate-down` - Run dbmate rollback
- `migrate-new` - Create new migration

## Task 0.5: Development Tooling

**Commits:**
- `chore: add air configuration for hot reload`
- `chore: add .env.example for backend`

---

## Total Backend Commits (Phase 0): 12
