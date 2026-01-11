package models

import (
	"time"
)

// Exercise represents an exercise created by a user
type Exercise struct {
	ID        int64      `json:"id" db:"id_exercise"`
	UserID    int64      `json:"user_id" db:"user_id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// ExerciseCreateRequest represents the request body for creating an exercise
type ExerciseCreateRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

// ExerciseUpdateRequest represents the request body for updating an exercise
type ExerciseUpdateRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

// ExerciseResponse represents the exercise data returned in responses
// @TODO: Agregar mas campos al ejercicio:
// Musculo principal
// Musculo secundario
// Tipo de ejercicio (peso, cardio, flexibilidad, etc.)
type ExerciseResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts an Exercise to ExerciseResponse
func (e *Exercise) ToResponse() *ExerciseResponse {
	return &ExerciseResponse{
		ID:        e.ID,
		UserID:    e.UserID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
	}
}
