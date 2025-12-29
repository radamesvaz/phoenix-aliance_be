package handler

import (
	"encoding/json"
	"net/http"

	"phoenix-alliance-be/internal/models"
	"phoenix-alliance-be/internal/service"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	userService service.UserService
	config      interface {
		GetJWTSecret() string
		GetJWTExpiry() int
	}
}

// AuthConfig interface for JWT configuration
type AuthConfig interface {
	GetJWTSecret() string
	GetJWTExpiry() int
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService service.UserService, cfg AuthConfig) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		config:      cfg,
	}
}

// Signup handles user registration
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req models.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	if len(req.Password) < 8 {
		respondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			respondWithError(w, http.StatusConflict, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	token, user, err := h.userService.LoginUser(&req, h.config.GetJWTSecret(), h.config.GetJWTExpiry())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
