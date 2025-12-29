#!/bin/bash

# Script to create PostgreSQL database for Phoenix Alliance
# This script will try to find PostgreSQL installation and create the database

DB_NAME="${DB_NAME:-phoenix_alliance}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

echo "Creating PostgreSQL database: $DB_NAME"

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
    echo ""
    echo "Please do one of the following:"
    echo "1. Add PostgreSQL bin directory to your PATH"
    echo "   Windows: C:\\Program Files\\PostgreSQL\\<version>\\bin"
    echo "   Linux/Mac: Usually /usr/bin or /usr/local/bin"
    echo "2. Install PostgreSQL from: https://www.postgresql.org/download/"
    echo "3. Use pgAdmin to create the database manually"
    echo ""
    echo "Or set PGPASSWORD and use psql directly:"
    echo "  export PGPASSWORD=your_password"
    echo "  psql -U postgres -c \"CREATE DATABASE $DB_NAME;\""
    exit 1
fi

# Set password if provided
if [ -n "$DB_PASSWORD" ]; then
    export PGPASSWORD="$DB_PASSWORD"
fi

# Build connection string
CONN_STRING="-h $DB_HOST -p $DB_PORT -U $DB_USER"

# Check if database already exists
echo ""
echo "Checking if database exists..."
DB_EXISTS=$($PSQL_CMD $CONN_STRING -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'" 2>/dev/null)

if [ "$DB_EXISTS" = "1" ]; then
    echo "Database '$DB_NAME' already exists!"
    read -p "Do you want to drop and recreate it? (y/N): " response
    if [ "$response" = "y" ] || [ "$response" = "Y" ]; then
        echo "Dropping existing database..."
        $PSQL_CMD $CONN_STRING -c "DROP DATABASE $DB_NAME;" 2>/dev/null
    else
        echo "Keeping existing database. Exiting."
        exit 0
    fi
fi

# Create database
echo ""
echo "Creating database '$DB_NAME'..."
if $PSQL_CMD $CONN_STRING -c "CREATE DATABASE $DB_NAME;" 2>/dev/null; then
    echo "Database '$DB_NAME' created successfully!"
    echo ""
    echo "Next steps:"
    echo "1. Run migrations using: ./scripts/run_migrations.sh"
    echo "2. Or manually run each migration file"
else
    echo ""
    echo "ERROR: Failed to create database!"
    echo ""
    echo "Common issues:"
    echo "- PostgreSQL service might not be running"
    echo "- Wrong password for user '$DB_USER'"
    echo "- User '$DB_USER' doesn't have permission to create databases"
    echo ""
    echo "Try running manually:"
    echo "  $PSQL_CMD $CONN_STRING -c \"CREATE DATABASE $DB_NAME;\""
    exit 1
fi










