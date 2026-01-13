.PHONY: help dev setup verify migrate migrate-create sqlc-generate clean test build docker-build

# Default target
help:
	@echo "Superfly Development Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make setup           - Run dev environment setup"
	@echo "  make verify          - Verify environment is ready"
	@echo ""
	@echo "Development:"
	@echo "  make dev             - Run API server with live reload"
	@echo "  make dev-web         - Run Svelte frontend"
	@echo "  make sqlc-generate   - Generate Go code from SQL"
	@echo "  make migrate         - Run database migrations"
	@echo "  make migrate-create  - Create a new migration (NAME=migration_name)"
	@echo ""
	@echo "Database:"
	@echo "  make db-shell        - Connect to PostgreSQL"
	@echo "  make db-reset        - Drop and recreate database"
	@echo ""
	@echo "Build:"
	@echo "  make build           - Build API server binary"
	@echo "  make build-web       - Build frontend"
	@echo "  make docker-build    - Build Docker images"
	@echo ""
	@echo "Testing:"
	@echo "  make test            - Run tests"
	@echo "  make test-coverage   - Run tests with coverage"
	@echo ""
	@echo "Clean:"
	@echo "  make clean           - Clean build artifacts"

# Variables
DATABASE_URL ?= postgresql://superfly:superfly_dev_password@localhost:5432/superfly?sslmode=disable
MIGRATION_DIR = db/migrations
GOOSE = goose -dir $(MIGRATION_DIR) postgres "$(DATABASE_URL)"

# Setup development environment
setup:
	@chmod +x dev-setup.sh
	@./dev-setup.sh

# Verify environment
verify:
	@chmod +x verify-setup.sh
	@./verify-setup.sh

# Run API server with live reload (requires air)
dev:
	@echo "Starting API server with live reload..."
	@cd cmd/api && air

# Run frontend dev server
dev-web:
	@echo "Starting Svelte dev server..."
	@cd web && npm run dev

# Generate Go code from SQL queries using sqlc
sqlc-generate:
	@echo "Generating Go code from SQL..."
	@sqlc generate

# Run database migrations
migrate:
	@echo "Running database migrations..."
	@$(GOOSE) up

# Create a new migration
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=your_migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	@$(GOOSE) create $(NAME) sql

# Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	@$(GOOSE) down

# Reset database (drop all tables and rerun migrations)
db-reset:
	@echo "Resetting database..."
	@$(GOOSE) reset
	@$(GOOSE) up

# Connect to PostgreSQL shell
db-shell:
	@PGPASSWORD=superfly_dev_password psql -h localhost -U superfly -d superfly

# Build API server binary
build:
	@echo "Building API server..."
	@go build -o bin/superfly-api ./cmd/api

# Build frontend
build-web:
	@echo "Building frontend..."
	@cd web && npm run build

# Build Docker images
docker-build:
	@echo "Building API Docker image..."
	@docker build -f Dockerfile.api -t superfly/api:latest .
	@echo "Building Web Docker image..."
	@docker build -f Dockerfile.web -t superfly/web:latest ./web

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf tmp/
	@rm -f coverage.out coverage.html
	@cd web && rm -rf build/ .svelte-kit/

# Install development dependencies
install-deps:
	@echo "Installing Go dependencies..."
	@go mod download
	@echo "Installing frontend dependencies..."
	@cd web && npm install

# Format code
fmt:
	@echo "Formatting Go code..."
	@go fmt ./...
	@echo "Formatting frontend code..."
	@cd web && npm run format || true

# Lint code
lint:
	@echo "Linting Go code..."
	@golangci-lint run ./... || echo "golangci-lint not installed, skipping..."
	@echo "Linting frontend code..."
	@cd web && npm run lint || true

# Show database migration status
migrate-status:
	@$(GOOSE) status

# Port forward PostgreSQL (manual, useful if systemd service isn't running)
port-forward-db:
	@echo "Port forwarding PostgreSQL..."
	@kubectl port-forward -n superfly-system svc/postgres 5432:5432

# View logs of system pods
logs-system:
	@echo "Recent logs from superfly-system pods:"
	@kubectl logs -n superfly-system deployment/postgres --tail=50

# View logs of app pods
logs-apps:
	@echo "Recent logs from superfly-apps:"
	@kubectl get pods -n superfly-apps
	@echo ""
	@echo "Use: kubectl logs -n superfly-apps <pod-name>"

# Initialize Go module (first time setup)
init:
	@if [ ! -f go.mod ]; then \
		echo "Initializing Go module..."; \
		go mod init github.com/yourusername/superfly; \
	else \
		echo "Go module already initialized"; \
	fi
	@echo "Installing Go dependencies..."
	@go get github.com/jackc/pgx/v5
	@go get github.com/jackc/pgx/v5/pgxpool
	@go get github.com/google/uuid
	@go get k8s.io/client-go@latest
	@go get k8s.io/api@latest
	@go get k8s.io/apimachinery@latest
	@go mod tidy
