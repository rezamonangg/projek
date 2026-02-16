.PHONY: build run test test-e2e lint migrate-up migrate-down migrate-new dev clean

BINARY_NAME=projek
BUILD_DIR=./bin

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

run: build
	@echo "Starting $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

dev:
	@go run ./cmd/api

test:
	@echo "Running unit tests..."
	@go test -v -race -coverprofile=coverage.out ./...

test-e2e:
	@echo "Running E2E tests..."
	@go test -v -tags=e2e ./test/...

lint:
	@echo "Running linter..."
	@golangci-lint run ./... || true

migrate-up:
	@echo "Running database migrations up..."
	@dbmate -e DATABASE_URL up

migrate-down:
	@echo "Rolling back database..."
	@dbmate -e DATABASE_URL down

migrate-new:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-new NAME=create_users_table"; exit 1; fi
	@echo "Creating migration $(NAME)..."
	@dbmate -e DATABASE_URL new $(NAME)

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@rm -f *.test

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

setup-dev:
	@echo "Setting up development environment..."
	@which dbmate > /dev/null || (echo "dbmate not found. Install: brew install dbmate" && exit 1)
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: brew install golangci-lint" && exit 1)
