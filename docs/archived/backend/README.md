# Backend Summary

## Overview
Backend implementation for Projek - self-hosted project management suite.

## Tech Stack
- **Language:** Go 1.26
- **Router:** Chi (go-chi/chi)
- **Database:** PostgreSQL 15+ with pgx
- **Cache/Sessions:** Redis with go-redis
- **Migrations:** dbmate
- **Logging:** Zerolog
- **Metrics:** Prometheus client
- **Testing:** testify, uber/mock, testcontainers-go

## Architecture
- Domain-driven folder structure
- Clean architecture (repository, service, handler)
- Manual dependency injection
- Interface-based design for testability

## Phase Breakdown

| Phase | Description | Commits |
|-------|-------------|---------|
| 0 | Setup | 12 |
| 1 | Infrastructure | 22 |
| 2 | Authentication | 27 |
| 3 | Projects & Kanban | 42 |
| 4 | Wiki & Files | 26 |
| 5 | Observability | 18 |
| 6 | Testing | 25 |
| 7 | Deployment | 10 |

**Total Backend Commits: 182**

## Directory Structure
```
backend/
├── cmd/api/           # Entry point
├── internal/          # Application code
│   ├── auth/          # Authentication
│   ├── common/        # Shared utilities
│   ├── community/     # Community management
│   ├── member/        # Members & invitations
│   ├── project/       # Projects, epics, boards
│   ├── task/          # Tasks & labels
│   ├── wiki/          # Wiki pages
│   ├── file/          # File storage
│   ├── email/         # Email service
│   ├── admin/         # Admin APIs
│   ├── metrics/       # Prometheus metrics
│   └── router/        # HTTP router
├── migrations/        # dbmate migrations
├── test/              # E2E tests
└── pkg/               # Public packages (if any)
```

## Key Dependencies
```
github.com/go-chi/chi/v5
github.com/jackc/pgx/v5
github.com/redis/go-redis/v9
github.com/rs/zerolog
github.com/spf13/viper
github.com/go-playground/validator/v10
golang.org/x/crypto/bcrypt
github.com/prometheus/client_golang
```
