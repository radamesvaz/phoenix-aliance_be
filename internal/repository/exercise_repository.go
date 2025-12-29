package repository

import (
	"database/sql"
	"errors"

	"phoenix-alliance-be/internal/models"
)

// ExerciseRepository defines the interface for exercise data operations
type ExerciseRepository interface {
	Create(exercise *models.Exercise) error
	GetByID(id int64) (*models.Exercise, error)
	GetByUserID(userID int64) ([]*models.Exercise, error)
	GetByIDAndUserID(id, userID int64) (*models.Exercise, error)
}

type exerciseRepository struct {
	db *sql.DB
}

// NewExerciseRepository creates a new exercise repository
func NewExerciseRepository(db *sql.DB) ExerciseRepository {
	return &exerciseRepository{db: db}
}

// Create creates a new exercise
func (r *exerciseRepository) Create(exercise *models.Exercise) error {
	query := `
		INSERT INTO exercises (user_id, name, created_at)
		VALUES ($1, $2, $3)
		RETURNING id_exercise, user_id, name, created_at
	`

	err := r.db.QueryRow(
		query,
		exercise.UserID,
		exercise.Name,
		exercise.CreatedAt,
	).Scan(
		&exercise.ID,
		&exercise.UserID,
		&exercise.Name,
		&exercise.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves an exercise by ID
func (r *exerciseRepository) GetByID(id int64) (*models.Exercise, error) {
	exercise := &models.Exercise{}
	query := `SELECT id_exercise, user_id, name, created_at FROM exercises WHERE id_exercise = $1`

	err := r.db.QueryRow(query, id).Scan(
		&exercise.ID,
		&exercise.UserID,
		&exercise.Name,
		&exercise.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("exercise not found")
		}
		return nil, err
	}

	return exercise, nil
}

// GetByUserID retrieves all exercises for a user
func (r *exerciseRepository) GetByUserID(userID int64) ([]*models.Exercise, error) {
	query := `
		SELECT id_exercise, user_id, name, created_at 
		FROM exercises 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*models.Exercise
	for rows.Next() {
		exercise := &models.Exercise{}
		err := rows.Scan(
			&exercise.ID,
			&exercise.UserID,
			&exercise.Name,
			&exercise.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	return exercises, rows.Err()
}

// GetByIDAndUserID retrieves an exercise by ID and ensures it belongs to the user
func (r *exerciseRepository) GetByIDAndUserID(id, userID int64) (*models.Exercise, error) {
	exercise := &models.Exercise{}
	query := `SELECT id_exercise, user_id, name, created_at FROM exercises WHERE id_exercise = $1 AND user_id = $2`

	err := r.db.QueryRow(query, id, userID).Scan(
		&exercise.ID,
		&exercise.UserID,
		&exercise.Name,
		&exercise.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("exercise not found")
		}
		return nil, err
	}

	return exercise, nil
}
