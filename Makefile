# Baby Tracker Makefile

# Variables
BINARY_NAME=baby-tracker
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/baby-tracker/main.go
WEB_DIR=web

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build variables
VERSION?=dev
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help message
	@echo "Baby Tracker Makefile"
	@echo "====================="
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development targets
.PHONY: run
run: ## Run the server
	$(GOCMD) run $(MAIN_PATH) serve

.PHONY: dev
dev: ## Run with hot reload (requires air)
	@if ! command -v air > /dev/null; then \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

.PHONY: dev-frontend
dev-frontend: ## Run frontend development server
	cd $(WEB_DIR) && npm run dev

# Build targets
.PHONY: build
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_PATH)"

.PHONY: build-frontend
build-frontend: ## Build frontend for production
	@echo "Building frontend..."
	cd $(WEB_DIR) && npm run build

.PHONY: build-all
build-all: build-frontend build ## Build both frontend and backend

.PHONY: build-docker
build-docker: ## Build Docker image
	docker build -t $(BINARY_NAME):latest .

# Database targets
.PHONY: db-test
db-test: ## Test database connection
	$(GOCMD) run $(MAIN_PATH) db test

.PHONY: migrate
migrate: ## Run database migrations
	$(GOCMD) run $(MAIN_PATH) db migrate

.PHONY: migrate-down
migrate-down: ## Rollback database migrations
	$(GOCMD) run $(MAIN_PATH) db migrate down

.PHONY: db-reset
db-reset: ## Reset database (drop all tables and re-migrate)
	@echo "WARNING: This will delete all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		$(GOCMD) run $(MAIN_PATH) db reset; \
	fi

# User management targets
.PHONY: create-user
create-user: ## Create a new user (use with ARGS="-u username -b 'Baby Name'")
	$(GOCMD) run $(MAIN_PATH) create-user $(ARGS)

# Testing targets
.PHONY: test
test: ## Run all tests
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-integration
test-integration: ## Run integration tests
	$(GOTEST) -v -tags=integration ./...

# Dependency management
.PHONY: deps
deps: deps-backend deps-frontend ## Install all dependencies

.PHONY: deps-backend
deps-backend: ## Install backend dependencies
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: deps-frontend
deps-frontend: ## Install frontend dependencies
	cd $(WEB_DIR) && npm install

.PHONY: deps-update
deps-update: ## Update all dependencies
	$(GOGET) -u ./...
	$(GOMOD) tidy
	cd $(WEB_DIR) && npm update

# Linting and formatting
.PHONY: lint
lint: lint-backend lint-frontend ## Run all linters

.PHONY: lint-backend
lint-backend: ## Run Go linter
	@if ! command -v golangci-lint > /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin; \
	fi
	golangci-lint run

.PHONY: lint-frontend
lint-frontend: ## Run frontend linter
	cd $(WEB_DIR) && npm run lint

.PHONY: fmt
fmt: ## Format Go code
	$(GOCMD) fmt ./...

# Cleaning targets
.PHONY: clean
clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -rf $(BINARY_PATH)
	rm -rf tmp/
	rm -rf coverage.out coverage.html

.PHONY: clean-all
clean-all: clean ## Clean everything including dependencies
	rm -rf vendor/
	rm -rf $(WEB_DIR)/node_modules
	rm -rf $(WEB_DIR)/dist

# Docker targets
.PHONY: docker-up
docker-up: ## Start Docker containers
	docker-compose up -d

.PHONY: docker-down
docker-down: ## Stop Docker containers
	docker-compose down

.PHONY: docker-logs
docker-logs: ## View Docker logs
	docker-compose logs -f

.PHONY: docker-reset
docker-reset: docker-down ## Reset Docker environment
	docker-compose down -v
	docker-compose up -d

# Utility targets
.PHONY: setup
setup: deps ## Initial project setup
	@echo "Setting up Baby Tracker..."
	@if [ ! -f .env ]; then \
		echo "Creating .env file..."; \
		cp .env.example .env; \
		echo "Please edit .env with your configuration"; \
	fi
	@echo "Running initial migrations..."
	$(GOCMD) run $(MAIN_PATH) db migrate
	@echo "Setup complete!"

.PHONY: check
check: lint test ## Run all checks (lint + test)

.PHONY: watch-backend
watch-backend: ## Watch backend files for changes (requires entr)
	@if ! command -v entr > /dev/null; then \
		echo "Please install entr: https://github.com/eradman/entr"; \
		exit 1; \
	fi
	find . -name "*.go" -not -path "./vendor/*" -not -path "./web/*" | entr -r make run

.PHONY: generate
generate: ## Run go generate
	$(GOCMD) generate ./...

.PHONY: mod-verify
mod-verify: ## Verify dependencies
	$(GOMOD) verify

.PHONY: mod-vendor
mod-vendor: ## Create vendor directory
	$(GOMOD) vendor

# Release targets
.PHONY: release
release: clean build-all ## Build release version
	@echo "Building release version $(VERSION)..."
	mkdir -p releases
	tar -czf releases/$(BINARY_NAME)-$(VERSION).tar.gz $(BINARY_PATH) README.md

.PHONY: version
version: ## Display version information
	@if [ -f $(BINARY_PATH) ]; then \
		./$(BINARY_PATH) version; \
	else \
		echo "Binary not built. Run 'make build' first."; \
	fi

# Development database shortcuts
.PHONY: db-shell
db-shell: ## Open database shell
	@if [ "$(DB_TYPE)" = "postgres" ]; then \
		psql -h localhost -U postgres -d baby_tracker; \
	else \
		sqlite3 baby-tracker.db; \
	fi

# Quick commands for development
.PHONY: quick-start
quick-start: deps setup dev ## Quick start for development

.PHONY: seed
seed: ## Seed database with test data
	$(GOCMD) run $(MAIN_PATH) db seed

# CI/CD targets
.PHONY: ci
ci: deps lint test build ## Run CI pipeline

.PHONY: cd
cd: build-all ## Run CD pipeline
	@echo "Deploying application..."
	# Add deployment steps here
