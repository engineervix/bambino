# Baby Tracker Makefile

# Variables
BINARY_NAME=baby-tracker
GO_MODULE=github.com/yourusername/baby-tracker

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Directories
BACKEND_DIR=.
FRONTEND_DIR=web

.PHONY: all build clean test deps run-backend run-frontend run dev setup help

# Default target
all: test build

# Build the application
build: build-frontend build-backend

build-backend:
	@echo "Building backend..."
	$(GOBUILD) -o $(BINARY_NAME) -v cmd/server/main.go

build-frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm run build

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(FRONTEND_DIR)/dist

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Download dependencies
deps:
	@echo "Downloading Go dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Installing frontend dependencies..."
	cd $(FRONTEND_DIR) && npm install

# Run backend only
run-backend:
	@echo "Starting backend server..."
	$(GOCMD) run cmd/server/main.go

# Run frontend only
run-frontend:
	@echo "Starting frontend development server..."
	cd $(FRONTEND_DIR) && npm run dev

# Run both frontend and backend (requires 2 terminals)
run:
	@echo "Start backend and frontend in separate terminals:"
	@echo "  Terminal 1: make run-backend"
	@echo "  Terminal 2: make run-frontend"

# Development mode - watch for changes
dev:
	@echo "Starting development mode..."
	@echo "Installing air for hot reload..."
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	@echo "Starting backend with hot reload..."
	air

# Database operations
db-test:
	@echo "Testing database connection..."
	$(GOCMD) run scripts/test-db.go

db-reset:
	@echo "Resetting database..."
	RESET_DB=true $(GOCMD) run cmd/server/main.go

create-user:
	@echo "Creating user..."
	@echo "Usage: make create-user ARGS=\"--username=parent --password=pass123\""
	$(GOCMD) run scripts/create-user.go $(ARGS)

# Initial project setup
setup: deps
	@echo "Setting up project..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from template"; \
		echo "Please edit .env with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi
	@echo "Setup complete!"

# Format code
fmt:
	@echo "Formatting Go code..."
	$(GOCMD) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

# Generate code coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Docker operations
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

# Help
help:
	@echo "Baby Tracker - Available commands:"
	@echo ""
	@echo "  make setup          - Initial project setup"
	@echo "  make deps           - Download all dependencies"
	@echo "  make build          - Build both frontend and backend"
	@echo "  make run-backend    - Run backend server"
	@echo "  make run-frontend   - Run frontend dev server"
	@echo "  make dev            - Run backend with hot reload (requires air)"
	@echo "  make test           - Run tests"
	@echo "  make fmt            - Format Go code"
	@echo "  make lint           - Lint code"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "Database commands:"
	@echo "  make db-test        - Test database connection"
	@echo "  make db-reset       - Reset database (WARNING: deletes all data)"
	@echo "  make create-user ARGS=\"--username=X --password=Y\" - Create a user"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run Docker container"