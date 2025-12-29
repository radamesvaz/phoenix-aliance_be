package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"phoenix-alliance-be/internal/middleware"
	"phoenix-alliance-be/internal/models"
	"phoenix-alliance-be/internal/service"
)

// WorkoutHandler handles workout-related requests
type WorkoutHandler struct {
	workoutService service.WorkoutService
	setService     service.SetService
}

// NewWorkoutHandler creates a new workout handler
func NewWorkoutHandler(workoutService service.WorkoutService, setService service.SetService) *WorkoutHandler {
	return &WorkoutHandler{
		workoutService: workoutService,
		setService:     setService,
	}
}

// CreateWorkout handles POST /workouts
func (h *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.WorkoutCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	workout, err := h.workoutService.CreateWorkout(userID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, workout)
}

// CreateSet handles POST /workouts/{id}/sets
func (h *WorkoutHandler) CreateSet(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	workoutID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid workout ID")
		return
	}

	var req models.SetCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.Weight < 0 {
		respondWithError(w, http.StatusBadRequest, "Weight must be non-negative")
		return
	}
	if req.Reps < 1 {
		respondWithError(w, http.StatusBadRequest, "Reps must be at least 1")
		return
	}
	if req.RPE != nil && (*req.RPE < 1 || *req.RPE > 10) {
		respondWithError(w, http.StatusBadRequest, "RPE must be between 1 and 10")
		return
	}

	set, err := h.setService.CreateSet(userID, workoutID, &req)
	if err != nil {
		if err.Error() == "workout not found" || err.Error() == "exercise not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, set)
}

// GetWorkouts handles GET /workouts
func (h *WorkoutHandler) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	workouts, err := h.workoutService.GetWorkouts(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, workouts)
}

// GetWorkout handles GET /workouts/{id}
func (h *WorkoutHandler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	workoutID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid workout ID")
		return
	}

	workout, err := h.workoutService.GetWorkoutByID(userID, workoutID)
	if err != nil {
		if err.Error() == "workout not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, workout)
}

// GetWorkoutSets handles GET /workouts/{id}/sets
func (h *WorkoutHandler) GetWorkoutSets(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	workoutID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid workout ID")
		return
	}

	// Verify workout belongs to user
	_, err = h.workoutService.GetWorkoutByID(userID, workoutID)
	if err != nil {
		if err.Error() == "workout not found" {
			respondWithError(w, http.StatusNotFound, "workout not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get sets for workout
	sets, err := h.setService.GetWorkoutSets(workoutID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, sets)
}

