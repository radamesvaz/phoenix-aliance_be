package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string
	Expiry    int // in hours
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowAllOrigins   bool
	AllowedOrigins    []string
	AllowedMethods    string
	AllowedHeaders    string
	AllowCredentials  bool
	MaxAgeSeconds     int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "phoenix_alliance"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiry:    getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
		CORS: loadCORSConfig(),
	}

	// Validate required fields
	if config.JWT.SecretKey == "your-secret-key-change-in-production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in environment variables")
	}

	return config, nil
}

func loadCORSConfig() CORSConfig {
	originsRaw := getEnv("CORS_ALLOWED_ORIGINS", "*")
	origins := splitCSV(originsRaw)
	allowAll := false
	if len(origins) == 1 && origins[0] == "*" {
		allowAll = true
		origins = nil
	}

	return CORSConfig{
		AllowAllOrigins:  allowAll,
		AllowedOrigins:   origins,
		AllowedMethods:   getEnv("CORS_ALLOWED_METHODS", "GET, POST, PUT, DELETE, OPTIONS"),
		AllowedHeaders:   getEnv("CORS_ALLOWED_HEADERS", "Content-Type, Authorization"),
		AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", false),
		MaxAgeSeconds:    getEnvAsInt("CORS_MAX_AGE_SECONDS", 300),
	}
}

// DSN returns the database connection string
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}

	parts := make([]string, 0, 4)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := strings.TrimSpace(s[start:i])
			if part != "" {
				parts = append(parts, part)
			}
			start = i + 1
		}
	}
	return parts
}

