default:
    just --list

# 🚀 Development

# Start both backend and frontend development servers
dev:
    #!/usr/bin/env bash
    echo "🚀 Starting development servers..."
    echo "Backend: http://localhost:8080"
    echo "Frontend: http://localhost:5173"
    trap 'kill 0' EXIT
    just dev-backend &
    just dev-frontend &
    wait

# Start backend development server
dev-backend:
    go run cmd/baby-tracker/main.go serve

# Start frontend development server  
dev-frontend:
    cd web && npm run dev

# 🏗️ Building

# Build frontend and backend for production
build:
    #!/usr/bin/env bash
    echo "🏗️ Building frontend..."
    just build-frontend
    echo "🏗️ Building backend..."
    just build-backend
    echo "✅ Build complete!"

# Build frontend only
build-frontend:
    cd web && npm run build

# Build backend only
build-backend:
    go build -o bin/baby-tracker cmd/baby-tracker/main.go

# 📦 Dependencies

# Install all dependencies (backend + frontend)
deps:
    #!/usr/bin/env bash
    echo "📦 Installing Go dependencies..."
    go mod download
    go mod tidy
    echo "📦 Installing frontend dependencies..."
    cd web && npm install

# Update all dependencies
deps-update:
    #!/usr/bin/env bash
    echo "📦 Updating Go dependencies..."
    go get -u ./...
    go mod tidy
    echo "📦 Updating frontend dependencies..."
    cd web && npm update

# 🗄️ Database

# Run database migrations
migrate:
    go run cmd/baby-tracker/main.go db migrate

# Rollback database migrations
migrate-down:
    go run cmd/baby-tracker/main.go db migrate down

# Reset database (⚠️ deletes all data!)
db-reset:
    #!/usr/bin/env bash
    echo "⚠️  WARNING: This will delete ALL data!"
    read -p "Are you sure? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        go run cmd/baby-tracker/main.go db reset
        echo "✅ Database reset complete"
    else
        echo "❌ Database reset cancelled"
    fi

# Create a new user (requires -u username -b "Baby Name")
create-user *ARGS:
    go run cmd/baby-tracker/main.go create-user {{ARGS}}

# Open database shell
db-shell:
    #!/usr/bin/env bash
    if [ -f baby-tracker.db ]; then
        sqlite3 baby-tracker.db
    else
        echo "Database file not found. Run 'just migrate' first."
    fi

# 🧪 Testing

# Run all tests
test:
    go test -v ./...

# Run tests with coverage
test-coverage:
    #!/usr/bin/env bash
    echo "📊 Running tests with coverage..."
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    echo "📈 Coverage report: coverage.html"

# Run tests with race detection
test-race:
    go test -v -race ./...

# 🧹 Cleanup

# Clean build artifacts
clean:
    #!/usr/bin/env bash
    echo "🧹 Cleaning build artifacts..."
    go clean
    rm -rf bin/
    rm -rf web/dist/
    rm -rf tmp/

# Clean everything including dependencies
clean-all:
    #!/usr/bin/env bash
    echo "🧹 Cleaning everything..."
    just clean
    rm -rf web/node_modules/
    rm -rf vendor/

# Clean test artifacts
test-clean:
    rm -f coverage.out coverage.html
    go clean -testcache

# � Utilities

# Run linting on backend code
lint:
    #!/usr/bin/env bash
    if ! command -v golangci-lint > /dev/null; then
        echo "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    fi
    golangci-lint run

# Format Go code
fmt:
    go fmt ./...

# Run all checks (lint + test)
check: lint test

# Project setup for first time
setup:
    #!/usr/bin/env bash
    echo "🔧 Setting up Baby Tracker..."
    just deps
    if [ ! -f .env ]; then
        echo "📝 Creating .env file..."
        if [ -f .env.example ]; then
            cp .env.example .env
            echo "⚠️  Please edit .env with your configuration"
        else
            echo "SESSION_SECRET=your-secret-key-here" > .env
            echo "DB_TYPE=sqlite" >> .env
            echo "PORT=8080" >> .env
        fi
    fi
    echo "🗄️ Running initial migrations..."
    just migrate
    echo "✅ Setup complete!"
