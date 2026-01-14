package router

import (
	"net/http"

	"phoenix-alliance-be/internal/config"
	"phoenix-alliance-be/internal/handler"
	"phoenix-alliance-be/internal/middleware"
	"phoenix-alliance-be/internal/service"

	"github.com/gorilla/mux"
)

// SetupRouter configures and returns the application router
func SetupRouter(
	cfg *config.Config,
	userService service.UserService,
	exerciseService service.ExerciseService,
	workoutService service.WorkoutService,
	setService service.SetService,
) *mux.Router {
	router := mux.NewRouter()

	// Apply CORS middleware to all routes
	router.Use(middleware.CORSMiddleware(cfg))

	// Create handlers
	authHandler := handler.NewAuthHandler(userService, &jwtConfigAdapter{cfg: cfg})
	exerciseHandler := handler.NewExerciseHandler(exerciseService, setService)
	workoutHandler := handler.NewWorkoutHandler(workoutService, setService)

	// Public routes (no authentication required)
	router.HandleFunc("/signup", authHandler.Signup).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")

	// Protected routes (authentication required)
	api := router.PathPrefix("/").Subrouter()
	api.Use(middleware.AuthMiddleware(cfg))

	// Exercise routes
	api.HandleFunc("/exercises", exerciseHandler.CreateExercise).Methods("POST", "OPTIONS")
	api.HandleFunc("/exercises", exerciseHandler.GetExercises).Methods("GET", "OPTIONS")
	api.HandleFunc("/exercises/{id}", exerciseHandler.UpdateExercise).Methods("PUT", "OPTIONS")
	api.HandleFunc("/exercises/{id}", exerciseHandler.DeleteExercise).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/exercises/{id}/history", exerciseHandler.GetExerciseHistory).Methods("GET", "OPTIONS")
	api.HandleFunc("/exercises/{id}/progress", exerciseHandler.GetExerciseProgress).Methods("GET", "OPTIONS")

	// Workout routes
	api.HandleFunc("/workouts", workoutHandler.CreateWorkout).Methods("POST", "OPTIONS")
	api.HandleFunc("/workouts", workoutHandler.GetWorkouts).Methods("GET", "OPTIONS")
	api.HandleFunc("/workouts/{id}", workoutHandler.GetWorkout).Methods("GET", "OPTIONS")
	api.HandleFunc("/workouts/{id}", workoutHandler.DeleteWorkout).Methods("DELETE", "OPTIONS")
	api.HandleFunc("/workouts/{id}/sets", workoutHandler.CreateSet).Methods("POST", "OPTIONS")
	api.HandleFunc("/workouts/{id}/sets", workoutHandler.GetWorkoutSets).Methods("GET", "OPTIONS")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
}

// jwtConfigAdapter adapts config.Config to handler.AuthConfig interface
type jwtConfigAdapter struct {
	cfg *config.Config
}

func (a *jwtConfigAdapter) GetJWTSecret() string {
	return a.cfg.JWT.SecretKey
}

func (a *jwtConfigAdapter) GetJWTExpiry() int {
	return a.cfg.JWT.Expiry
}
