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
		Date:      req.Date,
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

