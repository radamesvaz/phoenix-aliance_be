package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"phoenix-alliance-be/internal/config"
	"phoenix-alliance-be/internal/database"
	"phoenix-alliance-be/internal/repository"
	"phoenix-alliance-be/internal/router"
	"phoenix-alliance-be/internal/service"
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

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)
	exerciseRepo := repository.NewExerciseRepository(database.DB)
	workoutRepo := repository.NewWorkoutRepository(database.DB)
	setRepo := repository.NewSetRepository(database.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)
	exerciseService := service.NewExerciseService(exerciseRepo)
	workoutService := service.NewWorkoutService(workoutRepo)
	setService := service.NewSetService(setRepo, exerciseRepo, workoutRepo)

	// Setup router
	r := router.SetupRouter(cfg, userService, exerciseService, workoutService, setService)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Channel to listen for errors from server
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal or server error
	select {
	case err := <-serverErrors:
		log.Fatalf("Server failed to start: %v", err)
	case sig := <-shutdown:
		log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Stop accepting new requests and wait for active requests to complete
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
			srv.Close()
		}

		log.Println("Server stopped")

		// Close database connection
		log.Println("Closing database connection...")
		database.Close()

		log.Println("Graceful shutdown completed")
	}
}
