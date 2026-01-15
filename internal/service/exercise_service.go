package service

import (
	"errors"
	"time"

	"phoenix-alliance-be/internal/models"
	"phoenix-alliance-be/internal/repository"
)

// ExerciseService defines the interface for exercise business logic
type ExerciseService interface {
	CreateExercise(userID int64, req *models.ExerciseCreateRequest) (*models.ExerciseResponse, error)
	GetExercises(userID int64) ([]*models.ExerciseResponse, error)
	GetExerciseByID(userID, exerciseID int64) (*models.ExerciseResponse, error)
	UpdateExercise(userID, exerciseID int64, req *models.ExerciseUpdateRequest) (*models.ExerciseResponse, error)
	DeleteExercise(userID, exerciseID int64) error
}

type exerciseService struct {
	exerciseRepo repository.ExerciseRepository
}

// NewExerciseService creates a new exercise service
func NewExerciseService(exerciseRepo repository.ExerciseRepository) ExerciseService {
	return &exerciseService{exerciseRepo: exerciseRepo}
}

// CreateExercise creates a new exercise for a user
func (s *exerciseService) CreateExercise(userID int64, req *models.ExerciseCreateRequest) (*models.ExerciseResponse, error) {
	exercise := &models.Exercise{
		UserID:    userID,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := s.exerciseRepo.Create(exercise); err != nil {
		return nil, errors.New("failed to create exercise")
	}

	return exercise.ToResponse(), nil
}

// GetExercises retrieves all exercises for a user
func (s *exerciseService) GetExercises(userID int64) ([]*models.ExerciseResponse, error) {
	exercises, err := s.exerciseRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve exercises")
	}

	responses := make([]*models.ExerciseResponse, len(exercises))
	for i, exercise := range exercises {
		responses[i] = exercise.ToResponse()
	}

	return responses, nil
}

// GetExerciseByID retrieves an exercise by ID for a user
func (s *exerciseService) GetExerciseByID(userID, exerciseID int64) (*models.ExerciseResponse, error) {
	exercise, err := s.exerciseRepo.GetByIDAndUserID(exerciseID, userID)
	if err != nil {
		return nil, errors.New("exercise not found")
	}

	return exercise.ToResponse(), nil
}

// UpdateExercise updates an existing exercise for a user
func (s *exerciseService) UpdateExercise(userID, exerciseID int64, req *models.ExerciseUpdateRequest) (*models.ExerciseResponse, error) {
	// First, verify the exercise exists and belongs to the user
	exercise, err := s.exerciseRepo.GetByIDAndUserID(exerciseID, userID)
	if err != nil {
		return nil, errors.New("exercise not found")
	}

	// Update the exercise name
	exercise.Name = req.Name

	// Save the updated exercise
	if err := s.exerciseRepo.Update(exercise); err != nil {
		return nil, errors.New("failed to update exercise")
	}

	return exercise.ToResponse(), nil
}

// DeleteExercise performs a soft delete on an exercise for a user
func (s *exerciseService) DeleteExercise(userID, exerciseID int64) error {
	// First, verify the exercise exists and belongs to the user
	_, err := s.exerciseRepo.GetByIDAndUserID(exerciseID, userID)
	if err != nil {
		return errors.New("exercise not found")
	}	// Perform soft delete
	if err := s.exerciseRepo.Delete(exerciseID, userID); err != nil {
		if err.Error() == "exercise not found" {
			return errors.New("exercise not found")
		}
		return errors.New("failed to delete exercise")
	}
	return nil
}
