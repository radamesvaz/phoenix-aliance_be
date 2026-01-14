#!/bin/bash

# Phoenix Alliance - Strength Training Tracker Backend
# Run script to manage database and server

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to show usage
show_usage() {
    echo "Usage: $0 {start|stop|restart|logs|status|psql|migrate|seed|test}"
    echo ""
    echo "Commands:"
    echo "  start   - Start database and server"
    echo "  stop    - Stop database container"
    echo "  restart - Restart database container"
    echo "  logs    - Show database logs"
    echo "  status  - Show container status"
    echo "  psql    - Connect to database"
    echo "  migrate - Run database migrations (up by default)"
    echo "    migrate up     - Apply all pending migrations"
    echo "    migrate down   - Rollback last migration"
    echo "    migrate version - Show current migration version"
    echo "    migrate force <version> - Force migration version"
    echo "  seed    - Seed database with sample data"
    echo "  test    - Run all tests in the project"
    echo ""
    exit 1
}

# Function to check Docker
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo -e "${RED}❌ Docker is not running!${NC}"
        echo "Please start Docker Desktop and try again."
        exit 1
    fi
    echo -e "${GREEN}✓ Docker is running${NC}"
}

# Function to check .env file
check_env() {
    if [ ! -f .env ]; then
        echo -e "${RED}❌ .env file not found${NC}"
        echo "Please create a .env file before running the application."
        echo "See README.md for the required configuration."
        exit 1
    fi
}

# Function to check and free port 5432
check_and_free_port() {
    # Check if port 5432 is in use
    if command -v netstat >/dev/null 2>&1; then
        # Windows/Git Bash
        port_check=$(netstat -ano 2>/dev/null | grep -E ":[[:space:]]*5432[[:space:]]" | grep LISTENING || true)
    elif command -v ss >/dev/null 2>&1; then
        # Linux
        port_check=$(ss -tlnp 2>/dev/null | grep ":5432 " || true)
    else
        return 0  # Can't check, continue anyway
    fi
    
    if [ -z "$port_check" ]; then
        return 0  # Port is free
    fi
    
    echo -e "${YELLOW}⚠ Port 5432 is already in use${NC}"
    echo "Checking what's using it..."
    
    # Try to find Docker containers using port 5432
    docker_containers=$(docker ps --format "{{.Names}}\t{{.Ports}}" 2>/dev/null | grep ":5432" || true)
    
    if [ -n "$docker_containers" ]; then
        echo -e "${BLUE}Found Docker containers using port 5432:${NC}"
        echo "$docker_containers" | while IFS=$'\t' read -r name ports; do
            echo "  - $name ($ports)"
        done
        
        echo ""
        echo -e "${YELLOW}Attempting to stop conflicting containers...${NC}"
        
        # Stop containers that are using port 5432
        echo "$docker_containers" | while IFS=$'\t' read -r name ports; do
            if [ "$name" != "phoenix_alliance_db" ]; then
                echo "  Stopping: $name"
                docker stop "$name" 2>/dev/null || true
            fi
        done
        
        # Wait a moment for ports to be released
        sleep 2
        
        # Verify port is now free
        if command -v netstat >/dev/null 2>&1; then
            port_check=$(netstat -ano 2>/dev/null | grep -E ":[[:space:]]*5432[[:space:]]" | grep LISTENING || true)
        elif command -v ss >/dev/null 2>&1; then
            port_check=$(ss -tlnp 2>/dev/null | grep ":5432 " || true)
        fi
        
        if [ -z "$port_check" ]; then
            echo -e "${GREEN}✓ Port 5432 is now free${NC}"
            return 0
        fi
    fi
    
    # If we get here, port is still in use and it's not a Docker container we can stop
    echo -e "${RED}❌ Port 5432 is still in use${NC}"
    echo ""
    echo "This could be:"
    echo "  - A local PostgreSQL service"
    echo "  - Another Docker container that couldn't be stopped automatically"
    echo "  - Another application"
    echo ""
    echo "To find what's using the port:"
    if command -v netstat >/dev/null 2>&1; then
        echo "  netstat -ano | findstr :5432"
        echo "  Then: tasklist /FI \"PID eq <PID>\""
    else
        echo "  sudo lsof -i :5432"
    fi
    echo ""
    return 1
}

# Function to start database
start_database() {
    echo "Checking database container..."
    if ! docker-compose ps | grep -q "phoenix_alliance_db.*Up"; then
        # Check and free port 5432 if needed
        if ! check_and_free_port; then
            exit 1
        fi
        
        echo "Starting PostgreSQL container..."
        
        # Try to start, capture output to variable (works better on Windows)
        # Start postgres first
        output=$(docker-compose up -d postgres 2>&1) || {
            # Check if it's a port conflict
            if echo "$output" | grep -qiE "port.*already.*allocated|address.*already.*in.*use|bind.*failed"; then
                echo ""
                echo -e "${RED}❌ Port 5432 is already in use${NC}"
                echo ""
                echo "This usually means you have PostgreSQL running locally."
                echo ""
                echo -e "${YELLOW}To stop the local PostgreSQL service:${NC}"
                echo ""
                echo -e "${BLUE}Option 1: Windows PowerShell (as Administrator)${NC}"
                echo "  Get-Service postgresql* | Stop-Service"
                echo ""
                echo -e "${BLUE}Option 2: Find and stop manually${NC}"
                echo "  Get-Service | Where-Object {\$_.Name -like '*postgresql*'}"
                echo ""
                echo -e "${BLUE}Option 3: Windows Command Prompt (as Administrator)${NC}"
                echo "  net stop postgresql-x64-16"
                echo "  (Replace '16' with your PostgreSQL version)"
                echo ""
                echo -e "${BLUE}Option 4: Use Services GUI${NC}"
                echo "  1. Press Win+R, type: services.msc"
                echo "  2. Find 'postgresql-x64-*' service"
                echo "  3. Right-click → Stop"
                echo ""
                echo -e "${BLUE}Option 5: Use the provided script${NC}"
                echo "  PowerShell: .\\scripts\\stop-local-postgres.ps1"
                echo "  Bash: ./scripts/stop-local-postgres.sh"
                echo ""
                echo -e "${YELLOW}After stopping the service, run: ./run.sh start${NC}"
                echo ""
                exit 1
            else
                echo -e "${RED}❌ Failed to start database container${NC}"
                echo ""
                echo "Error output:"
                echo "$output"
                echo ""
                echo "Check logs with: $0 logs"
                exit 1
            fi
        }
        
        # Show docker-compose output
        echo "$output"
        
        echo "Waiting for database to be ready..."
        max_attempts=30
        attempt=0
        
        while [ $attempt -lt $max_attempts ]; do
            if docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
                echo -e "${GREEN}✓ Database is ready${NC}"
                break
            fi
            attempt=$((attempt + 1))
            echo -n "."
            sleep 1
        done
        
        if [ $attempt -eq $max_attempts ]; then
            echo -e "\n${RED}❌ Database failed to start${NC}"
            echo "Check logs with: $0 logs"
            exit 1
        fi
        echo ""
        
        # Run migrations service (it will wait for postgres to be healthy)
        echo "Running database migrations..."
        migrate_output=$(docker-compose up migrate 2>&1)
        migrate_exit_code=$?
        
        if [ $migrate_exit_code -ne 0 ] || echo "$migrate_output" | grep -qiE "error|failed"; then
            echo -e "${YELLOW}⚠ Migrations service had issues, trying manual migration...${NC}"
            run_migrations
        else
            echo -e "${GREEN}✓ Migrations completed${NC}"
        fi
    else
        echo -e "${GREEN}✓ Database container is already running${NC}"
    fi
    
    # Check if database exists and is accessible
    echo "Verifying database connection..."
    if docker-compose exec -T postgres psql -U postgres -d phoenix_alliance -c "SELECT 1;" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Database connection verified${NC}"
        
        # Check if tables exist (check for users table as indicator)
        table_count=$(docker-compose exec -T postgres psql -U postgres -d phoenix_alliance -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users';" 2>/dev/null || echo "0")
        
        if [ "$table_count" = "0" ]; then
            echo -e "${YELLOW}⚠ No tables found. Running migrations manually...${NC}"
            run_migrations
        else
            echo -e "${GREEN}✓ Database schema is up to date${NC}"
            # Check migration version to ensure it's up to date
            echo "Checking migration version..."
            version_output=$(go run cmd/migrate/main.go version 2>&1 || echo "")
            if [ -n "$version_output" ]; then
                echo "  $version_output"
            fi
        fi
    else
        echo -e "${YELLOW}⚠ Database might not be fully initialized yet${NC}"
        echo "Waiting a bit more..."
        sleep 2
    fi
}

# Function to run migrations
run_migrations() {
    echo "Running database migrations using golang-migrate..."
    go run cmd/migrate/main.go
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Migrations completed${NC}"
    else
        echo -e "${RED}❌ Migrations failed${NC}"
        exit 1
    fi
}

# Function to start everything
start() {
    echo "🚀 Starting Phoenix Alliance Backend..."
    echo ""
    
    check_docker
    check_env
    start_database
    
    echo ""
    echo "=========================================="
    echo "Starting Go server..."
    echo "=========================================="
    echo ""
    
    # Run the server
    go run cmd/server/main.go
}

# Function to stop database
stop() {
    echo "Stopping database container..."
    docker-compose down
    echo -e "${GREEN}✓ Database stopped${NC}"
}

# Function to restart database
restart() {
    echo "Restarting database container..."
    docker-compose restart postgres
    echo -e "${GREEN}✓ Database restarted${NC}"
}

# Function to show logs
logs() {
    docker-compose logs -f postgres
}

# Function to show status
status() {
    echo "Container status:"
    docker-compose ps
    echo ""
    echo "Database connection test:"
    if docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Database is ready${NC}"
    else
        echo -e "${RED}❌ Database is not ready${NC}"
    fi
}

# Function to connect to database
psql() {
    docker-compose exec postgres psql -U postgres -d phoenix_alliance
}

# Function to seed database
seed() {
    check_docker
    check_env
    echo "Seeding database with sample data..."
    go run scripts/seed.go
    echo -e "${GREEN}✓ Database seeded${NC}"
}

# Function to run tests
test() {
    echo "Running all tests..."
    echo ""
    go test ./... -v
    if [ $? -eq 0 ]; then
        echo ""
        echo -e "${GREEN}✓ All tests passed${NC}"
    else
        echo ""
        echo -e "${RED}❌ Some tests failed${NC}"
        exit 1
    fi
}

# Main command handler
case "${1:-}" in
    start)
        start
        ;;
    stop)
        check_docker
        stop
        ;;
    restart)
        check_docker
        restart
        ;;
    logs)
        check_docker
        logs
        ;;
    status)
        check_docker
        status
        ;;
    psql)
        check_docker
        psql
        ;;
    migrate)
        check_docker
        check_env
        echo "Running database migrations..."
        if [ "$2" = "down" ]; then
            echo "Rolling back migrations..."
            go run cmd/migrate/main.go down
        elif [ "$2" = "version" ]; then
            go run cmd/migrate/main.go version
        elif [ "$2" = "force" ] && [ -n "$3" ]; then
            echo "Forcing migration version to $3..."
            go run cmd/migrate/main.go force "$3"
        else
            run_migrations
        fi
        ;;
    seed)
        seed
        ;;
    test)
        test
        ;;
    *)
        show_usage
        ;;
esac
