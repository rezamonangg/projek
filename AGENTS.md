# AGENTS.md

Coding guidelines for AI agents working on this repository.

## Project Overview

Full-stack project management application with Go backend and SvelteKit frontend.

**Stack:**
- Backend: Go 1.25+ (Chi router, PostgreSQL via pgx, Redis, Zerolog)
- Frontend: SvelteKit 2.x + Svelte 5 (runes), TypeScript 5.x, Tailwind CSS 4.x
- Services: PostgreSQL 15, Redis 7, Prometheus (via Docker Compose)

---

## Build / Test / Lint Commands

### Backend (Go)

```bash
# Development
cd backend && go run ./cmd/api              # Run dev server
make backend-dev                             # Via Makefile

# Build
make backend-build                           # Build binary to backend/bin/projek

# Tests
cd backend && go test ./...                  # Run all tests
cd backend && go test -v ./internal/auth/... # Run tests in package
cd backend && go test -run TestName ./...    # Run single test
cd backend && go test -v -race -coverprofile=coverage.out ./...  # With race/coverage
make backend-test                            # Via Makefile

# E2E Tests
cd backend && go test -v -tags=e2e ./test/...
make backend-test-e2e

# Linting
cd backend && golangci-lint run ./...
make backend-lint

# Dependencies
cd backend && go mod tidy && go mod download
```

### Frontend (SvelteKit)

```bash
# Development
cd frontend && npm run dev                   # Start dev server (Vite)

# Build
cd frontend && npm run build                 # Production build
cd frontend && npm run check                 # TypeScript + Svelte check

# Tests
cd frontend && npm test                      # Run Playwright tests
cd frontend && npm run test:ui               # Playwright UI mode
cd frontend && npm run test:headed           # Playwright headed mode

# Linting & Formatting
cd frontend && npm run lint                  # ESLint check
cd frontend && npm run lint:fix              # ESLint fix
cd frontend && npm run format                # Prettier format
cd frontend && npm run format:check          # Prettier check
```

### Combined (Root Makefile)

```bash
make dev                                     # Start both backend and frontend
make docker-up                               # Start Docker services (Postgres, Redis, Prometheus)
make docker-down                             # Stop Docker services
```

---

## Code Style Guidelines

### Go Backend

**Imports:** Group stdlib first, then third-party, then internal packages. Use goimports.

```go
import (
    "context"
    "net/http"
    
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    
    "github.com/monachy/projek/internal/common"
    "github.com/monachy/projek/internal/member"
)
```

**Naming:**
- PascalCase for exported (Handler, Service, Repository)
- camelCase for unexported
- Interface names: Repository, Service (suffix with type)
- Test files: `*_test.go`, mocks: `mock_*.go`

**Types:**
- Use UUIDs from `github.com/google/uuid` for IDs
- JSON tags: snake_case (`json:"created_at"`)
- DB tags: snake_case (`db:"password_hash"`)
- Validation: use `validate:"required,email"` tags with go-playground/validator

**Error Handling:**
- Return errors, don't panic in business logic
- Use `common.Error()` for HTTP responses
- Wrap errors with context: `fmt.Errorf("failed to create: %w", err)`
- Repository errors bubble up to service layer

**Structure:**
```go
// internal/domain/model.go
type Model struct { ... }
type CreateInput struct { ... }

// internal/domain/repository.go
type Repository interface { ... }

// internal/domain/service.go
type Service struct { ... }

// internal/domain/handler.go
type Handler struct { ... }
func (h *Handler) Routes() chi.Router { ... }
```

**Logging:** Use zerolog via `internal/common/logger.go`:
```go
log.Info().Str("user_id", id.String()).Msg("user created")
log.Error().Err(err).Str("context", "auth").Msg("login failed")
```

### TypeScript / Svelte Frontend

**Formatting:**
- Tabs (not spaces)
- Single quotes
- No trailing commas
- 100 char line width
- Prettier + prettier-plugin-svelte

**Svelte 5 (Runes):**
```svelte
<script lang="ts">
	let count = $state(0);              // Reactive state
	let doubled = $derived(count * 2);  // Derived values
	
	$effect(() => {                     // Effects
		console.log('count changed:', count);
	});
</script>
```

**TypeScript:**
- Strict mode enabled
- Explicit types for function params and returns
- Interfaces for props: `interface Props { variant?: 'primary' | 'secondary' }`
- Use `$lib` alias: `import { Button } from '$lib/components'`

**Naming:**
- PascalCase: Components, Types, Interfaces
- camelCase: variables, functions, props
- Svelte components: PascalCase files

**Imports:**
```typescript
// Svelte imports first
import { onMount } from 'svelte';

// Type imports
import type { Member } from '$lib/types';

// Component imports
import { Button } from '$lib/components/common';

// API/utility imports
import { memberApi } from '$lib/api';
```

**Error Handling:**
- Use try/catch for async operations
- Check API response structure before accessing data
- Use optional chaining: `data?.items?.length`

**API Calls:** Pattern in `$lib/api/*.ts`
```typescript
export const projectApi = {
	async list(): Promise<ApiResponse<Project[]>> {
		const res = await fetch('/api/projects');
		return res.json();
	}
};
```

---

## Testing Guidelines

### Go
- Table-driven tests preferred
- Use testify: `assert`, `require`
- Mock external deps with go.uber.org/mock
- Test naming: `TestService_MethodName`

### Frontend
- Playwright for E2E: `tests/*.spec.ts`
- Vitest for unit: `src/**/*.{test,spec}.ts`
- Mock API calls in unit tests

---

## Common Patterns

**Backend Handler Response:**
```go
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    item, err := h.service.Get(r.Context(), uuid.MustParse(id))
    if err != nil {
        common.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
        return
    }
    common.Success(w, http.StatusOK, item)
}
```

**Frontend API Call:**
```typescript
async function load() {
	const response = await api.getItem(id);
	if (response.data) {
		item = response.data;
	}
}
```

---

## Git Workflow

- Branch: `feature/description` or `fix/description`
- Commit: present tense, lowercase after prefix
- Example: `feat: add user authentication`

## Environment Setup

Copy example env files before first run:
```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

Required: Go 1.25+, Node.js 20+, Docker, dbmate, golangci-lint
