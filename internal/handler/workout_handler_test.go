package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"phoenix-alliance-be/internal/middleware"
	"phoenix-alliance-be/internal/models"

	"github.com/gorilla/mux"
)

// mockWorkoutService is a mock implementation of WorkoutService
type mockWorkoutService struct {
	createFunc func(userID int64, req *models.WorkoutCreateRequest) (*models.WorkoutResponse, error)
	updateFunc func(userID, workoutID int64, req *models.WorkoutUpdateRequest) (*models.WorkoutResponse, error)
	deleteFunc func(userID, workoutID int64) error
}

func (m *mockWorkoutService) CreateWorkout(userID int64, req *models.WorkoutCreateRequest) (*models.WorkoutResponse, error) {
	if m.createFunc != nil {
		return m.createFunc(userID, req)
	}
	return nil, nil
}

func (m *mockWorkoutService) GetWorkoutByID(userID, workoutID int64) (*models.WorkoutResponse, error) {
	return nil, nil
}

func (m *mockWorkoutService) GetWorkouts(userID int64) ([]*models.WorkoutResponse, error) {
	return nil, nil
}

func (m *mockWorkoutService) UpdateWorkout(userID, workoutID int64, req *models.WorkoutUpdateRequest) (*models.WorkoutResponse, error) {
	if m.updateFunc != nil {
		return m.updateFunc(userID, workoutID, req)
	}
	return nil, nil
}

func (m *mockWorkoutService) DeleteWorkout(userID, workoutID int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(userID, workoutID)
	}
	return nil
}

// mockSetServiceWorkout is a stub for SetService used in workout handler tests
type mockSetServiceWorkout struct{}

func (m *mockSetServiceWorkout) CreateSet(userID, workoutID int64, req *models.SetCreateRequest) (*models.SetResponse, error) {
	return nil, nil
}

func (m *mockSetServiceWorkout) GetWorkoutSets(workoutID int64) ([]*models.SetResponse, error) {
	return nil, nil
}

func (m *mockSetServiceWorkout) GetExerciseHistory(userID, exerciseID int64) (*models.ExerciseHistoryResponse, error) {
	return nil, nil
}

func (m *mockSetServiceWorkout) GetExerciseProgress(userID, exerciseID int64, rangeType models.ProgressRange) (*models.ExerciseProgressResponse, error) {
	return nil, nil
}

func TestUpdateWorkout(t *testing.T) {
	userID := int64(1)
	workoutID := int64(200)
	originalName := "Leg Day"
	updatedName := "Lower Body"

	mockService := &mockWorkoutService{
		createFunc: func(uid int64, req *models.WorkoutCreateRequest) (*models.WorkoutResponse, error) {
			return &models.WorkoutResponse{
				ID:        workoutID,
				UserID:    uid,
				Name:      req.Name,
				CreatedAt: time.Now(),
			}, nil
		},
		updateFunc: func(uid, wid int64, req *models.WorkoutUpdateRequest) (*models.WorkoutResponse, error) {
			if uid != userID {
				t.Errorf("expected userID %d, got %d", userID, uid)
			}
			if wid != workoutID {
				t.Errorf("expected workoutID %d, got %d", workoutID, wid)
			}
			if req.Name != updatedName {
				t.Errorf("expected name %s, got %s", updatedName, req.Name)
			}
			return &models.WorkoutResponse{
				ID:        wid,
				UserID:    uid,
				Name:      req.Name,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewWorkoutHandler(mockService, &mockSetServiceWorkout{})

	t.Run("Create and Update Workout - ID remains the same", func(t *testing.T) {
		// Create workout
		createReq := models.WorkoutCreateRequest{Name: originalName}
		createBody, _ := json.Marshal(createReq)
		createRequest := httptest.NewRequest("POST", "/workouts", bytes.NewBuffer(createBody))
		createRequest.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(createRequest.Context(), middleware.UserIDKey, userID)
		createRequest = createRequest.WithContext(ctx)
		createRecorder := httptest.NewRecorder()

		handler.CreateWorkout(createRecorder, createRequest)

		if createRecorder.Code != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, createRecorder.Code)
		}

		var createdWorkout models.WorkoutResponse
		if err := json.Unmarshal(createRecorder.Body.Bytes(), &createdWorkout); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if createdWorkout.ID != workoutID {
			t.Errorf("expected workout ID %d, got %d", workoutID, createdWorkout.ID)
		}
		if createdWorkout.Name != originalName {
			t.Errorf("expected name %s, got %s", originalName, createdWorkout.Name)
		}

		// Update workout
		updateReq := models.WorkoutUpdateRequest{Name: updatedName}
		updateBody, _ := json.Marshal(updateReq)
		updateRequest := httptest.NewRequest("PUT", "/workouts/200", bytes.NewBuffer(updateBody))
		updateRequest.Header.Set("Content-Type", "application/json")
		ctx = context.WithValue(updateRequest.Context(), middleware.UserIDKey, userID)
		updateRequest = updateRequest.WithContext(ctx)
		updateRequest = mux.SetURLVars(updateRequest, map[string]string{"id": "200"})
		updateRecorder := httptest.NewRecorder()

		handler.UpdateWorkout(updateRecorder, updateRequest)

		if updateRecorder.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d. Body: %s", http.StatusOK, updateRecorder.Code, updateRecorder.Body.String())
		}

		var updatedWorkout models.WorkoutResponse
		if err := json.Unmarshal(updateRecorder.Body.Bytes(), &updatedWorkout); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if updatedWorkout.ID != workoutID {
			t.Errorf("expected workout ID to remain %d, got %d", workoutID, updatedWorkout.ID)
		}
		if updatedWorkout.Name != updatedName {
			t.Errorf("expected updated name %s, got %s", updatedName, updatedWorkout.Name)
		}
		if updatedWorkout.Name == originalName {
			t.Errorf("expected name to change from %s", originalName)
		}
	})

	t.Run("Update Workout - Workout not found", func(t *testing.T) {
		mockServiceWithError := &mockWorkoutService{
			updateFunc: func(uid, wid int64, req *models.WorkoutUpdateRequest) (*models.WorkoutResponse, error) {
				return nil, &serviceError{message: "workout not found"}
			},
		}
		handler := NewWorkoutHandler(mockServiceWithError, &mockSetServiceWorkout{})

		updateReq := models.WorkoutUpdateRequest{Name: updatedName}
		updateBody, _ := json.Marshal(updateReq)
		updateRequest := httptest.NewRequest("PUT", "/workouts/404", bytes.NewBuffer(updateBody))
		updateRequest.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(updateRequest.Context(), middleware.UserIDKey, userID)
		updateRequest = updateRequest.WithContext(ctx)
		updateRequest = mux.SetURLVars(updateRequest, map[string]string{"id": "404"})
		updateRecorder := httptest.NewRecorder()

		handler.UpdateWorkout(updateRecorder, updateRequest)

		if updateRecorder.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, updateRecorder.Code)
		}
	})

	t.Run("Update Workout - Missing Name", func(t *testing.T) {
		updateReq := models.WorkoutUpdateRequest{Name: ""}
		updateBody, _ := json.Marshal(updateReq)
		updateRequest := httptest.NewRequest("PUT", "/workouts/200", bytes.NewBuffer(updateBody))
		updateRequest.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(updateRequest.Context(), middleware.UserIDKey, userID)
		updateRequest = updateRequest.WithContext(ctx)
		updateRequest = mux.SetURLVars(updateRequest, map[string]string{"id": "200"})
		updateRecorder := httptest.NewRecorder()

		handler.UpdateWorkout(updateRecorder, updateRequest)

		if updateRecorder.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, updateRecorder.Code)
		}
	})
}
