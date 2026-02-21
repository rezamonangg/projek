# Projek

> Full-stack project management application with Go backend and SvelteKit frontend.

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-2.x-FF3E00?style=flat-square&logo=svelte)](https://kit.svelte.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=flat-square&logo=typescript)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=flat-square&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

---

## Quick Start

```bash
# 1. Setup environment
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# 2. Start infrastructure (Postgres, Redis)
make docker-up

# 3. Run migrations
make migrate

# 4. Start development servers
make dev
```

That's it. Backend runs on `:8080`, frontend on `:5173`.

---

## Tech Stack

| Layer | Technologies |
|-------|-------------|
| **Backend** | Go 1.25+, Chi router, pgx (PostgreSQL), Redis, Zerolog |
| **Frontend** | SvelteKit 2.x, Svelte 5 (runes), TypeScript 5.x, Tailwind CSS 4.x |
| **Infrastructure** | Docker Compose (Postgres 15, Redis 7, Prometheus) |
| **Tools** | dbmate (migrations), golangci-lint, Playwright (E2E) |

---

## Project Structure

```
projek/
├── backend/                 # Go API server
│   ├── cmd/
│   │   ├── api/            # Entry point
│   │   └── seed/           # DB seeding
│   ├── internal/           # Domain packages
│   │   ├── auth/           # Authentication
│   │   ├── member/         # Member domain
│   │   ├── common/         # Shared utilities
│   │   └── middleware/     # HTTP middleware
│   ├── migrations/         # dbmate migrations
│   └── test/              # E2E tests
├── frontend/              # SvelteKit app
│   ├── src/
│   │   ├── lib/           # Components, API, utilities
│   │   └── routes/        # Page routes
│   └── tests/             # Playwright tests
├── docs/                  # Documentation
├── docker-compose.yml     # Services definition
└── Makefile              # Unified commands
```

---

## Commands

### Development

| Command | Description |
|---------|-------------|
| `make dev` | Start backend + frontend with hot reload |
| `make backend` | Run backend only |
| `make frontend` | Run frontend only |

### Build & Deploy

| Command | Description |
|---------|-------------|
| `make build` | Build backend binary to `backend/bin/projek` |
| `make build-fe` | Build frontend to `frontend/build/` |
| `make start` | Run production binary |

### Testing

| Command | Description |
|---------|-------------|
| `make test` | Run all unit tests (backend + frontend) |
| `make test-be` | Backend tests with race detection & coverage |
| `make test-fe` | Frontend tests |
| `make test-e2e` | E2E tests (requires running DB) |

### Code Quality

| Command | Description |
|---------|-------------|
| `make lint` | Run all linters |
| `make lint-be` | golangci-lint |
| `make lint-fe` | ESLint |

### Database

| Command | Description |
|---------|-------------|
| `make migrate` | Apply pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-new NAME=x` | Create new migration |
| `make seed` | Seed development data |

### Infrastructure

| Command | Description |
|---------|-------------|
| `make docker-up` | Start Postgres, Redis, Prometheus |
| `make docker-db` | Start only Postgres + Redis |
| `make docker-down` | Stop all services |

---

## Architecture

### Backend (Clean Architecture)

Each domain follows layered separation:

```
internal/member/
├── handler.go      # HTTP handlers (I/O)
├── service.go      # Business logic
├── repository.go   # Database operations
└── model.go        # Domain models
```

**Flow:** `HTTP Request → Handler → Service → Repository → Database`

### Frontend (SvelteKit)

- **Routes**: File-based routing in `src/routes/`
- **API**: Centralized clients in `src/lib/api/`
- **Components**: Reusable UI in `src/lib/components/`
- **State**: Svelte 5 runes (`$state`, `$derived`, `$effect`)

---

## Development Workflow

### Add a New Feature

1. **Create migration** (if schema changes):
   ```bash
   make migrate-new NAME=add_feature_table
   make migrate
   ```

2. **Implement backend**:
   - Add domain in `backend/internal/<domain>/`
   - Follow: handler → service → repository
   - Write tests in `*_test.go`

3. **Implement frontend**:
   - Add API client in `frontend/src/lib/api/`
   - Create components in `frontend/src/lib/components/`
   - Add routes in `frontend/src/routes/`

4. **Verify**:
   ```bash
   make lint
   make test
   ```

### Reset Database

```bash
make migrate-down && make migrate && make seed
```

---

## Prerequisites

- Go 1.25+
- Node.js 20+
- Docker & Docker Compose
- dbmate (`brew install dbmate`)
- golangci-lint (`brew install golangci-lint`)

---

## Environment Variables

See `.env.example` files in `backend/` and `frontend/` directories.

---

## Contributing

1. Follow code style (see [AGENTS.md](AGENTS.md))
2. Write tests for new features
3. Run `make lint && make test` before committing
4. Use conventional commits (`feat:`, `fix:`, `docs:`)

---

## License

[MIT](LICENSE)
