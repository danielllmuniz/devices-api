.PHONY: help migration run_terndotenv

MIGRATIONS_PATH=internal/store/pgstore/migrations

## Display this help message
help:
	@echo ""
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

## Run tern with dotenv loaded and create a new migration - usage: make migration name=create_devices_table
migration: run_terndotenv ## Create a new migration with tern
	cd $(MIGRATIONS_PATH) && tern new $(name)

## Run the command go run ./cmd/terndotenv
run_terndotenv: ## Load environment variables using dotenv
	go run ./cmd/terndotenv
