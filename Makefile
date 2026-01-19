# Check if the environment file exists
ENVFILE := .env
ifneq ("$(wildcard $(ENVFILE))","")
	include $(ENVFILE)
	export $(shell sed 's/=.*//' $(ENVFILE))
endif

# Outlet Makefile
EXECUTABLE=outlet

.PHONY: help build build-all build-api build-worker build-mcp run migrate-up migrate-down migrate-status migrate-reset clean test swagger

# Default target
help:
	@echo "Outlet - AI-Powered Brand Platform"
	@echo ""
	@echo "Available targets:"
	@echo "  build          - Build the main application binary (go-zero API)"
	@echo "  build-all      - Build the all-in-one binary (API + Worker + MCP)"
	@echo "  build-api      - Build API server only"
	@echo "  build-worker   - Build Temporal worker only"
	@echo "  build-mcp      - Build MCP server only"
	@echo "  run            - Run the main application"
	@echo "  run-all        - Run the all-in-one server"
	@echo "  migrate-up     - Run all pending database migrations"
	@echo "  migrate-down   - Rollback one migration"
	@echo "  migrate-status - Show migration status"
	@echo "  migrate-reset  - Reset database (down all, then up all)"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  deps           - Download dependencies"
	@echo "  gen            - Generate code (go-zero)"
	@echo "  gen-sdk        - Generate SDK clients (TypeScript, Python, Go, PHP)"
	@echo "  models         - Generate PostgreSQL models for all tables"
	@echo "  swagger        - Validate and display Swagger documentation info"
	@echo "  install-updater - Install go-selfupdate CLI tool"
	@echo "  create-update  - Create update files for current version"
	@echo "  create-update-dir - Create update files for cross-compiled binaries"

# Version info for ldflags
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -X outlet/internal/version.Version=$(VERSION) -X outlet/internal/version.Commit=$(COMMIT) -X outlet/internal/version.BuildDate=$(BUILD_DATE)

# Build the main application (go-zero)
build:
	@echo "Building Outlet $(VERSION)..."
	go build -ldflags "$(LDFLAGS)" -o bin/outlet .

# Build all-in-one binary (API + Worker + MCP)
build-all:
	@echo "Building Outlet All-in-One (API + Worker + MCP)..."
	go build -o bin/all ./cmd/all

# Build API server only
build-api:
	@echo "Building Outlet API server..."
	go build -o bin/api ./cmd/api

# Build Temporal worker only
build-worker:
	@echo "Building Outlet Worker..."
	go build -o bin/worker ./cmd/worker

# Build MCP server only
build-mcp:
	@echo "Building Outlet MCP server..."
	go build -o bin/mcp ./cmd/mcp

# Run the main application
run: build
	@echo "Starting Outlet..."
	./bin/outlet

# Run all-in-one server
run-all: build-all
	@echo "Starting Outlet All-in-One Server..."
	./bin/all

# Database migration commands using psql
migrate-up:
	@echo "Running database migrations..."
	@if [ -z "$$POSTGRES_EXTERNAL_DSN" ]; then \
		echo "Error: POSTGRES_EXTERNAL_DSN not set"; \
		exit 1; \
	fi
	@echo "Applying analytics system migration..."
	@psql "$$POSTGRES_EXTERNAL_DSN" -f ./internal/db/migrations/add_analytics_system.sql
	@echo "Applying dashboard publishing migration..."
	@psql "$$POSTGRES_EXTERNAL_DSN" -f ./internal/db/migrations/add_dashboard_publishing_and_layouts.sql
	@echo "✓ All migrations applied successfully"

migrate-schema:
	@echo "Applying base schema..."
	@if [ -z "$$POSTGRES_EXTERNAL_DSN" ]; then \
		echo "Error: POSTGRES_EXTERNAL_DSN not set"; \
		exit 1; \
	fi
	@psql "$$POSTGRES_EXTERNAL_DSN" -f ./internal/db/schema.sql
	@echo "✓ Base schema applied"

migrate-status:
	@echo "Checking migration tables..."
	@if [ -z "$$POSTGRES_EXTERNAL_DSN" ]; then \
		echo "Error: POSTGRES_EXTERNAL_DSN not set"; \
		exit 1; \
	fi
	@psql "$$POSTGRES_EXTERNAL_DSN" -c "\dt analytics_*"
	@psql "$$POSTGRES_EXTERNAL_DSN" -c "\dt dashboard_*"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Code generation
gen: ## Generate API code from .api file
	@echo "🧹 Cleaning auto-generated handlers..."
	@rm -rf internal/handler
	@echo "Generating Go API code..."
	goctl api go -api $(EXECUTABLE).api -dir . --style gozero
	@echo "Generating TypeScript API code..."
	goctl api ts -api $(EXECUTABLE).api -dir ./app/src/lib/api/generate
	@if [ -f cmd/sourcegen/main.go ]; then \
		echo "Generating source registry..."; \
		go run cmd/sourcegen/main.go internal/sources internal/platform/registry.go; \
	fi
	@echo "✅ API code generation complete!"

# SDK generation
gen-sdk: ## Generate SDK clients from .api file (4 languages)
	@echo "📦 Generating SDK clients..."
	@echo "Generating TypeScript SDK..."
	goctl api plugin -plugin "goctl-sdk -format typescript -prefix /sdk/v1 -name Outlet" \
		-api $(EXECUTABLE).api -dir ./sdk/typescript/src
	@echo "Generating Python SDK..."
	goctl api plugin -plugin "goctl-sdk -format python -prefix /sdk/v1 -name Outlet" \
		-api $(EXECUTABLE).api -dir ./sdk/python/src/outlet_sdk
	@echo "Generating Go SDK..."
	goctl api plugin -plugin "goctl-sdk -format golang -prefix /sdk/v1 -name Outlet" \
		-api $(EXECUTABLE).api -dir ./sdk/go
	@echo "Generating PHP SDK..."
	goctl api plugin -plugin "goctl-sdk -format php -prefix /sdk/v1 -name Outlet" \
		-api $(EXECUTABLE).api -dir ./sdk/php/src
	@echo "✅ SDK generation complete (4 languages)!"

# Database code generation
sqlc-gen: ## Generate type-safe Go code from SQL queries
	@echo "Generating sqlc code..."
	sqlc generate
	@echo "✓ sqlc code generation complete!"

# Outlet database migration commands
db-create: ## Create the outlet database
	@echo "Creating outlet database..."
	@if [ -z "$$POSTGRES_HOST" ]; then POSTGRES_HOST=localhost; fi
	@if [ -z "$$POSTGRES_PORT" ]; then POSTGRES_PORT=5432; fi
	@if [ -z "$$POSTGRES_USER" ]; then POSTGRES_USER=postgres; fi
	@psql -h $${POSTGRES_HOST:-localhost} -p $${POSTGRES_PORT:-5432} -U $${POSTGRES_USER:-postgres} -c "CREATE DATABASE outlet;" || true
	@psql -h $${POSTGRES_HOST:-localhost} -p $${POSTGRES_PORT:-5432} -U $${POSTGRES_USER:-postgres} -c "CREATE DATABASE outlet_test;" || true
	@echo "✓ Databases created"

db-migrate: ## Run outlet schema migrations
	@echo "Running outlet schema migrations..."
	@psql -h $${POSTGRES_HOST:-localhost} -p $${POSTGRES_PORT:-5432} -U $${POSTGRES_USER:-postgres} -d outlet -f ./internal/db/migrations/000001_initial_schema.up.sql
	@echo "✓ Schema migrations applied"

db-rollback: ## Rollback outlet schema migrations
	@echo "Rolling back outlet schema migrations..."
	@psql -h $${POSTGRES_HOST:-localhost} -p $${POSTGRES_PORT:-5432} -U $${POSTGRES_USER:-postgres} -d outlet -f ./internal/db/migrations/000001_initial_schema.down.sql
	@echo "✓ Schema migrations rolled back"

db-seed: ## Seed the outlet database with test data
	@echo "Seeding outlet database..."
	@psql -h $${POSTGRES_HOST:-localhost} -p $${POSTGRES_PORT:-5432} -U $${POSTGRES_USER:-postgres} -d outlet -f ./internal/db/seeds/001_seed_data.sql
	@echo "✓ Database seeded"

db-reset: db-rollback db-migrate db-seed ## Reset outlet database (rollback, migrate, seed)
	@echo "✓ Database reset complete"

db-test-migrate: ## Run migrations on test database
	@echo "Running migrations on test database..."
	@psql -h $${POSTGRES_HOST:-localhost} -p $${POSTGRES_PORT:-5432} -U $${POSTGRES_USER:-postgres} -d outlet_test -f ./internal/db/migrations/000001_initial_schema.up.sql
	@echo "✓ Test database migrations applied"

test-db: ## Run database tests
	@echo "Running database tests..."
	go test -v ./internal/db/sqlc/...
	@echo "✓ Database tests complete"

# Proto generation
proto-gen: ## Generate gRPC code from proto files
	@echo "Generating gRPC code from proto files..."
	@mkdir -p internal/grpc/llmpb
	protoc --go_out=internal/grpc/llmpb --go_opt=paths=source_relative \
		--go-grpc_out=internal/grpc/llmpb --go-grpc_opt=paths=source_relative \
		-Iproto proto/llm.proto
	@echo "✓ Proto code generation complete!"

# Full generation (API + sqlc + proto)
gen-all: gen sqlc-gen proto-gen
	@echo "✓ All code generation complete!"


# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@mkdir -p bin
	@echo "Development setup complete!"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t outlet .

docker-run:
	@echo "Running Docker container..."
	docker run -p 9888:9888 --env-file .env outlet

# Database setup for development
db-setup:
	@echo "Setting up development database..."
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Database setup complete!"

# Full development environment
dev: dev-setup db-setup migrate-up
	@echo "Development environment ready!"
	@echo "Run 'make run' to start the application"

# Self-update commands
.PHONY: install-updater
install-updater: ## Install go-selfupdate CLI tool
	@echo "Installing go-selfupdate CLI tool..."
	go install github.com/sanbornm/go-selfupdate/cmd/go-selfupdate@latest

.PHONY: create-update
create-update: ## Create update files for current version
	@echo "Creating update files..."
	@if [ -z "$(VERSION)" ]; then echo "Usage: make create-update VERSION=1.0.1"; exit 1; fi
	go-selfupdate $(EXECUTABLE) $(VERSION)

.PHONY: create-update-dir
create-update-dir: ## Create update files for cross-compiled binaries in directory
	@echo "Creating update files for cross-compiled binaries..."
	@if [ -z "$(VERSION)" ] || [ -z "$(DIR)" ]; then echo "Usage: make create-update-dir VERSION=1.0.1 DIR=/path/to/binaries"; exit 1; fi
	go-selfupdate $(DIR) $(VERSION)
	