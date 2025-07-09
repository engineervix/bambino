#!/bin/sh
set -e

# Run database migrations
echo "Running migrations..."
./bambino db migrate

# Run whatever command was passed
exec "$@"
