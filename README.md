# Baby Tracker

A self-hosted baby activity tracking app for personal use. Tracks feeding, sleeping, diapers, and other baby activities with complete data ownership.

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
- **Database**: SQLite (default) or PostgreSQL
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
./bin/baby-tracker db migrate
```

Alternatively, you can run the migrations without building the binary:
```bash
go run cmd/baby-tracker/main.go db migrate
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

In production mode, you can run the server via:
```bash
ENV=production ./bin/baby-tracker serve
```
The application will be available at `http://localhost:8080`.

### Creating a User

Create an initial user and baby. The date should be in `YYYY-MM-DD` format.

```bash
./bin/baby-tracker create-user -u <username> -b <babyname> -d <date_of_birth>
```

### Command-Line Help

You can get help for any command by passing the `--help` flag.

```bash
./bin/baby-tracker --help
./bin/baby-tracker create-user --help
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
