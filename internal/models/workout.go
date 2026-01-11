package models

import (
	"time"
)

// Workout represents a workout session
type Workout struct {
	ID        int64     `json:"id" db:"id_workout"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// WorkoutCreateRequest represents the request body for creating a workout
type WorkoutCreateRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

// WorkoutResponse represents the workout data returned in responses
type WorkoutResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a Workout to WorkoutResponse
func (w *Workout) ToResponse() *WorkoutResponse {
	return &WorkoutResponse{
		ID:        w.ID,
		UserID:    w.UserID,
		Name:      w.Name,
		CreatedAt: w.CreatedAt,
	}
}

