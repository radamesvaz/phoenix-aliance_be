#!/bin/bash

# Script to initialize database with migrations
# This script is used by Docker to run migrations on first startup
# It can also be run manually to apply migrations

set -e

echo "Waiting for PostgreSQL to be ready..."

# Wait for PostgreSQL to be ready
until pg_isready -h localhost -p 5432 -U postgres; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "PostgreSQL is ready!"

# Get migration files in order
MIGRATION_DIR="/docker-entrypoint-initdb.d"
if [ ! -d "$MIGRATION_DIR" ]; then
    echo "Using migrations from current directory..."
    MIGRATION_DIR="./migrations"
fi

MIGRATION_FILES=($(ls "$MIGRATION_DIR"/*_up.sql 2>/dev/null | sort))

if [ ${#MIGRATION_FILES[@]} -eq 0 ]; then
    echo "No migration files found!"
    exit 1
fi

echo "Found ${#MIGRATION_FILES[@]} migration file(s) to run"

# Run each migration
for file in "${MIGRATION_FILES[@]}"; do
    echo "Running migration: $(basename $file)..."
    psql -v ON_ERROR_STOP=1 -U postgres -d phoenix_alliance -f "$file"
    echo "✓ Applied: $(basename $file)"
done

echo ""
echo "Database initialization completed successfully!"












