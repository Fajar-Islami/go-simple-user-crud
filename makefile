.PHONY: help
help: ## Show help command
	@printf "Makefile Command\n";
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: docs 
docs: ## Generate Documents
	swag init -g ./cmd/main.go --output docs


.PHONY: migrate
migrate: ## Create Migrations file
	@if [ -z "${name}" ]; then \
		echo "Error: name is required \t example : make migrate name="name_file_migration";" \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations ':hammer: ${name}'



migrate-up: ## Up migration
	go run cmd/migrate/main.go

migrate-rollback: ## Up rollback
	go run cmd/migrate/main.go -rollback