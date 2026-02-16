.PHONY: build backend frontend run-all test test-e2e lint migrate-up migrate-down migrate-new dev clean docker-up docker-down

# Backend
BINARY_NAME=projek
BUILD_DIR=./backend/bin

# Root commands
backend-dev:
	cd backend && go run ./cmd/api

backend-build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@cd backend && go build -o ../$(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

backend-run: backend-build
	@echo "Starting $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

backend-test:
	@echo "Running backend unit tests..."
	@cd backend && go test -v -race -coverprofile=coverage.out ./...

backend-test-e2e:
	@echo "Running backend E2E tests..."
	@cd backend && go test -v -tags=e2e ./test/...

backend-lint:
	@echo "Running backend linter..."
	@cd backend && golangci-lint run ./... || true

backend-migrate-up:
	@echo "Running database migrations up..."
	@cd backend && dbmate -e DATABASE_URL up

backend-migrate-down:
	@echo "Rolling back database..."
	@cd backend && dbmate -e DATABASE_URL down

backend-migrate-new:
	@if [ -z "$(NAME)" ]; then echo "Usage: make backend-migrate-new NAME=create_users_table"; exit 1; fi
	@echo "Creating migration $(NAME)..."
	@cd backend && dbmate -e DATABASE_URL new $(NAME)

backend-clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@cd backend && rm -f coverage.out *.test

backend-deps:
	@echo "Installing backend dependencies..."
	@cd backend && go mod download && go mod tidy

backend-setup-dev:
	@echo "Setting up backend development environment..."
	@which dbmate > /dev/null || (echo "dbmate not found. Install: brew install dbmate" && exit 1)
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: brew install golangci-lint" && exit 1)

backend-seed:
	@echo "Running development seed..."
	@cd backend && ENV=development go run ./cmd/seed

# Frontend
frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

frontend-test:
	cd frontend && npm test

frontend-lint:
	cd frontend && npm run lint

# Combined
dev:
	@echo "Starting backend and frontend in parallel..."
	@trap 'kill %1 %2 2>/dev/null; exit' INT; \
		cd backend && go run ./cmd/api & \
		cd frontend && npm run dev & \
		wait

run:
	@echo "Building and running backend with frontend..."
	@trap 'kill %1 %2 2>/dev/null; exit' INT; \
		$(BUILD_DIR)/$(BINARY_NAME) & \
		cd frontend && sleep 2 && npm run dev & \
		wait

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Aliases for backward compatibility
build: backend-build
test: backend-test
test-e2e: backend-test-e2e
lint: backend-lint
migrate-up: backend-migrate-up
migrate-down: backend-migrate-down
migrate-new: backend-migrate-new
clean: backend-clean
