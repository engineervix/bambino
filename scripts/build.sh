#!/bin/bash

set -e

echo "Building Baby Tracker..."

# Build frontend
echo "Building frontend..."
cd web
npm install
npm run build
cd ..

# Embed frontend into Go binary
echo "Building backend with embedded frontend..."
go build -o baby-tracker cmd/baby-tracker/main.go

echo "Build complete! Binary: ./baby-tracker"
echo ""
echo "Usage:"
echo "  ./baby-tracker serve            - Start the server"
echo "  ./baby-tracker create-user      - Create a new user"
echo "  ./baby-tracker db test          - Test database connection"
echo "  ./baby-tracker db migrate       - Run database migrations"
echo "  ./baby-tracker --help           - Show all commands"
echo ""
