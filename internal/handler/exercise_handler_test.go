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

// mockExerciseService is a mock implementation of ExerciseService
type mockExerciseService struct {
	createFunc func(userID int64, req *models.ExerciseCreateRequest) (*models.ExerciseResponse, error)
	updateFunc func(userID, exerciseID int64, req *models.ExerciseUpdateRequest) (*models.ExerciseResponse, error)
	deleteFunc func(userID, exerciseID int64) error
}

func (m *mockExerciseService) CreateExercise(userID int64, req *models.ExerciseCreateRequest) (*models.ExerciseResponse, error) {
	if m.createFunc != nil {
		return m.createFunc(userID, req)
	}
	return nil, nil
}

func (m *mockExerciseService) GetExercises(userID int64) ([]*models.ExerciseResponse, error) {
	return nil, nil
}

func (m *mockExerciseService) GetExerciseByID(userID, exerciseID int64) (*models.ExerciseResponse, error) {
	return nil, nil
}

func (m *mockExerciseService) UpdateExercise(userID, exerciseID int64, req *models.ExerciseUpdateRequest) (*models.ExerciseResponse, error) {
	if m.updateFunc != nil {
		return m.updateFunc(userID, exerciseID, req)
	}
	return nil, nil
}

func (m *mockExerciseService) DeleteExercise(userID, exerciseID int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(userID, exerciseID)
	}
	return nil
}

// mockSetService is a mock implementation of SetService
type mockSetService struct{}

func (m *mockSetService) CreateSet(userID, workoutID int64, req *models.SetCreateRequest) (*models.SetResponse, error) {
	return nil, nil
}

func (m *mockSetService) GetWorkoutSets(workoutID int64) ([]*models.SetResponse, error) {
	return nil, nil
}

func (m *mockSetService) GetExerciseHistory(userID, exerciseID int64) (*models.ExerciseHistoryResponse, error) {
	return nil, nil
}

func (m *mockSetService) GetExerciseProgress(userID, exerciseID int64, rangeType models.ProgressRange) (*models.ExerciseProgressResponse, error) {
	return nil, nil
}

// TestUpdateExercise tests the UpdateExercise handler
func TestUpdateExercise(t *testing.T) {
	// Test data
	userID := int64(1)
	exerciseID := int64(100)
	originalName := "Bench Press"
	updatedName := "Incline Bench Press"

	// Create a mock exercise service
	mockService := &mockExerciseService{
		createFunc: func(uid int64, req *models.ExerciseCreateRequest) (*models.ExerciseResponse, error) {
			return &models.ExerciseResponse{
				ID:        exerciseID,
				UserID:    uid,
				Name:      req.Name,
				CreatedAt: time.Now(),
			}, nil
		},
		updateFunc: func(uid, eid int64, req *models.ExerciseUpdateRequest) (*models.ExerciseResponse, error) {
			// Verify the user ID and exercise ID match
			if uid != userID {
				t.Errorf("Expected userID %d, got %d", userID, uid)
			}
			if eid != exerciseID {
				t.Errorf("Expected exerciseID %d, got %d", exerciseID, eid)
			}
			// Verify the name was updated
			if req.Name != updatedName {
				t.Errorf("Expected name %s, got %s", updatedName, req.Name)
			}
			return &models.ExerciseResponse{
				ID:        eid, // Same ID as before
				UserID:    uid,
				Name:      req.Name,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	mockSetService := &mockSetService{}

	handler := NewExerciseHandler(mockService, mockSetService)

	t.Run("Create and Update Exercise - ID remains the same", func(t *testing.T) {
		// Step 1: Create an exercise
		createReq := models.ExerciseCreateRequest{
			Name: originalName,
		}
		createBody, _ := json.Marshal(createReq)
		createRequest := httptest.NewRequest("POST", "/exercises", bytes.NewBuffer(createBody))
		createRequest.Header.Set("Content-Type", "application/json")
		// Add user ID to context (simulating middleware)
		ctx := context.WithValue(createRequest.Context(), middleware.UserIDKey, userID)
		createRequest = createRequest.WithContext(ctx)
		createRecorder := httptest.NewRecorder()

		handler.CreateExercise(createRecorder, createRequest)

		if createRecorder.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, createRecorder.Code)
		}

		var createdExercise models.ExerciseResponse
		if err := json.Unmarshal(createRecorder.Body.Bytes(), &createdExercise); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if createdExercise.ID != exerciseID {
			t.Errorf("Expected exercise ID %d, got %d", exerciseID, createdExercise.ID)
		}
		if createdExercise.Name != originalName {
			t.Errorf("Expected name %s, got %s", originalName, createdExercise.Name)
		}

		// Step 2: Update the exercise
		updateReq := models.ExerciseUpdateRequest{
			Name: updatedName,
		}
		updateBody, _ := json.Marshal(updateReq)
		updateRequest := httptest.NewRequest("PUT", "/exercises/100", bytes.NewBuffer(updateBody))
		updateRequest.Header.Set("Content-Type", "application/json")
		// Add user ID to context (simulating middleware)
		ctx = context.WithValue(updateRequest.Context(), middleware.UserIDKey, userID)
		updateRequest = updateRequest.WithContext(ctx)
		// Set mux vars for route parameters
		updateRequest = mux.SetURLVars(updateRequest, map[string]string{"id": "100"})
		updateRecorder := httptest.NewRecorder()

		handler.UpdateExercise(updateRecorder, updateRequest)

		if updateRecorder.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, updateRecorder.Code, updateRecorder.Body.String())
		}

		var updatedExercise models.ExerciseResponse
		if err := json.Unmarshal(updateRecorder.Body.Bytes(), &updatedExercise); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Step 3: Validate that the ID is the same
		if updatedExercise.ID != exerciseID {
			t.Errorf("Expected exercise ID to remain %d, got %d", exerciseID, updatedExercise.ID)
		}

		// Step 4: Validate that the name was updated
		if updatedExercise.Name != updatedName {
			t.Errorf("Expected updated name %s, got %s", updatedName, updatedExercise.Name)
		}

		// Step 5: Validate that the name is different from the original
		if updatedExercise.Name == originalName {
			t.Errorf("Expected name to change from %s, but it remained the same", originalName)
		}
	})

	t.Run("Update Exercise - Invalid Exercise ID", func(t *testing.T) {
		mockServiceWithError := &mockExerciseService{
			updateFunc: func(uid, eid int64, req *models.ExerciseUpdateRequest) (*models.ExerciseResponse, error) {
				return nil, &serviceError{message: "exercise not found"}
			},
		}

		handler := NewExerciseHandler(mockServiceWithError, mockSetService)

		updateReq := models.ExerciseUpdateRequest{
			Name: updatedName,
		}
		updateBody, _ := json.Marshal(updateReq)
		updateRequest := httptest.NewRequest("PUT", "/exercises/999", bytes.NewBuffer(updateBody))
		updateRequest.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(updateRequest.Context(), middleware.UserIDKey, userID)
		updateRequest = updateRequest.WithContext(ctx)
		// Set mux vars for route parameters
		updateRequest = mux.SetURLVars(updateRequest, map[string]string{"id": "999"})
		updateRecorder := httptest.NewRecorder()

		handler.UpdateExercise(updateRecorder, updateRequest)

		if updateRecorder.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, updateRecorder.Code)
		}
	})

	t.Run("Update Exercise - Missing Name", func(t *testing.T) {
		updateReq := models.ExerciseUpdateRequest{
			Name: "",
		}
		updateBody, _ := json.Marshal(updateReq)
		updateRequest := httptest.NewRequest("PUT", "/exercises/100", bytes.NewBuffer(updateBody))
		updateRequest.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(updateRequest.Context(), middleware.UserIDKey, userID)
		updateRequest = updateRequest.WithContext(ctx)
		// Set mux vars for route parameters
		updateRequest = mux.SetURLVars(updateRequest, map[string]string{"id": "100"})
		updateRecorder := httptest.NewRecorder()

		handler.UpdateExercise(updateRecorder, updateRequest)

		if updateRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, updateRecorder.Code)
		}
	})
}

// TestDeleteExercise tests the DeleteExercise handler
func TestDeleteExercise(t *testing.T) {
	userID := int64(1)
	exerciseID := int64(100)

	mockSetService := &mockSetService{}

	t.Run("Delete Exercise - Success", func(t *testing.T) {
		mockServiceWithDelete := &mockExerciseService{
			deleteFunc: func(uid, eid int64) error {
				if uid != userID {
					t.Errorf("Expected userID %d, got %d", userID, uid)
				}
				if eid != exerciseID {
					t.Errorf("Expected exerciseID %d, got %d", exerciseID, eid)
				}
				return nil
			},
		}

		handler := NewExerciseHandler(mockServiceWithDelete, mockSetService)

		deleteRequest := httptest.NewRequest("DELETE", "/exercises/100", nil)
		ctx := context.WithValue(deleteRequest.Context(), middleware.UserIDKey, userID)
		deleteRequest = deleteRequest.WithContext(ctx)
		deleteRequest = mux.SetURLVars(deleteRequest, map[string]string{"id": "100"})
		deleteRecorder := httptest.NewRecorder()

		handler.DeleteExercise(deleteRecorder, deleteRequest)

		if deleteRecorder.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusNoContent, deleteRecorder.Code, deleteRecorder.Body.String())
		}
	})

	t.Run("Delete Exercise - Exercise not found", func(t *testing.T) {
		mockServiceWithError := &mockExerciseService{
			deleteFunc: func(uid, eid int64) error {
				return &serviceError{message: "exercise not found"}
			},
		}

		handler := NewExerciseHandler(mockServiceWithError, mockSetService)

		deleteRequest := httptest.NewRequest("DELETE", "/exercises/999", nil)
		ctx := context.WithValue(deleteRequest.Context(), middleware.UserIDKey, userID)
		deleteRequest = deleteRequest.WithContext(ctx)
		deleteRequest = mux.SetURLVars(deleteRequest, map[string]string{"id": "999"})
		deleteRecorder := httptest.NewRecorder()

		handler.DeleteExercise(deleteRecorder, deleteRequest)

		if deleteRecorder.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, deleteRecorder.Code)
		}
	})

	t.Run("Delete Exercise - Invalid Exercise ID", func(t *testing.T) {
		mockService := &mockExerciseService{}
		handler := NewExerciseHandler(mockService, mockSetService)

		deleteRequest := httptest.NewRequest("DELETE", "/exercises/invalid", nil)
		ctx := context.WithValue(deleteRequest.Context(), middleware.UserIDKey, userID)
		deleteRequest = deleteRequest.WithContext(ctx)
		deleteRequest = mux.SetURLVars(deleteRequest, map[string]string{"id": "invalid"})
		deleteRecorder := httptest.NewRecorder()

		handler.DeleteExercise(deleteRecorder, deleteRequest)

		if deleteRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, deleteRecorder.Code)
		}
	})
}

// serviceError is a simple error type for testing
type serviceError struct {
	message string
}

func (e *serviceError) Error() string {
	return e.message
}
