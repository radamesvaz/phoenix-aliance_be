package service

import (
	"errors"
	"time"

	"phoenix-alliance-be/internal/auth"
	"phoenix-alliance-be/internal/models"
	"phoenix-alliance-be/internal/repository"
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(req *models.UserCreateRequest) (*models.UserResponse, error)
	LoginUser(req *models.UserLoginRequest, jwtSecret string, jwtExpiry int) (string, *models.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// CreateUser creates a new user
func (s *userService) CreateUser(req *models.UserCreateRequest) (*models.UserResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user.ToResponse(), nil
}

// LoginUser authenticates a user and returns a JWT token
func (s *userService) LoginUser(req *models.UserLoginRequest, jwtSecret string, jwtExpiry int) (string, *models.UserResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		return "", nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, jwtSecret, jwtExpiry)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user.ToResponse(), nil
}

