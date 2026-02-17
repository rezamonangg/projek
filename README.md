# Projek

Full-stack project management application with Go backend and SvelteKit frontend.

## Tech Stack

**Backend:**
- Go 1.25+ with Chi router
- PostgreSQL (via pgx)
- Redis for caching/sessions
- dbmate for migrations
- golangci-lint for linting

**Frontend:**
- SvelteKit 2.x + Svelte 5 (runes)
- TypeScript 5.x
- Tailwind CSS 4.x
- Playwright for E2E testing

**Infrastructure:**
- Docker Compose (Postgres, Redis, Prometheus)

## Project Structure

```
projek/
├── backend/              # Go backend
│   ├── cmd/
│   │   ├── api/         # Main API server entrypoint
│   │   └── seed/        # Database seeding utility
│   ├── internal/        # Internal packages
│   │   ├── auth/        # Authentication handlers
│   │   ├── common/      # Shared utilities (logger, responses)
│   │   ├── member/      # Member domain (handlers, service, repository)
│   │   └── middleware/  # HTTP middleware
│   ├── migrations/      # Database migrations (dbmate)
│   ├── test/           # E2E tests
│   └── go.mod
├── frontend/           # SvelteKit frontend
│   ├── src/
│   │   ├── lib/        # Shared components, utilities, API clients
│   │   ├── routes/     # SvelteKit routes
│   │   └── app.html    # HTML template
│   ├── static/         # Static assets
│   ├── tests/          # Playwright tests
│   └── package.json
├── docs/               # Documentation
├── docker-compose.yml  # Docker services
└── Makefile           # Unified commands
```

## Prerequisites

- Go 1.25+
- Node.js 20+
- Docker & Docker Compose
- dbmate (`brew install dbmate`)
- golangci-lint (`brew install golangci-lint`)

## Quick Start

1. **Copy environment files:**
   ```bash
   cp backend/.env.example backend/.env
   cp frontend/.env.example frontend/.env
   ```
   
   Note: Frontend uses `.env` for development and `.env.production` for production builds (handled automatically by SvelteKit).

2. **Start infrastructure:**
   ```bash
   make docker-up
   ```

3. **Run database migrations:**
   ```bash
   make migrate
   ```

4. **Start development servers:**
   ```bash
   make dev
   ```
   This runs both backend and frontend concurrently with hot reload.

## Makefile Commands

### Development

| Command | Description |
|---------|-------------|
| `make dev` | Run backend + frontend concurrently with hot reload |
| `make backend` | Run backend dev server only (`go run ./cmd/api`) |
| `make frontend` | Run frontend dev server only (`npm run dev`) |

### Build

| Command | Description |
|---------|-------------|
| `make build` | Build backend binary to `backend/bin/projek` |
| `make build-fe` | Build frontend for production (outputs to `frontend/build/`) |
| `make start` | Build and run production backend binary |

### Testing

| Command | Description |
|---------|-------------|
| `make test` | Run ALL unit tests (backend + frontend) |
| `make test-be` | Run backend unit tests with race detector and coverage |
| `make test-fe` | Run frontend tests |
| `make test-e2e` | Run backend E2E tests (requires database) |

### Code Quality

| Command | Description |
|---------|-------------|
| `make lint` | Run ALL linters (backend + frontend) |
| `make lint-be` | Run golangci-lint on backend |
| `make lint-fe` | Run ESLint on frontend |

### Database

| Command | Description |
|---------|-------------|
| `make migrate` | Run all pending database migrations up |
| `make migrate-down` | Rollback the last database migration |
| `make migrate-new NAME=create_users` | Create a new migration file |
| `make seed` | Seed development database with test data |

### Infrastructure

| Command | Description |
|---------|-------------|
| `make docker-up` | Start Postgres, Redis, Prometheus in Docker |
| `make docker-down` | Stop Docker services |

### Maintenance

| Command | Description |
|---------|-------------|
| `make clean` | Remove build artifacts and coverage files |
| `make help` | Show all available commands |

## Development Workflow

### Adding a Feature

1. Create database migration (if needed):
   ```bash
   make migrate-new NAME=add_user_roles
   make migrate
   ```

2. Implement backend:
   - Add domain logic in `backend/internal/<domain>/`
   - Follow existing patterns (handler → service → repository)
   - Add tests in `*_test.go` files

3. Implement frontend:
   - Add API client in `frontend/src/lib/api/`
   - Create components in `frontend/src/lib/components/`
   - Add routes in `frontend/src/routes/`

4. Run tests:
   ```bash
   make test
   make lint
   ```

### Common Tasks

**Reset database:**
```bash
make migrate-down  # Rollback
make migrate       # Re-apply
make seed          # Seed data
```

**Run only backend tests:**
```bash
make test-be
```

**Check code before commit:**
```bash
make lint
make test
```

## Architecture

### Backend (Clean Architecture)

Each domain follows a layered approach:

```
internal/member/
├── handler.go      # HTTP handlers (input/output)
├── service.go      # Business logic
├── repository.go   # Database operations
└── model.go        # Domain models
```

**Request flow:** HTTP Request → Handler → Service → Repository → Database

### Frontend (SvelteKit)

- **Routes**: File-based routing in `src/routes/`
- **API Clients**: Centralized in `src/lib/api/`
- **Components**: Reusable UI in `src/lib/components/`
- **State**: Svelte 5 runes (`$state`, `$derived`, `$effect`)

## Environment Variables

See `.env.example` files in `backend/` and `frontend/` directories for required variables.

## Contributing

1. Follow existing code style (see AGENTS.md for detailed guidelines)
2. Write tests for new features
3. Run `make lint` and `make test` before committing
4. Use conventional commit messages

## License

[Your License Here]
