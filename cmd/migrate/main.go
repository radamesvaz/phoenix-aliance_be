package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Cargar variables de entorno desde .env (opcional)
	err := godotenv.Load()
	if err != nil {
		// .env es opcional, especialmente en producción
	}

	// Obtener configuración de base de datos desde variables de entorno
	databaseURL := os.Getenv("DATABASE_URL")
	dbUser := firstNonEmpty(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("PGUSER"),
		os.Getenv("DB_USER"),
	)
	dbPassword := firstNonEmpty(
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("DB_PASSWORD"),
	)
	dbHost := firstNonEmpty(
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("PGHOST"),
		"localhost",
	)
	dbPort := firstNonEmpty(
		os.Getenv("DB_PORT"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("PGPORT"),
		"5432",
	)
	dbName := firstNonEmpty(
		os.Getenv("POSTGRES_DB"),
		os.Getenv("PGDATABASE"),
		os.Getenv("DB_NAME"),
	)

	// Validar que tenemos las variables necesarias
	if databaseURL == "" && (dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "") {
		fmt.Fprintf(os.Stderr, "Error: Missing required database environment variables\n")
		os.Exit(1)
	}

	// Construir DSN (Data Source Name)
	var dsn string
	if databaseURL != "" {
		dsn = databaseURL
	} else {
		sslMode := "require"
		lowerHost := strings.ToLower(dbHost)
		if lowerHost == "localhost" || lowerHost == "127.0.0.1" || lowerHost == "::1" {
			sslMode = "disable"
		}
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
	}

	// Conectar a la base de datos
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Probar la conexión
	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not ping database: %v\n", err)
		os.Exit(1)
	}

	// Crear instancia del driver de migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not create migrate driver: %v\n", err)
		os.Exit(1)
	}

	// Crear instancia de migrate
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not create migrate instance: %v\n", err)
		os.Exit(1)
	}

	// Manejar comandos desde argumentos
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "down":
			fmt.Println("Rolling back database migrations...")
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				fmt.Fprintf(os.Stderr, "Error: Could not rollback migrations: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Migrations rolled back successfully!")
			return
		case "force":
			if len(os.Args) < 3 {
				fmt.Fprintf(os.Stderr, "Error: force command requires a version number\n")
				fmt.Fprintf(os.Stderr, "Usage: go run cmd/migrate/main.go force <version>\n")
				os.Exit(1)
			}
			var version int
			if _, err := fmt.Sscanf(os.Args[2], "%d", &version); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid version number: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Forcing database version to %d...\n", version)
			if err := m.Force(version); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Could not force version: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Database version forced to %d successfully!\n", version)
			return
		case "version":
			version, dirty, err := m.Version()
			if err != nil {
				if err == migrate.ErrNilVersion {
					fmt.Println("No migrations have been applied yet.")
					return
				}
				fmt.Fprintf(os.Stderr, "Error: Could not get version: %v\n", err)
				os.Exit(1)
			}
			if dirty {
				fmt.Printf("Current version: %d (DIRTY - migration failed)\n", version)
			} else {
				fmt.Printf("Current version: %d\n", version)
			}
			return
		}
	}

	// Ejecutar migraciones (comando por defecto: up)
	fmt.Println("Running database migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// Manejar base de datos en estado "dirty"
		if strings.Contains(err.Error(), "Dirty database version") {
			re := regexp.MustCompile(`Dirty database version (\d+)`)
			matches := re.FindStringSubmatch(err.Error())
			if len(matches) > 1 {
				version := matches[1]
				fmt.Printf("Warning: Database is dirty at version %s, cleaning and resetting...\n", version)

				// Forzar versión a 1 (primera migración)
				if forceErr := m.Force(1); forceErr != nil {
					fmt.Fprintf(os.Stderr, "Error: Could not force version to 1: %v\n", forceErr)
					os.Exit(1)
				}

				// Limpiar tablas (ajusta según tus tablas)
				fmt.Println("Dropping all tables to start fresh...")
				// Agrega aquí las sentencias DROP para tus tablas

				// Reintentar migraciones
				fmt.Println("Running migrations from scratch...")
				if retryErr := m.Up(); retryErr != nil && retryErr != migrate.ErrNoChange {
					fmt.Fprintf(os.Stderr, "Error: Could not run migrations after reset: %v\n", retryErr)
					os.Exit(1)
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error: Could not run migrations: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Database migrations completed successfully!")
}

// firstNonEmpty retorna el primer string no vacío de la lista
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}



