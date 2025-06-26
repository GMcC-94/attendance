include .env
export

# Environment Variables
DB_URL=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
GOOSE=goose -dir ./migrations postgres "$(DB_URL)"

# Targets
.PHONY: run migrate-up migrate-down migrate-status build docker-up docker-down docker-logs

## Run the Go server
run:
	go run main.go

## Build the Go server
build:
	go build -o server main.go

## Run migrations up
migrate-up:
	$(GOOSE) up

## Rollback last migration
migrate-down:
	$(GOOSE) down 1

## Show migration status
migrate-status:
	$(GOOSE) status

## Fix migration files
migrate-fix:
	$(GOOSE) fix

## Create a new migration (Usage: make new-migration name=create_students_table)
new-migration:
	@if [ -z "$(name)" ]; then \
		echo "Migration name required. Usage: make new-migration name=add_table"; \
	else \
		goose -dir ./migrations create $(name) sql; \
	fi

## Start Docker containers
docker-up:
	docker-compose up -d

## Stop Docker containers
docker-down:
	docker-compose down

## View Docker logs
docker-logs:
	docker-compose logs -f