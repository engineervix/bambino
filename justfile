export PATH := "./node_modules/.bin:" + env_var('PATH')
DOCKER_COMPOSE_DEV := "docker compose -f docker/docker-compose.dev.yml --env-file .env"
DOCKER_COMPOSE_PROD := "docker compose --env-file .prod.env"

default:
    just --list

# docker-compose up
up *build:
    if [[ "{{build}}" =~ ^(-b|b|build|--build)$ ]]; then \
        {{DOCKER_COMPOSE_DEV}} up -d --build; \
    elif [[ "{{build}}" = "" ]]; then \
        {{DOCKER_COMPOSE_DEV}} up -d; \
    else \
        echo "{{build}} doesn't match any of -b, b, build or --build"; \
    fi

# docker-compose exec [container] [command(s)]
exec container +command:
    {{DOCKER_COMPOSE_DEV}} exec {{container}} "{{command}}"

# docker-compose logs [container] [-f (Follow log output)]
logs container *follow:
    if [[ "{{follow}}" =~ ^(-f|f|follow|--follow)$ ]]; then \
        {{DOCKER_COMPOSE_DEV}} logs {{container}} -f; \
    elif [[ "{{follow}}" = "" ]]; then \
        {{DOCKER_COMPOSE_DEV}} logs {{container}}; \
    else \
        echo "{{follow}} doesn't match any of -f, f, follow or --follow"; \
    fi

# docker-compose stop
stop:
    docker-compose stop

# docker-compose down [-v]
down *volumes:
    if [[ "{{volumes}}" =~ ^(-v|v|--vol|vol|volumes|--volumes)$ ]]; then \
        {{DOCKER_COMPOSE_DEV}} down -v; \
    elif [[ "{{volumes}}" = "" ]]; then \
        {{DOCKER_COMPOSE_DEV}} down; \
    else \
        echo "{{volumes}} doesn't match any of -v, v, vol, --vol, volumes or --volumes"; \
    fi

# 🔨 Build the frontend and the Go application
build:
    @echo "🖼️ Building frontend assets..."
    npm run build
    @echo "🔨 Building Go application binary..."
    CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/bambino ./cmd/bambino/main.go

# Run database migrations
db-migrate:
    @go run cmd/bambino/main.go db migrate

# 🚀 Run the Go application
run:
    @echo "🚀 Starting Go application..."
    go run ./cmd/bambino/main.go serve

# Start all services in production mode
prod-up:
    @echo "Starting production services..."
    @{{DOCKER_COMPOSE_PROD}} up -d

# Stop all services in production mode
prod-down:
    @echo "Stopping production services..."
    @{{DOCKER_COMPOSE_PROD}} down

# View logs for all production services
prod-logs:
    @{{DOCKER_COMPOSE_PROD}} logs -f

# Execute a command inside the running production container
prod-exec +CMD:
    @{{DOCKER_COMPOSE_PROD}} exec bambino {{CMD}}

# 🧪 Run tests with better output formatting
test:
    @echo "🧪 Running tests..."
    gotestsum --format=pkgname-and-test-fails

# 📊 Run tests with coverage analysis
test-coverage:
    @echo "📊 Running tests with coverage analysis..."
    go test -coverprofile=coverage.out ./...
    @echo "📈 Coverage summary:"
    go tool cover -func=coverage.out

# 🧹 Clean test artifacts
test-clean:
    @echo "🧹 Cleaning test artifacts..."
    rm -f coverage.out core_coverage.out
    go clean -testcache
