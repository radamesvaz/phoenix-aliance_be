package service

import (
    "errors"
    "time"

    "phoenix-alliance-be/internal/models"
    "phoenix-alliance-be/internal/repository"
)

// WorkoutService defines the interface for workout business logic
type WorkoutService interface {
    CreateWorkout(userID int64, req *models.WorkoutCreateRequest) (*models.WorkoutResponse, error)
    GetWorkoutByID(userID, workoutID int64) (*models.WorkoutResponse, error)
    GetWorkouts(userID int64) ([]*models.WorkoutResponse, error)
    DeleteWorkout(userID, workoutID int64) error
}

type workoutService struct {
    workoutRepo repository.WorkoutRepository
}

// NewWorkoutService creates a new workout service
func NewWorkoutService(workoutRepo repository.WorkoutRepository) WorkoutService {
    return &workoutService{workoutRepo: workoutRepo}
}

// CreateWorkout creates a new workout for a user
func (s *workoutService) CreateWorkout(userID int64, req *models.WorkoutCreateRequest) (*models.WorkoutResponse, error) {
    workout := &models.Workout{
        UserID:    userID,
        Name:      req.Name,
        CreatedAt: time.Now(),
    }

    if err := s.workoutRepo.Create(workout); err != nil {
        return nil, errors.New("failed to create workout")
    }

    return workout.ToResponse(), nil
}

// GetWorkoutByID retrieves a workout by ID for a user
func (s *workoutService) GetWorkoutByID(userID, workoutID int64) (*models.WorkoutResponse, error) {
    workout, err := s.workoutRepo.GetByIDAndUserID(workoutID, userID)
    if err != nil {
        return nil, errors.New("workout not found")
    }

    return workout.ToResponse(), nil
}

// GetWorkouts retrieves all workouts for a user
func (s *workoutService) GetWorkouts(userID int64) ([]*models.WorkoutResponse, error) {
    workouts, err := s.workoutRepo.GetByUserID(userID)
    if err != nil {
        return nil, errors.New("failed to retrieve workouts")
    }
    responses := make([]*models.WorkoutResponse, len(workouts))
    for i, workout := range workouts {
        responses[i] = workout.ToResponse()
    }
    return responses, nil
}

// DeleteWorkout performs a soft delete on a workout for a user
func (s *workoutService) DeleteWorkout(userID, workoutID int64) error {
    // Verify workout exists and belongs to user
    if _, err := s.workoutRepo.GetByIDAndUserID(workoutID, userID); err != nil {
        return errors.New("workout not found")
    }

    // Soft delete
    if err := s.workoutRepo.Delete(workoutID, userID); err != nil {
        if err.Error() == "workout not found" {
            return errors.New("workout not found")
        }
        return errors.New("failed to delete workout")
    }

    return nil
}