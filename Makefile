# Variables
PKG := ./...  # Test all packages
GOOSE := $(HOME)/go/bin/goose
MIGRATIONS_DIR := ./migrations
DB_URL := "postgres://root:root@localhost:5432?sslmode=disable&database=league"

# Default target: Show help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make test              - Run all tests"
	@echo "  make color-test        - Run tests with colored output"
	@echo "  make create-migration  - Create a new database migration"
	@echo "  make migrate-up        - Run all pending database migrations"
	@echo "  make migrate-down      - Rollback the last database migration"

# Default target: Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v $(PKG)

# Run tests with colored output using gotestsum
.PHONY: color-test
color-test:
	@echo "Running tests with colorized output..."
	gotestsum --format=short-verbose --color=always

# Create a new migration
.PHONY: create-migration
create-migration:
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $$name sql

# Run migrations
.PHONY: migrate-up
migrate-up:
	@echo "Running migrations..."
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres $(DB_URL) up

# Rollback the last migration
.PHONY: migrate-down
migrate-down:
	@echo "Rolling back last migration..."
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres $(DB_URL) down
