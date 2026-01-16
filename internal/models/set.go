package models

import (
	"time"
)

// Set represents a training set
type Set struct {
	ID          int64     `json:"id" db:"id_set"`
	WorkoutID   int64     `json:"workout_id" db:"workout_id"`
	ExerciseID  int64     `json:"exercise_id" db:"exercise_id"`
	Weight      float64   `json:"weight" db:"weight"`
	Reps        int       `json:"reps" db:"reps"`
	RestSeconds *int      `json:"rest_seconds,omitempty" db:"rest_seconds"`
	Notes       *string   `json:"notes,omitempty" db:"notes"`
	RPE         *int      `json:"rpe,omitempty" db:"rpe"` // Rate of Perceived Exertion (1-10)
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// SetCreateRequest represents the request body for creating a set
type SetCreateRequest struct {
	ExerciseID  int64   `json:"exercise_id" validate:"required"`
	Weight      float64 `json:"weight" validate:"required,min=0"`
	Reps        int     `json:"reps" validate:"required,min=1"`
	RestSeconds *int    `json:"rest_seconds,omitempty" validate:"omitempty,min=0"`
	Notes       *string `json:"notes,omitempty"`
	RPE         *int    `json:"rpe,omitempty" validate:"omitempty,min=1,max=10"`
}

// SetResponse represents the set data returned in responses
type SetResponse struct {
	ID          int64     `json:"id"`
	WorkoutID   int64     `json:"workout_id"`
	ExerciseID  int64     `json:"exercise_id"`
	Weight      float64   `json:"weight"`
	Reps        int       `json:"reps"`
	RestSeconds *int      `json:"rest_seconds,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	RPE         *int      `json:"rpe,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToResponse converts a Set to SetResponse
func (s *Set) ToResponse() *SetResponse {
	return &SetResponse{
		ID:          s.ID,
		WorkoutID:   s.WorkoutID,
		ExerciseID:  s.ExerciseID,
		Weight:      s.Weight,
		Reps:        s.Reps,
		RestSeconds: s.RestSeconds,
		Notes:       s.Notes,
		RPE:         s.RPE,
		CreatedAt:   s.CreatedAt,
	}
}
