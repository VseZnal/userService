DATABASE_CONNECTION_STRING=postgresql://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DB}?sslmode=disable

MIGRATE := docker run --network host -v $(shell pwd)/database/migrations/:/migrations/ \
migrate/migrate:v4.10.0 -path=/migrations/ -database "$(DATABASE_CONNECTION_STRING)"

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# ===================== Database management ========================================

.PHONY: migrate
migrate: ## run all snew database migrations
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## revert database to the last migration step
	@echo "Reverting database to the last migration step..."
	@$(MIGRATE) down 1

.PHONY: migrate-new
migrate-new: ## create a new database migration
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE) create -ext sql -dir backend/migrations/ $${name// /_}

.PHONY: migrate-reset
migrate-reset: ## reset database and re-run all migrations
	@echo "Resetting database..."
	echo "y" | $(MIGRATE) down -all
	@echo "Running all database migrations..."
	@$(MIGRATE) up

.PHONY: feed
feed: ## insert test values into a database
	cat ./database/scripts/feed.sql | docker exec -i postgres psql "postgresql://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:5432/${PG_DB}?sslmode=disable"
