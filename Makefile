.PHONY: help migration run_terndotenv generate_sqlc



MIGRATIONS_PATH=internal/store/pgstore/migrations

## Display this help message
help:
	@echo ""
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

## Run tern with dotenv loaded and create a new migration - usage: make migration name=create_devices_table
migrate_generate: run_terndotenv ## Create a new migration with tern
	cd $(MIGRATIONS_PATH) && tern new $(name)

## Run migrations with tern
migrate_run: ## Run migrations
	go run ./cmd/terndotenv

## Rollback migrations with tern - usage: make migrate_rollback steps=1
migrate_rollback: ## Rollback migrations
	go run ./cmd/terndotenv steps=$(steps)

## Generate code using sqlc
queries_run: ## Run sqlc generate
	cd internal/store/pgstore/ && sqlc generate -f ./sqlc.yaml

## Run tests
tests: ## Run tests
	go test ./...

## Run Tests with coverage
tests_coverage: ## Run tests with coverage
	go test -cover ./...

## Run application
run: ## Run application
	air -c .air.toml

## Run docker compose
up: ## Run docker compose
	docker-compose up -d

down: ## Stop docker compose
	docker-compose down -v

CONTAINER_NAME=devices-api-dev
EXEC=docker exec -it $(CONTAINER_NAME)
## Run migrations
migrate: ## Run migrations
	$(EXEC) go run ./cmd/terndotenv

## Rebuild the container
rebuild: ## Rebuild the container
	docker-compose down -v
	docker-compose up -d --build