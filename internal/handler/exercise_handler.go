package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"phoenix-alliance-be/internal/middleware"
	"phoenix-alliance-be/internal/models"
	"phoenix-alliance-be/internal/service"

	"github.com/gorilla/mux"
)

// ExerciseHandler handles exercise-related requests
type ExerciseHandler struct {
	exerciseService service.ExerciseService
	setService      service.SetService
}

// NewExerciseHandler creates a new exercise handler
func NewExerciseHandler(exerciseService service.ExerciseService, setService service.SetService) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService: exerciseService,
		setService:      setService,
	}
}

// CreateExercise handles POST /exercises
func (h *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Agregar mas campos al ejercicio:
	// Musculo principal
	// Musculo secundario
	// Tipo de ejercicio (peso, cardio, flexibilidad, etc.)
	var req models.ExerciseCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Exercise name is required")
		return
	}

	exercise, err := h.exerciseService.CreateExercise(userID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, exercise)
}

// GetExercises handles GET /exercises
func (h *ExerciseHandler) GetExercises(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	exercises, err := h.exerciseService.GetExercises(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, exercises)
}

// GetExerciseHistory handles GET /exercises/{id}/history
func (h *ExerciseHandler) GetExerciseHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	exerciseID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	history, err := h.setService.GetExerciseHistory(userID, exerciseID)
	if err != nil {
		if err.Error() == "exercise not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, history)
}

// GetExerciseProgress handles GET /exercises/{id}/progress?range=week|month|year
func (h *ExerciseHandler) GetExerciseProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	exerciseID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	// @TODO: Agregar rangos de tiempo en un enum
	rangeParam := r.URL.Query().Get("range")
	if rangeParam == "" {
		rangeParam = "month" // Default to month
	}

	rangeType := models.ProgressRange(rangeParam)
	if rangeType != models.ProgressRangeWeek && rangeType != models.ProgressRangeMonth && rangeType != models.ProgressRangeYear {
		respondWithError(w, http.StatusBadRequest, "Invalid range. Must be 'week', 'month', or 'year'")
		return
	}

	progress, err := h.setService.GetExerciseProgress(userID, exerciseID, rangeType)
	if err != nil {
		if err.Error() == "exercise not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, progress)
}

// UpdateExercise handles PUT /exercises/{id}
func (h *ExerciseHandler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	exerciseID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	var req models.ExerciseUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Exercise name is required")
		return
	}

	exercise, err := h.exerciseService.UpdateExercise(userID, exerciseID, &req)
	if err != nil {
		if err.Error() == "exercise not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, exercise)
}

// DeleteExercise handles DELETE /exercises/{id}
func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	vars := mux.Vars(r)
	exerciseID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID")
		return
	}

	err = h.exerciseService.DeleteExercise(userID, exerciseID)
	if err != nil {
		if err.Error() == "exercise not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
