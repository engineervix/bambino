# Baby Tracker

A self-hosted baby activity tracking application with complete data ownership and privacy.

## Features

- Track feeding, pumping, diapers, sleep, growth, health records, and milestones
- Dark mode by default for nighttime use
- Mobile-first responsive design
- Timer functionality for activities
- Single binary deployment with embedded frontend

## Quick Start

### Development

1. Backend:
```bash
go run cmd/server/main.go
```

2. Frontend:
```bash
cd web
npm install
npm run dev
```

### Production Build

```bash
./scripts/build.sh
```

### Create User

```bash
go run scripts/create-user.go --username=parent --password=yourpassword
```

## Technology Stack

- Backend: Go with Echo framework
- Frontend: Vue 3 with Vuetify
- Database: SQLite (default) or PostgreSQL
- Authentication: Session-based

## License

MIT
