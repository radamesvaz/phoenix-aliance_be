#!/bin/bash

# Script to run all database migrations
# This script will try to find PostgreSQL installation and run all migration files

DB_NAME="${DB_NAME:-phoenix_alliance}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

echo "Running database migrations for: $DB_NAME"

# Try to find psql in PATH first
if command -v psql &> /dev/null; then
    PSQL_CMD="psql"
    echo "Found psql in PATH"
else
    # Try common installation paths (Windows)
    PSQL_PATHS=(
        "/c/Program Files/PostgreSQL/16/bin/psql.exe"
        "/c/Program Files/PostgreSQL/15/bin/psql.exe"
        "/c/Program Files/PostgreSQL/14/bin/psql.exe"
        "/c/Program Files/PostgreSQL/13/bin/psql.exe"
        "/c/Program Files/PostgreSQL/12/bin/psql.exe"
    )
    
    PSQL_CMD=""
    for path in "${PSQL_PATHS[@]}"; do
        if [ -f "$path" ]; then
            PSQL_CMD="$path"
            echo "Found psql at: $path"
            break
        fi
    done
    
    # Also try Linux/Mac paths
    if [ -z "$PSQL_CMD" ]; then
        LINUX_PATHS=(
            "/usr/bin/psql"
            "/usr/local/bin/psql"
            "/opt/homebrew/bin/psql"
        )
        for path in "${LINUX_PATHS[@]}"; do
            if [ -f "$path" ]; then
                PSQL_CMD="$path"
                echo "Found psql at: $path"
                break
            fi
        done
    fi
fi

if [ -z "$PSQL_CMD" ]; then
    echo ""
    echo "ERROR: PostgreSQL psql command not found!"
    echo "Please add PostgreSQL bin directory to your PATH"
    exit 1
fi

# Set password if provided
if [ -n "$DB_PASSWORD" ]; then
    export PGPASSWORD="$DB_PASSWORD"
fi

# Build connection string
CONN_STRING="-h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME"

# Get migration files in order
MIGRATION_DIR="migrations"
if [ ! -d "$MIGRATION_DIR" ]; then
    echo "ERROR: migrations directory not found!"
    exit 1
fi

MIGRATION_FILES=($(ls "$MIGRATION_DIR"/*_up.sql 2>/dev/null | sort))

if [ ${#MIGRATION_FILES[@]} -eq 0 ]; then
    echo "No migration files found in $MIGRATION_DIR/ directory!"
    exit 1
fi

echo ""
echo "Found ${#MIGRATION_FILES[@]} migration file(s) to run"

# Run each migration
SUCCESS_COUNT=0
FAIL_COUNT=0

for file in "${MIGRATION_FILES[@]}"; do
    echo ""
    echo "Running: $(basename $file)..."
    if $PSQL_CMD $CONN_STRING -f "$file" 2>/dev/null; then
        echo "✓ Successfully applied: $(basename $file)"
        ((SUCCESS_COUNT++))
    else
        echo "✗ Failed to apply: $(basename $file)"
        ((FAIL_COUNT++))
    fi
done

echo ""
echo "========================================"
echo "Migration Summary:"
echo "  Successful: $SUCCESS_COUNT"
echo "  Failed: $FAIL_COUNT"
echo "========================================"

if [ $FAIL_COUNT -gt 0 ]; then
    exit 1
fi










