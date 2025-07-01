include .env
export

# Environment Variables
DB_URL=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
GOOSE=goose -dir ./migrations postgres "$(DB_URL)"

# Targets
.PHONY: run migrate-up migrate-down migrate-status build docker-up docker-down docker-logs

## Run the Go server
run:
	@cd backend && go run main.go

## Build the Go server
build:
	go build -o server main.go

## Run migrations up
migrate-up:
	cd backend && $(GOOSE) up

## Rollback last migration
migrate-down:
	cd backend && $(GOOSE) down 1

## Show migration status
migrate-status:
	cd backend && $(GOOSE) status

## Fix migration files
migrate-fix:
	cd backend && $(GOOSE) fix

	## Fix migration files
migrate-reset:
	cd backend && $(GOOSE) reset

## Create a new migration (Usage: make new-migration name=create_students_table)
new-migration:
	@if [ -z "$(name)" ]; then \
		echo "Migration name required. Usage: make new-migration name=add_table"; \
	else \
		cd backend && goose -dir ./migrations create $(name) sql; \
	fi

## Start Docker containers
docker-up:
	cd backend && docker-compose up -d

## Stop Docker containers
docker-down:
	cd backend && docker-compose down

## View Docker logs
docker-logs:
	cd backend && docker-compose logs -f