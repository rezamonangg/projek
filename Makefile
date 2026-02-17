# Project Makefile - Unified commands for backend (Go) and frontend (SvelteKit)

.PHONY: dev backend frontend build build-fe start test test-be test-fe test-e2e test-e2e-auth test-e2e-admin test-e2e-project lint lint-be lint-fe migrate migrate-down migrate-new seed docker-up docker-down clean help

# =============================================================================
# Development - Start services in development mode with hot reload
# =============================================================================

## Run backend and frontend concurrently in development mode
dev:
	@echo "Starting backend and frontend in parallel..."
	@trap 'kill %1 %2 2>/dev/null; exit' INT; \
		cd backend && go run ./cmd/api & \
		cd frontend && npm run dev & \
		wait

## Run backend development server (hot reload via go run)
backend:
	@cd backend && go run ./cmd/api

## Run frontend development server (Vite hot reload)
frontend:
	@cd frontend && npm run dev

# =============================================================================
# Production Build
# =============================================================================

## Build backend binary to backend/bin/projek
build:
	@echo "Building backend binary..."
	@mkdir -p backend/bin
	@cd backend && go build -o ./bin/projek ./cmd/api
	@echo "Binary built: backend/bin/projek"

## Build frontend for production (outputs to frontend/build/)
build-fe:
	@echo "Building frontend..."
	@cd frontend && npm run build
	@echo "Frontend built: frontend/build/"

## Build and run production backend binary
start: build
	@echo "Starting production binary..."
	@./backend/bin/projek

# =============================================================================
# Testing
# =============================================================================

## Run all unit tests (backend + frontend)
test: test-be test-fe

## Run backend unit tests with race detector and coverage
test-be:
	@echo "Running backend unit tests..."
	@cd backend && go test -v -race -coverprofile=coverage.out ./...

## Run frontend unit tests
test-fe:
	@echo "Running frontend tests..."
	@cd frontend && npm test

## Run backend E2E tests (requires Docker for testcontainers)
test-e2e:
	@echo "Running backend E2E tests (this may take a few minutes)..."
	@cd backend && go test -v -tags=e2e -timeout 15m ./test/e2e/...

## Run auth E2E tests only
test-e2e-auth:
	@echo "Running auth E2E tests..."
	@cd backend && go test -v -tags=e2e -timeout 5m -run TestAuth ./test/e2e/...

## Run admin E2E tests only
test-e2e-admin:
	@echo "Running admin E2E tests..."
	@cd backend && go test -v -tags=e2e -timeout 5m -run TestAdmin ./test/e2e/...

## Run project flow E2E tests only
test-e2e-project:
	@echo "Running project flow E2E tests..."
	@cd backend && go test -v -tags=e2e -timeout 10m -run "Test(Project|Task|Epic|Board|Label)" ./test/e2e/...

# =============================================================================
# Code Quality
# =============================================================================

## Run all linters (backend + frontend)
lint: lint-be lint-fe

## Run golangci-lint on backend code
lint-be:
	@echo "Running backend linter..."
	@cd backend && golangci-lint run ./...

## Run ESLint on frontend code
lint-fe:
	@echo "Running frontend linter..."
	@cd frontend && npm run lint

# =============================================================================
# Database Migrations (using dbmate)
# =============================================================================

## Run all pending database migrations up
migrate:
	@echo "Running database migrations up..."
	@cd backend && dbmate --migrations-dir ./migrations -e DATABASE_URL up

## Rollback the last database migration
migrate-down:
	@echo "Rolling back last migration..."
	@cd backend && dbmate --migrations-dir ./migrations -e DATABASE_URL down

## Create a new database migration file
## Usage: make migrate-new NAME=create_users_table
migrate-new:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required"; \
		echo "Usage: make migrate-new NAME=create_users_table"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)..."
	@cd backend && dbmate --migrations-dir ./migrations -e DATABASE_URL new $(NAME)

## Seed development database with test data
seed:
	@echo "Running development seed..."
	@cd backend && ENV=development go run ./cmd/seed

# =============================================================================
# Infrastructure
# =============================================================================

## Start Docker services (Postgres, Redis, Prometheus)
docker-up:
	@docker-compose up -d

docker-db:
	@docker-compose up redis postgres -d

## Stop Docker services
docker-down:
	@docker-compose down

# =============================================================================
# Maintenance
# =============================================================================

## Remove build artifacts, coverage files, and test binaries
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf backend/bin/
	@rm -f backend/coverage.out
	@rm -rf frontend/build/
	@echo "Clean complete"

# =============================================================================
# Help
# =============================================================================

## Show this help message
help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make dev          Run backend + frontend concurrently"
	@echo "  make backend      Run backend dev server only"
	@echo "  make frontend     Run frontend dev server only"
	@echo ""
	@echo "Build:"
	@echo "  make build        Build backend binary to backend/bin/projek"
	@echo "  make build-fe     Build frontend to frontend/build/"
	@echo "  make start        Build and run production binary"
	@echo ""
	@echo "Testing:"
	@echo "  make test            Run all unit tests (backend + frontend)"
	@echo "  make test-be         Run backend unit tests"
	@echo "  make test-fe         Run frontend tests"
	@echo "  make test-e2e        Run all backend E2E tests (requires Docker)"
	@echo "  make test-e2e-auth   Run auth E2E tests only"
	@echo "  make test-e2e-admin  Run admin E2E tests only"
	@echo "  make test-e2e-project Run project flow E2E tests only"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint         Run all linters (backend + frontend)"
	@echo "  make lint-be      Run golangci-lint on backend"
	@echo "  make lint-fe      Run ESLint on frontend"
	@echo ""
	@echo "Database:"
	@echo "  make migrate         Run pending migrations up"
	@echo "  make migrate-down    Rollback last migration"
	@echo "  make migrate-new     Create new migration (NAME=required)"
	@echo "  make seed            Seed development database with test data"
	@echo ""
	@echo "Infrastructure:"
	@echo "  make docker-up    Start Postgres, Redis, Prometheus"
	@echo "  make docker-down  Stop Docker services"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean        Remove build artifacts and coverage files"
