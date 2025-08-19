## Makefile
.PHONY: help dev build run test clean

help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Start development server with hot reload
	@echo "Starting development server..."
	@air

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/attendance-server cmd/server/main.go

run: build ## Run the built application
	@echo "Starting server..."
	@./bin/attendance-server

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ tmp/

db-start: ## Start PostgreSQL with Docker
	@echo "Starting PostgreSQL..."
	@docker run --name attendance-postgres -e POSTGRES_DB=attendanceDB -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15-alpine

db-stop: ## Stop PostgreSQL Docker container
	@echo "Stopping PostgreSQL..."
	@docker stop attendance-postgres || true
	@docker rm attendance-postgres || true