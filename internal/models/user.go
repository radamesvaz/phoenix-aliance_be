package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        int64     `json:"id" db:"id_user"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Never return password in JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserCreateRequest represents the request body for creating a user
type UserCreateRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// UserLoginRequest represents the request body for login
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents the user data returned in responses
type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

