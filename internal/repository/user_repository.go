package repository

import (
	"database/sql"
	"errors"

	"phoenix-alliance-be/internal/models"

)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, password, created_at)
		VALUES ($1, $2, $3)
		RETURNING id_user, email, created_at
	`

	err := r.db.QueryRow(
		query,
		user.Email,
		user.Password,
		user.CreatedAt,
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id_user, email, password, created_at FROM users WHERE id_user = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id_user, email, password, created_at FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
