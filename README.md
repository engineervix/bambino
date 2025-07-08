# Baby Tracker

A self-hosted baby activity tracking app for personal use. Tracks feeding, sleeping, diapers, and other baby activities with complete data ownership.

## What it does

- Track feeds, pumps, diapers, sleep, growth, health records, and milestones
- Mobile-first design with dark mode for nighttime use
- Timer functionality for activities
- Single binary deployment with embedded frontend

## Tech stack

- **Backend**: Go with [Echo](https://echo.labstack.com/) framework
- **Frontend**: [Vue 3](https://vuejs.org/) + [Vuetify](https://vuetifyjs.com/)
- **Database**: SQLite (default) or PostgreSQL
- **Auth**: Session-based

## Quick start

### Development
```bash
# Install dependencies
just deps

# Start development servers (backend + frontend)
just dev

# Or start them separately
just dev-backend    # Go server on :8080
just dev-frontend   # Vite dev server on :5173
```

### Production
```bash
# Build everything
just build

# Or build and run in Docker
just docker-up
```

## Setup

1. Install dependencies: `just deps`
2. Copy `.env.example` to `.env` and configure
3. Run migrations: `just migrate`
4. Create first user: `just create-user -u admin -b "Baby Name"`
5. Start development: `just dev`

## Database

- **SQLite** (default): Data stored in `baby-tracker.db`
- **PostgreSQL**: Configure in `.env` with `DB_TYPE=postgres`

### Common tasks
```bash
just migrate         # Run migrations
just migrate-down    # Rollback migrations
just db-reset        # Reset database (deletes all data!)
just create-user     # Create new user
```

## Testing

```bash
just test            # Run all tests
just test-coverage   # Run tests with coverage
just test-race       # Run with race detection
```

## Deployment

Single binary with embedded frontend:
```bash
just build           # Creates bin/baby-tracker
./bin/baby-tracker serve
```
