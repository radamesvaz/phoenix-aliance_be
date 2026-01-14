.PHONY: build run test clean migrate-up migrate-down seed docker-up docker-down docker-logs docker-reset deps

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-down-volumes:
	docker-compose down -v

docker-logs:
	docker-compose logs -f postgres

docker-reset: docker-down-volumes docker-up
	@echo "Database reset complete. Migrations will run automatically."

docker-psql:
	docker-compose exec postgres psql -U postgres -d phoenix_alliance

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run migrations up (uses cmd/migrate/main.go)
migrate-up:
	go run cmd/migrate/main.go up

# Run migrations down (uses cmd/migrate/main.go)
migrate-down:
	go run cmd/migrate/main.go down

# Seed the database
seed:
	go run scripts/seed.go

# Install dependencies
deps:
	go mod download
	go mod tidy

# Setup: Start database and seed (if needed)
setup: docker-up
	@echo "Waiting for database to be ready..."
	@sleep 3
	@echo "Setup complete! Run 'make run' to start the server."