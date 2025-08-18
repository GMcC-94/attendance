# Default target
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development setup
install-deps: ## Install all development dependencies
	@echo "Installing development dependencies..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golangci-lint/golangci-lint/cmd/golangci-lint@latest
	@npm install -g @tailwindcss/cli@latest || echo "Optional: install tailwindcss for custom builds"

setup: install-deps ## Complete project setup
	@echo "Setting up project..."
	@cp .env.example .env
	@echo "Please edit .env file with your configuration"
	@$(MAKE) templ-generate
	@echo "Setup complete! Run 'make dev' to start development server"

# Template generation
templ-generate: ## Generate templ templates
	@echo "Generating templ templates..."
	@templ generate

templ-watch: ## Watch and regenerate templ templates
	@echo "Watching templ files for changes..."
	@templ generate --watch

# Build commands
build: templ-generate ## Build the application binary
	@echo "Building application..."
	@go build -ldflags="-s -w" -o bin/$(BINARY_NAME) cmd/server/main.go

build-race: templ-generate ## Build with race detection (for testing)
	@echo "Building with race detection..."
	@go build -race -o bin/$(BINARY_NAME)-race cmd/server/main.go

# Run commands
run: build ## Run the built application
	@echo "Starting server..."
	@./bin/$(BINARY_NAME)

dev: ## Start development server with hot reload
	@echo "Starting development server with hot reload..."
	@air

dev-full: ## Start full development environment (templates + server)
	@echo "Starting full development environment..."
	@make -j2 templ-watch dev

# Database commands
migrate-create: ## Create a new migration file (usage: make migrate-create NAME=migration_name)
	@echo "Creating migration: $(NAME)"
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(NAME)

migrate-up: ## Run database migrations up
	@echo "Running database migrations up..."
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

migrate-down: ## Run database migrations down
	@echo "Running database migrations down..."
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down

migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	@echo "Forcing migration to version $(VERSION)..."
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" force $(VERSION)

migrate-version: ## Show current migration version
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" version

# Testing
test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	@go test -race -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test -v ./tests/integration/...

# Code quality
fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@templ fmt web/templates/

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run
	@templ fmt --check web/templates/

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):latest .

docker-build-dev: ## Build development Docker image
	@echo "Building development Docker image..."
	@docker build -f docker/Dockerfile.dev -t $(DOCKER_IMAGE):dev .

docker-run: ## Run application with Docker Compose
	@echo "Starting application with Docker Compose..."
	@docker-compose up --build

docker-run-detached: ## Run application with Docker Compose in background
	@echo "Starting application with Docker Compose (detached)..."
	@docker-compose up -d --build

docker-down: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	@docker-compose down

docker-logs: ## Show Docker Compose logs
	@docker-compose logs -f

# Database with Docker
db-start: ## Start PostgreSQL with Docker
	@echo "Starting PostgreSQL database..."
	@docker run --name attendance-postgres -e POSTGRES_DB=attendance_db -e POSTGRES_USER=attendance -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15-alpine

db-stop: ## Stop PostgreSQL Docker container
	@echo "Stopping PostgreSQL database..."
	@docker stop attendance-postgres || true
	@docker rm attendance-postgres || true

db-reset: db-stop db-start ## Reset database (stop, remove, and start fresh)
	@echo "Waiting for database to be ready..."
	@sleep 5
	@$(MAKE) migrate-up

# Cleaning
clean: ## Clean build artifacts and generated files
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f web/templates/*_templ.go
	@docker system prune -f

clean-all: clean ## Clean everything including Docker images
	@echo "Cleaning Docker images..."
	@docker rmi $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):dev 2>/dev/null || true

# Production helpers
build-prod: ## Build for production with optimizations
	@echo "Building for production..."
	@$(MAKE) templ-generate
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=$(shell git describe --tags --always)" -o bin/$(BINARY_NAME) cmd/server/main.go

deploy-prep: build-prod ## Prepare for deployment
	@echo "Preparing for deployment..."
	@$(MAKE) test
	@$(MAKE) lint
	@echo "Deployment preparation complete"

# Development utilities
watch-logs: ## Watch application logs
	@tail -f logs/app.log 2>/dev/null || echo "No log file found"

serve-docs: ## Serve documentation locally
	@echo "Serving documentation on http://localhost:8081"
	@cd docs && python3 -m http.server 8081 2>/dev/null || python -m SimpleHTTPServer 8081

# Security scanning
security-scan: ## Run security vulnerability scan
	@echo "Running security scan..."
	@go install github.com/securecodewarrior/govulncheck@latest
	@govulncheck ./...

# Performance profiling
profile-cpu: ## Run CPU profiling
	@echo "Starting server with CPU profiling..."
	@go run -cpuprofile=cpu.prof cmd/server/main.go

profile-mem: ## Run memory profiling
	@echo "Starting server with memory profiling..."
	@go run -memprofile=mem.prof cmd/server/main.go

# Git hooks
install-git-hooks: ## Install git pre-commit hooks
	@echo "Installing git hooks..."
	@echo "#!/bin/sh\nmake fmt lint test" > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed"

# Environment management
env-local: ## Set up local environment
	@echo "Setting up local environment..."
	@cp .env.example .env.local
	@echo "Edit .env.local for local development"

env-prod: ## Validate production environment
	@echo "Validating production environment..."
	@go run cmd/server/main.go -check-config

# Backup and restore
backup-db: ## Backup database
	@echo "Creating database backup..."
	@pg_dump $(DATABASE_URL) > backup_$(shell date +%Y%m%d_%H%M%S).sql

# Statistics
stats: ## Show project statistics
	@echo "Project Statistics:"
	@echo "===================="
	@echo "Go files: $(shell find . -name '*.go' | wc -l)"
	@echo "Templ files: $(shell find . -name '*.templ' | wc -l)"
	@echo "Total lines of Go code: $(shell find . -name '*.go' -exec cat {} \; | wc -l)"
	@echo "Total lines of Templ code: $(shell find . -name '*.templ' -exec cat {} \; | wc -l)"
	@echo "Dependencies: $(shell go list -m all | wc -l)"

# Quick commands for daily development
quick-start: install-deps setup db-start migrate-up dev ## Complete setup and start development

restart: ## Quick restart of development server
	@pkill -f "$(BINARY_NAME)" || true
	@$(MAKE) dev

update-deps: ## Update all dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@$(MAKE) templ-generate