package main

import (
	"log"
	"os"
	"time"

	"phoenix-alliance-be/internal/auth"
	"phoenix-alliance-be/internal/config"
	"phoenix-alliance-be/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Get seed user credentials from environment variables
	seedEmail := getEnv("SEED_USER_EMAIL", "test@example.com")
	seedPassword := getEnv("SEED_USER_PASSWORD", "password123")

	// Create a test user
	hashedPassword, err := auth.HashPassword(seedPassword)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	var userID int64
	err = database.DB.QueryRow(`
		INSERT INTO users (email, password, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id_user
	`, seedEmail, hashedPassword, time.Now()).Scan(&userID)

	if err != nil {
		log.Printf("Error creating user (may already exist): %v", err)
		// Try to get existing user ID
		err = database.DB.QueryRow("SELECT id_user FROM users WHERE email = $1", seedEmail).Scan(&userID)
		if err != nil {
			log.Fatalf("Failed to get user ID: %v", err)
		}
		log.Printf("Using existing user: %s", seedEmail)
	} else {
		log.Printf("Created test user: %s", seedEmail)
	}

	// Create some exercises
	exercises := []struct {
		name string
	}{
		{"Bench Press"},
		{"Squat"},
		{"Deadlift"},
		{"Overhead Press"},
	}

	var exerciseIDs []int64
	for _, ex := range exercises {
		var exerciseID int64
		err = database.DB.QueryRow(`
			INSERT INTO exercises (user_id, name, created_at)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
			RETURNING id_exercise
		`, userID, ex.name, time.Now()).Scan(&exerciseID)
		if err != nil {
			// Exercise might already exist, try to get its ID
			err = database.DB.QueryRow(`
				SELECT id_exercise FROM exercises WHERE user_id = $1 AND name = $2
			`, userID, ex.name).Scan(&exerciseID)
			if err != nil {
				log.Printf("Error creating/getting exercise %s: %v", ex.name, err)
				continue
			}
		}
		exerciseIDs = append(exerciseIDs, exerciseID)
		log.Printf("Created/found exercise: %s (ID: %d)", ex.name, exerciseID)
	}

	if len(exerciseIDs) == 0 {
		log.Fatalf("No exercises created. Cannot create workout.")
	}

	// Create a workout
	workoutDate := time.Now().AddDate(0, 0, -1) // Yesterday
	var workoutID int64
	err = database.DB.QueryRow(`
		INSERT INTO workouts (user_id, date, created_at)
		VALUES ($1, $2, $3)
		RETURNING id_workout
	`, userID, workoutDate, time.Now()).Scan(&workoutID)
	if err != nil {
		log.Printf("Error creating workout: %v", err)
	} else {
		log.Printf("Created workout (ID: %d)", workoutID)
	}

	// Create some sets for the first exercise (Bench Press)
	benchPressID := exerciseIDs[0]
	sets := []struct {
		weight float64
		reps   int
		rpe    *int
	}{
		{60.0, 10, intPtr(7)},
		{65.0, 8, intPtr(8)},
		{70.0, 6, intPtr(9)},
	}

	for i, s := range sets {
		restSeconds := 120
		_, err = database.DB.Exec(`
			INSERT INTO sets (workout_id, exercise_id, weight, reps, rest_seconds, rpe, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, workoutID, benchPressID, s.weight, s.reps, restSeconds, s.rpe, time.Now().Add(time.Duration(i)*time.Minute))
		if err != nil {
			log.Printf("Error creating set: %v", err)
		} else {
			log.Printf("Created set: %.1fkg x %d reps", s.weight, s.reps)
		}
	}

	log.Println("Seed data created successfully!")
}

func intPtr(i int) *int {
	return &i
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
