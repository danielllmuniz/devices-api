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
test: ## Run tests
	go test -v ./...