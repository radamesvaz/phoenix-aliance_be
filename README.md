# Phoenix Alliance - Strength Training Tracker Backend

A Go-based backend API for tracking strength training workouts, exercises, and progression metrics.

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
cmd/
  └── server/          # Application entry point
internal/
  ├── models/          # Domain models and DTOs
  ├── repository/      # Data access layer (interfaces + implementations)
  ├── service/         # Business logic layer
  ├── handler/         # HTTP handlers (presentation layer)
  ├── router/          # Route configuration
  ├── middleware/      # HTTP middleware (auth, CORS)
  ├── auth/            # Authentication utilities (JWT, password hashing)
  ├── config/          # Configuration management
  └── database/        # Database connection
migrations/             # SQL migration files
scripts/                # Utility scripts (seed data, etc.)
```

## 📋 Features

- **User Management**: Registration and JWT-based authentication
- **Exercise Tracking**: Create and manage custom exercises
- **Workout Management**: Track workout sessions with dates
- **Set Tracking**: Record sets with weight, reps, rest time, RPE, and notes
- **Progress Analytics**: View exercise history and progress metrics over time
- **Metrics Calculation**: Automatic calculation of volume, averages, and trends

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- Make (optional, for convenience commands)

### Quick Start

The easiest way to get started:

```bash
# 1. Clone the repository
git clone <repository-url>
cd phoenix-aliance_be

# 2. Install Go dependencies
go mod download

# 3. Run everything (database + server)
./run.sh start
```

   The `run.sh` script will:
   - Check if Docker is running
   - Verify `.env` file exists
   - Start PostgreSQL container
   - Wait for database to be ready
   - Start the Go server

### Installation (Detailed)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd phoenix-aliance_be
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the root directory:
   ```env
   # Server Configuration
   SERVER_HOST=localhost
   SERVER_PORT=8080

   # Database Configuration (Docker)
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=phoenix_alliance
   DB_SSLMODE=disable

   # JWT Configuration
   # IMPORTANT: Change this to a secure random string in production!
   # You can generate one with: openssl rand -base64 32
   JWT_SECRET=your-super-secret-jwt-key-change-in-production
   JWT_EXPIRY_HOURS=24
   ```

4. **Start PostgreSQL with Docker**
   ```bash
   # Start PostgreSQL container (runs in background)
   docker-compose up -d
   
   # Check if container is running
   docker-compose ps
   
   # View logs
   docker-compose logs -f postgres
   ```

   This will:
   - Create and start a PostgreSQL 16 container
   - Create the database `phoenix_alliance` automatically
   - Run migrations automatically on first startup (only `*_up.sql` files)
   - Expose PostgreSQL on port `5432`

   **Note:** Migrations are automatically executed when the container is created for the first time. If you need to re-run them, reset the volume (`docker-compose down -v`) and start again.

5. **Run the application**

   **Option A: Using the run script (Recommended)**
   ```bash
   ./run.sh start
   ```
   
   This script will:
   - Check if Docker is running
   - Verify `.env` file exists
   - Start PostgreSQL container
   - Wait for database to be ready
   - Start the Go server
   
   **Available commands:**
   ```bash
   ./run.sh start    # Start database and server
   ./run.sh stop     # Stop database container
   ./run.sh restart  # Restart database container
   ./run.sh logs     # Show database logs
   ./run.sh status   # Show container status
   ./run.sh psql     # Connect to database
   ./run.sh seed     # Seed database with sample data
   ```
   
   **Option B: Manual start**
   ```bash
   # Start database
   docker-compose up -d
   
   # Run the server
   go run cmd/server/main.go
   ```

   The server will start on `http://localhost:8080`

### Docker Commands

**Start database:**
```bash
docker-compose up -d
```

**Stop database:**
```bash
docker-compose down
```

**Stop and remove all data (fresh start):**
```bash
docker-compose down -v
```

**View database logs:**
```bash
docker-compose logs -f postgres
```

**Connect to database:**
```bash
# Using docker exec
docker-compose exec postgres psql -U postgres -d phoenix_alliance

# Or using psql from your machine (if installed)
psql -h localhost -p 5432 -U postgres -d phoenix_alliance
```

**Reset database (removes all data and reruns migrations):**
```bash
docker-compose down -v
docker-compose up -d
```

### Manual Migration

If you need to run migrations manually (e.g., after adding new migration files):

```bash
# Using docker exec
docker-compose exec postgres psql -U postgres -d phoenix_alliance -f /migrations/001_create_users_table.up.sql

# Or copy migrations and run from your machine
for file in migrations/*_up.sql; do
  docker-compose exec -T postgres psql -U postgres -d phoenix_alliance < "$file"
done
```

### Troubleshooting

#### Docker not running

If you get errors like `docker: command not found` or `docker-compose: command not found`:

1. **Install Docker Desktop:**
   - Download from: https://www.docker.com/products/docker-desktop
   - Make sure Docker Desktop is running (you should see the Docker icon in your system tray)

2. **Verify Docker is running:**
   ```bash
   docker --version
   docker-compose --version
   ```

#### Port 5432 already in use

If you get an error that port 5432 is already in use, you likely have PostgreSQL running locally.

**To stop local PostgreSQL service:**

**Option 1: Using PowerShell (as Administrator)**
```powershell
# Find PostgreSQL services
Get-Service | Where-Object {$_.Name -like '*postgresql*'}

# Stop all PostgreSQL services
Get-Service postgresql* | Stop-Service
```

**Option 2: Using the provided script**
```powershell
# PowerShell (as Administrator)
.\scripts\stop-local-postgres.ps1

# Or Bash/Git Bash
./scripts/stop-local-postgres.sh
```

**Option 3: Using Command Prompt (as Administrator)**
```cmd
net stop postgresql-x64-16
```
(Replace `16` with your PostgreSQL version number)

**Option 4: Using Services GUI**
1. Press `Win+R`, type `services.msc`
2. Find service named `postgresql-x64-*`
3. Right-click → Stop

**Option 5: Check what's using the port**
```bash
# Windows
netstat -ano | findstr :5432

# Then stop the process using the PID shown
taskkill /PID <pid> /F
```

#### Database connection errors

If your Go application can't connect to the database:

1. **Check if container is running:**
   ```bash
   docker-compose ps
   ```

2. **Check container logs:**
   ```bash
   docker-compose logs postgres
   ```

3. **Verify connection settings in `.env`:**
   - Make sure `DB_HOST=localhost`
   - Make sure `DB_PORT=5432` (or the port you configured)
   - Make sure `DB_USER=postgres` and `DB_PASSWORD=postgres`

4. **Test connection:**
   ```bash
   docker-compose exec postgres psql -U postgres -d phoenix_alliance -c "SELECT 1;"
   ```

#### Migrations not running automatically

If migrations didn't run on first startup:

1. **Check if migrations directory is mounted correctly:**
   ```bash
   docker-compose exec postgres ls -la /migrations
   ```

2. **Run migrations manually:**
   ```bash
   for file in migrations/*_up.sql; do
     docker-compose exec -T postgres psql -U postgres -d phoenix_alliance < "$file"
   done
   ```

3. **Or reset the database:**
   ```bash
   docker-compose down -v
   docker-compose up -d
   ```

### Seed Data (Optional)

To populate the database with sample data:

```bash
# Make sure Docker container is running
docker-compose up -d

# Run the seed script
go run scripts/seed.go
```

This creates:
- A test user: `test@example.com` / `password123`
- Sample exercises (Bench Press, Squat, Deadlift, Overhead Press)
- A sample workout with sets

**Note:** The seed script connects to the database using the settings from your `.env` file, so make sure Docker is running and the database is accessible.

### Using Makefile (Optional)

For convenience, you can use the provided Makefile commands:

```bash
# Start database
make docker-up

# Stop database
make docker-down

# Reset database (removes all data)
make docker-reset

# View database logs
make docker-logs

# Connect to database
make docker-psql

# Run the server
make run

# Run tests
make test

# Seed database
make seed

# Full setup (start DB and wait for it)
make setup
```

### Using run.sh Script

The `./run.sh` script is the simplest way to manage the application:

```bash
# Start everything (database + server)
./run.sh start
```

**Available commands:**
- `./run.sh start` - Start database and server
- `./run.sh stop` - Stop database container
- `./run.sh restart` - Restart database container
- `./run.sh logs` - Show database logs (follow mode)
- `./run.sh status` - Show container status and health
- `./run.sh psql` - Connect to database with psql
- `./run.sh seed` - Seed database with sample data

The `start` command will:
- ✅ Check Docker is running
- ✅ Verify `.env` file exists
- ✅ Start PostgreSQL container
- ✅ Wait for database to be ready
- ✅ Start the Go server

Perfect for development! 🚀

## 📡 API Endpoints

### Authentication

#### POST `/signup`
Register a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### POST `/login`
Authenticate a user and receive a JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "jwt-token-here",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Exercises

#### POST `/exercises` (Protected)
Create a new exercise.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "name": "Bench Press"
}
```

#### GET `/exercises` (Protected)
Get all exercises for the authenticated user.

**Headers:**
```
Authorization: Bearer <jwt-token>
```

#### GET `/exercises/{id}/history` (Protected)
Get complete history and metrics for an exercise.

**Query Parameters:**
- None

**Response:**
```json
{
  "exercise_id": "uuid",
  "exercise_name": "Bench Press",
  "sets": [...],
  "metrics": {
    "total_sets": 10,
    "total_volume": 5000.0,
    "max_weight": 100.0,
    "max_reps": 12,
    "average_weight": 75.5,
    "average_reps": 8.2,
    "average_rest": 120.0,
    "average_rpe": 7.5,
    "first_recorded_at": "2024-01-01T00:00:00Z",
    "last_recorded_at": "2024-01-15T00:00:00Z"
  }
}
```

#### GET `/exercises/{id}/progress` (Protected)
Get progress data for a specific time range.

**Query Parameters:**
- `range`: `week`, `month`, or `year` (default: `month`)

**Response:**
```json
{
  "exercise_id": "uuid",
  "exercise_name": "Bench Press",
  "range": "month",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-02-01T00:00:00Z",
  "data_points": [
    {
      "date": "2024-01-01T00:00:00Z",
      "total_volume": 500.0,
      "max_weight": 80.0,
      "total_sets": 5,
      "average_rpe": 7.5
    }
  ],
  "summary": {...}
}
```

### Workouts

#### POST `/workouts` (Protected)
Create a new workout.

**Request Body:**
```json
{
  "date": "2024-01-15T10:00:00Z"
}
```

#### POST `/workouts/{id}/sets` (Protected)
Add a set to a workout.

**Request Body:**
```json
{
  "exercise_id": "uuid",
  "weight": 80.0,
  "reps": 8,
  "rest_seconds": 120,
  "notes": "Felt strong today",
  "rpe": 7
}
```

### Health Check

#### GET `/health`
Check if the server is running.

**Response:**
```
OK
```

## 🧪 Testing

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## 📦 Project Structure

```
phoenix-aliance_be/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── models/                  # Domain models
│   │   ├── user.go
│   │   ├── exercise.go
│   │   ├── workout.go
│   │   ├── set.go
│   │   └── progress.go
│   ├── repository/              # Data access layer
│   │   ├── user_repository.go
│   │   ├── exercise_repository.go
│   │   ├── workout_repository.go
│   │   └── set_repository.go
│   ├── service/                 # Business logic
│   │   ├── user_service.go
│   │   ├── exercise_service.go
│   │   ├── workout_service.go
│   │   └── set_service.go
│   ├── handler/                 # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── exercise_handler.go
│   │   ├── workout_handler.go
│   │   └── response.go
│   ├── router/
│   │   └── router.go            # Route setup
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication
│   │   └── cors.go              # CORS handling
│   ├── auth/
│   │   ├── jwt.go               # JWT utilities
│   │   └── password.go          # Password hashing
│   ├── config/
│   │   └── config.go            # Configuration management
│   └── database/
│       └── database.go          # Database connection
├── migrations/                  # SQL migrations
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_exercises_table.up.sql
│   ├── 002_create_exercises_table.down.sql
│   ├── 003_create_workouts_table.up.sql
│   ├── 003_create_workouts_table.down.sql
│   ├── 004_create_sets_table.up.sql
│   └── 004_create_sets_table.down.sql
├── scripts/
│   └── seed.go                  # Seed data script
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## 🔒 Security Notes

- Passwords are hashed using bcrypt
- JWT tokens are used for authentication
- All protected routes require a valid JWT token in the Authorization header
- User data is isolated (users can only access their own data)
- SQL injection protection via parameterized queries

## 🚧 Future Enhancements

- [ ] Add pagination for list endpoints
- [ ] Add filtering and sorting options
- [ ] Implement workout templates
- [ ] Add exercise categories and tags
- [ ] Add 1RM (one-rep max) calculations
- [ ] Add volume progression charts
- [ ] Add workout notes and comments
- [ ] Add social features (sharing workouts)
- [ ] Add export functionality (CSV, PDF)

## 📝 License

[Your License Here]

## 👥 Contributing

[Your Contributing Guidelines Here]

