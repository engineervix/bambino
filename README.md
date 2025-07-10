# Bambino

A self-hosted baby activity tracking app for personal use. Tracks feeding, sleeping, diapers, and other baby activities with complete data ownership.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/engineervix/bambino)
[![CI/CD](https://github.com/engineervix/bambino/actions/workflows/main.yml/badge.svg)](https://github.com/engineervix/bambino/actions/workflows/main.yml)

[![Node v22](https://img.shields.io/badge/Node-v22-teal.svg)](https://nodejs.org/en/blog/release/v22.0.0)
[![code style: prettier](https://img.shields.io/badge/code%20style-prettier-ff69b4.svg)](https://prettier.io/)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Features](#features)
- [Tech stack](#tech-stack)
- [Running the application](#running-the-application)
  - [Prerequisites](#prerequisites)
  - [Configuration](#configuration)
  - [Build](#build)
  - [Database](#database)
  - [Development Mode](#development-mode)
  - [Production Mode](#production-mode)
    - [With Docker (Recommended)](#with-docker-recommended)
    - [With Local Binary](#with-local-binary)
  - [Creating a User](#creating-a-user)
  - [Command-Line Help](#command-line-help)
- [Testing](#testing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## Features

- Track feeds, pumps, diapers, sleep, growth, health records, and milestones
- Mobile-first design with dark mode for nighttime use
- Timer functionality for activities
- Single binary deployment with embedded frontend

## Tech stack

- **Backend**: Go with [Echo](https://echo.labstack.com/) framework
- **Frontend**: [Vue 3](https://vuejs.org/) + [Vuetify](https://vuetifyjs.com/)
- **Database**: SQLite or PostgreSQL
- **Auth**: Session-based

## Running the application

### Prerequisites

- Go
- Node.js and npm
- [just](https://github.com/casey/just)

### Configuration

Copy the example environment files and update them with your configuration.

```bash
cp .env.example .env
cp .env.test.example .env.test
```

### Build

To build the application binary and frontend assets, run:

```bash
just build
```

### Database

The application can use either PostgreSQL or SQLite.

If you are using PostgreSQL, start the container:
```bash
just up
```

Then, run the database migrations. It's recommended to use the binary:
```bash
./bin/bambino db migrate
```

Alternatively, you can run the migrations without building the binary:
```bash
go run cmd/bambino/main.go db migrate
```

### Development Mode

For development, you can run the backend and frontend separately.

**Backend:**
```bash
just run
```

**Frontend:**
```bash
npm run dev
```
The application will be available at `http://localhost:5173`.

### Production Mode

There are two ways to run the application in production mode.

#### With Docker (Recommended)

For production, create a `.prod.env` file with your production secrets. **Do not commit this file to version control.**

Create a file named `.prod.env` and paste the following content into it, replacing the placeholder values:

```env
# The domain name for your application
DOMAIN_NAME=your-domain.com

# PostgreSQL connection details
POSTGRES_USER=baby
POSTGRES_DB=baby
POSTGRES_PASSWORD=generate-a-strong-password

# A long, random string for session signing
SESSION_SECRET=generate-a-long-random-secret-string

# Backblaze B2 credentials (optional, for backups)
# B2_APPLICATION_KEY_ID=
# B2_APPLICATION_KEY=
# B2_ENDPOINT=
# B2_BUCKET_NAME=
# B2_REGION=

# Notification URL for backup failures (optional)
# NOTIFICATION_URL=
```

You can then manage the application using `just`:

```bash
# Start all services in the background
just prod-up

# View the logs
just prod-logs

# Stop all services
just prod-down

# Execute a command in the running container
just prod-exec ls -la

# Get an interactive shell
just prod-shell
```

The application will be available on the domain you configure in `.prod.env`.

#### With Local Binary

You can also run a production-like instance locally without Docker. This uses the compiled binary and respects the `ENV` environment variable.

```bash
just build
ENV=production ./bin/bambino serve
```
The application will be available at `http://localhost:8080`.

### Creating a User

Create an initial user and baby. The date should be in `YYYY-MM-DD` format.

```bash
./bin/bambino create-user -u <username> -b <babyname> -d <date_of_birth>
```

### Command-Line Help

You can get help for any command by passing the `--help` flag.

```bash
./bin/bambino --help
./bin/bambino create-user --help
```

## Testing

Run the test suite with:

```bash
just test
```

To see test coverage statistics:

```bash
just test-coverage
```
