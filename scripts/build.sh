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
go build -o baby-tracker cmd/server/main.go

echo "Build complete! Binary: ./baby-tracker"
