package service

import (
	"errors"
	"testing"
	"time"

	"phoenix-alliance-be/internal/models"
)

// mockExerciseRepository is a mock implementation of ExerciseRepository
type mockExerciseRepository struct {
	createFunc           func(exercise *models.Exercise) error
	getByIDFunc          func(id int64) (*models.Exercise, error)
	getByUserIDFunc      func(userID int64) ([]*models.Exercise, error)
	getByIDAndUserIDFunc func(id, userID int64) (*models.Exercise, error)
	updateFunc           func(exercise *models.Exercise) error
	deleteFunc           func(id, userID int64) error
}

func (m *mockExerciseRepository) Create(exercise *models.Exercise) error {
	if m.createFunc != nil {
		return m.createFunc(exercise)
	}
	return nil
}

func (m *mockExerciseRepository) GetByID(id int64) (*models.Exercise, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

func (m *mockExerciseRepository) GetByUserID(userID int64) ([]*models.Exercise, error) {
	if m.getByUserIDFunc != nil {
		return m.getByUserIDFunc(userID)
	}
	return nil, nil
}

func (m *mockExerciseRepository) GetByIDAndUserID(id, userID int64) (*models.Exercise, error) {
	if m.getByIDAndUserIDFunc != nil {
		return m.getByIDAndUserIDFunc(id, userID)
	}
	return nil, nil
}

func (m *mockExerciseRepository) Update(exercise *models.Exercise) error {
	if m.updateFunc != nil {
		return m.updateFunc(exercise)
	}
	return nil
}

func (m *mockExerciseRepository) Delete(id, userID int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id, userID)
	}
	return nil
}

// TestUpdateExercise tests the UpdateExercise service method
func TestUpdateExercise(t *testing.T) {
	userID := int64(1)
	exerciseID := int64(100)
	originalName := "Bench Press"
	updatedName := "Incline Bench Press"

	t.Run("Success - Updates exercise and maintains ID", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				if id != exerciseID {
					t.Errorf("Expected exerciseID %d, got %d", exerciseID, id)
				}
				if uid != userID {
					t.Errorf("Expected userID %d, got %d", userID, uid)
				}
				return &models.Exercise{
					ID:        exerciseID,
					UserID:    userID,
					Name:      originalName,
					CreatedAt: time.Now(),
				}, nil
			},
			updateFunc: func(exercise *models.Exercise) error {
				// Validate that the exercise is updated correctly
				if exercise.ID != exerciseID {
					t.Errorf("Expected ID to remain %d, got %d", exerciseID, exercise.ID)
				}
				if exercise.UserID != userID {
					t.Errorf("Expected UserID %d, got %d", userID, exercise.UserID)
				}
				if exercise.Name != updatedName {
					t.Errorf("Expected name to be updated to '%s', got '%s'", updatedName, exercise.Name)
				}
				return nil
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.UpdateExercise(userID, exerciseID, &models.ExerciseUpdateRequest{
			Name: updatedName,
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected result, got nil")
		}
		if result.ID != exerciseID {
			t.Errorf("Expected ID to remain %d, got %d", exerciseID, result.ID)
		}
		if result.UserID != userID {
			t.Errorf("Expected UserID %d, got %d", userID, result.UserID)
		}
		if result.Name != updatedName {
			t.Errorf("Expected name '%s', got '%s'", updatedName, result.Name)
		}
		if result.Name == originalName {
			t.Errorf("Expected name to change from '%s', but it remained the same", originalName)
		}
	})

	t.Run("Error - Exercise not found", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				return nil, errors.New("exercise not found")
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.UpdateExercise(userID, 999, &models.ExerciseUpdateRequest{
			Name: updatedName,
		})

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "exercise not found" {
			t.Errorf("Expected error 'exercise not found', got '%s'", err.Error())
		}
		if result != nil {
			t.Error("Expected nil result on error")
		}
	})

	t.Run("Error - Exercise belongs to different user", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				// Simula que el ejercicio no pertenece al usuario
				return nil, errors.New("exercise not found")
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.UpdateExercise(userID, exerciseID, &models.ExerciseUpdateRequest{
			Name: updatedName,
		})

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "exercise not found" {
			t.Errorf("Expected error 'exercise not found', got '%s'", err.Error())
		}
		if result != nil {
			t.Error("Expected nil result on error")
		}
	})

	t.Run("Error - Database update fails", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				return &models.Exercise{
					ID:        exerciseID,
					UserID:    userID,
					Name:      originalName,
					CreatedAt: time.Now(),
				}, nil
			},
			updateFunc: func(exercise *models.Exercise) error {
				return errors.New("database connection lost")
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.UpdateExercise(userID, exerciseID, &models.ExerciseUpdateRequest{
			Name: updatedName,
		})

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "failed to update exercise" {
			t.Errorf("Expected error 'failed to update exercise', got '%s'", err.Error())
		}
		if result != nil {
			t.Error("Expected nil result on error")
		}
	})
}

// TestCreateExercise tests the CreateExercise service method
func TestCreateExercise(t *testing.T) {
	userID := int64(1)
	exerciseName := "Squat"

	t.Run("Success - Creates exercise with correct data", func(t *testing.T) {
		var createdExercise *models.Exercise
		mockRepo := &mockExerciseRepository{
			createFunc: func(exercise *models.Exercise) error {
				// Validate the exercise data before creation
				if exercise.UserID != userID {
					t.Errorf("Expected UserID %d, got %d", userID, exercise.UserID)
				}
				if exercise.Name != exerciseName {
					t.Errorf("Expected name '%s', got '%s'", exerciseName, exercise.Name)
				}
				if exercise.CreatedAt.IsZero() {
					t.Error("Expected CreatedAt to be set")
				}
				// Simulate database setting the ID
				exercise.ID = 100
				createdExercise = exercise
				return nil
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.CreateExercise(userID, &models.ExerciseCreateRequest{
			Name: exerciseName,
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected result, got nil")
		}
		if result.ID != 100 {
			t.Errorf("Expected ID 100, got %d", result.ID)
		}
		if result.UserID != userID {
			t.Errorf("Expected UserID %d, got %d", userID, result.UserID)
		}
		if result.Name != exerciseName {
			t.Errorf("Expected name '%s', got '%s'", exerciseName, result.Name)
		}
		if createdExercise == nil {
			t.Error("Expected exercise to be created")
		}
	})

	t.Run("Error - Database creation fails", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			createFunc: func(exercise *models.Exercise) error {
				return errors.New("database error")
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.CreateExercise(userID, &models.ExerciseCreateRequest{
			Name: exerciseName,
		})

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "failed to create exercise" {
			t.Errorf("Expected error 'failed to create exercise', got '%s'", err.Error())
		}
		if result != nil {
			t.Error("Expected nil result on error")
		}
	})
}

// TestGetExerciseByID tests the GetExerciseByID service method
func TestGetExerciseByID(t *testing.T) {
	userID := int64(1)
	exerciseID := int64(100)
	exerciseName := "Deadlift"

	t.Run("Success - Returns exercise that belongs to user", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				if id != exerciseID {
					t.Errorf("Expected exerciseID %d, got %d", exerciseID, id)
				}
				if uid != userID {
					t.Errorf("Expected userID %d, got %d", userID, uid)
				}
				return &models.Exercise{
					ID:        exerciseID,
					UserID:    userID,
					Name:      exerciseName,
					CreatedAt: time.Now(),
				}, nil
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.GetExerciseByID(userID, exerciseID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected result, got nil")
		}
		if result.ID != exerciseID {
			t.Errorf("Expected ID %d, got %d", exerciseID, result.ID)
		}
		if result.UserID != userID {
			t.Errorf("Expected UserID %d, got %d", userID, result.UserID)
		}
		if result.Name != exerciseName {
			t.Errorf("Expected name '%s', got '%s'", exerciseName, result.Name)
		}
	})

	t.Run("Error - Exercise not found", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				return nil, errors.New("exercise not found")
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.GetExerciseByID(userID, 999)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "exercise not found" {
			t.Errorf("Expected error 'exercise not found', got '%s'", err.Error())
		}
		if result != nil {
			t.Error("Expected nil result on error")
		}
	})

	t.Run("Error - Exercise belongs to different user", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				// Simula que el ejercicio no pertenece al usuario
				return nil, errors.New("exercise not found")
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.GetExerciseByID(userID, exerciseID)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "exercise not found" {
			t.Errorf("Expected error 'exercise not found', got '%s'", err.Error())
		}
		if result != nil {
			t.Error("Expected nil result on error")
		}
	})
}

// TestGetExercises tests the GetExercises service method
func TestGetExercises(t *testing.T) {
	userID := int64(1)

	t.Run("Success - Returns all exercises for user", func(t *testing.T) {
		expectedExercises := []*models.Exercise{
			{ID: 1, UserID: userID, Name: "Bench Press", CreatedAt: time.Now()},
			{ID: 2, UserID: userID, Name: "Squat", CreatedAt: time.Now()},
			{ID: 3, UserID: userID, Name: "Deadlift", CreatedAt: time.Now()},
		}

		mockRepo := &mockExerciseRepository{
			getByUserIDFunc: func(uid int64) ([]*models.Exercise, error) {
				if uid != userID {
					t.Errorf("Expected userID %d, got %d", userID, uid)
				}
				return expectedExercises, nil
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.GetExercises(userID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected result, got nil")
		}
		if len(result) != len(expectedExercises) {
			t.Errorf("Expected %d exercises, got %d", len(expectedExercises), len(result))
		}
		for i, ex := range result {
			if ex.ID != expectedExercises[i].ID {
				t.Errorf("Expected ID %d, got %d", expectedExercises[i].ID, ex.ID)
			}
			if ex.Name != expectedExercises[i].Name {
				t.Errorf("Expected name '%s', got '%s'", expectedExercises[i].Name, ex.Name)
			}
		}
	})

	t.Run("Success - Returns empty list when user has no exercises", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByUserIDFunc: func(uid int64) ([]*models.Exercise, error) {
				return []*models.Exercise{}, nil
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.GetExercises(userID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected empty slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected 0 exercises, got %d", len(result))
		}
	})

	t.Run("Error - Database query fails", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByUserIDFunc: func(uid int64) ([]*models.Exercise, error) {
				return nil, errors.New("database connection lost")
			},
		}

		service := NewExerciseService(mockRepo)
		result, err := service.GetExercises(userID)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "failed to retrieve exercises" {
			t.Errorf("Expected error 'failed to retrieve exercises', got '%s'", err.Error())
		}
		if result != nil {
			t.Error("Expected nil result on error")
		}
	})
}

// TestDeleteExercise tests the DeleteExercise service method
func TestDeleteExercise(t *testing.T) {
	userID := int64(1)
	exerciseID := int64(100)

	t.Run("Success - Soft deletes exercise", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				if id != exerciseID {
					t.Errorf("Expected exerciseID %d, got %d", exerciseID, id)
				}
				if uid != userID {
					t.Errorf("Expected userID %d, got %d", userID, uid)
				}
				return &models.Exercise{
					ID:        exerciseID,
					UserID:    userID,
					Name:      "Test Exercise",
					CreatedAt: time.Now(),
					DeletedAt: nil,
				}, nil
			},
			deleteFunc: func(id, uid int64) error {
				if id != exerciseID {
					t.Errorf("Expected exerciseID %d, got %d", exerciseID, id)
				}
				if uid != userID {
					t.Errorf("Expected userID %d, got %d", userID, uid)
				}
				return nil
			},
		}

		service := NewExerciseService(mockRepo)
		err := service.DeleteExercise(userID, exerciseID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Error - Exercise not found", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				return nil, errors.New("exercise not found")
			},
		}

		service := NewExerciseService(mockRepo)
		err := service.DeleteExercise(userID, 999)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "exercise not found" {
			t.Errorf("Expected error 'exercise not found', got '%s'", err.Error())
		}
	})

	t.Run("Error - Exercise belongs to different user", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				// Simula que el ejercicio no pertenece al usuario
				return nil, errors.New("exercise not found")
			},
		}

		service := NewExerciseService(mockRepo)
		err := service.DeleteExercise(userID, exerciseID)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "exercise not found" {
			t.Errorf("Expected error 'exercise not found', got '%s'", err.Error())
		}
	})

	t.Run("Error - Database delete fails", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				return &models.Exercise{
					ID:        exerciseID,
					UserID:    userID,
					Name:      "Test Exercise",
					CreatedAt: time.Now(),
				}, nil
			},
			deleteFunc: func(id, uid int64) error {
				return errors.New("database connection lost")
			},
		}

		service := NewExerciseService(mockRepo)
		err := service.DeleteExercise(userID, exerciseID)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "failed to delete exercise" {
			t.Errorf("Expected error 'failed to delete exercise', got '%s'", err.Error())
		}
	})

	t.Run("Error - Exercise already deleted", func(t *testing.T) {
		mockRepo := &mockExerciseRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Exercise, error) {
				return &models.Exercise{
					ID:        exerciseID,
					UserID:    userID,
					Name:      "Test Exercise",
					CreatedAt: time.Now(),
				}, nil
			},
			deleteFunc: func(id, uid int64) error {
				return errors.New("exercise not found")
			},
		}

		service := NewExerciseService(mockRepo)
		err := service.DeleteExercise(userID, exerciseID)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "exercise not found" {
			t.Errorf("Expected error 'exercise not found', got '%s'", err.Error())
		}
	})
}
