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
	Update(exercise *models.Exercise) error
	Delete(id, userID int64) error
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

// GetByID retrieves an exercise by ID (only non-deleted)
func (r *exerciseRepository) GetByID(id int64) (*models.Exercise, error) {
	exercise := &models.Exercise{}
	query := `SELECT id_exercise, user_id, name, created_at, deleted_at FROM exercises WHERE id_exercise = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id).Scan(
		&exercise.ID,
		&exercise.UserID,
		&exercise.Name,
		&exercise.CreatedAt,
		&exercise.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("exercise not found")
		}
		return nil, err
	}

	return exercise, nil
}

// GetByUserID retrieves all exercises for a user (only non-deleted)
func (r *exerciseRepository) GetByUserID(userID int64) ([]*models.Exercise, error) {
	query := `
		SELECT id_exercise, user_id, name, created_at, deleted_at
		FROM exercises 
		WHERE user_id = $1 AND deleted_at IS NULL
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
			&exercise.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	return exercises, rows.Err()
}

// GetByIDAndUserID retrieves an exercise by ID and ensures it belongs to the user (only non-deleted)
func (r *exerciseRepository) GetByIDAndUserID(id, userID int64) (*models.Exercise, error) {
	exercise := &models.Exercise{}
	query := `SELECT id_exercise, user_id, name, created_at, deleted_at FROM exercises WHERE id_exercise = $1 AND user_id = $2 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id, userID).Scan(
		&exercise.ID,
		&exercise.UserID,
		&exercise.Name,
		&exercise.CreatedAt,
		&exercise.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("exercise not found")
		}
		return nil, err
	}

	return exercise, nil
}

// Update updates an existing exercise (only non-deleted)
func (r *exerciseRepository) Update(exercise *models.Exercise) error {
	query := `
		UPDATE exercises 
		SET name = $1
		WHERE id_exercise = $2 AND user_id = $3 AND deleted_at IS NULL
		RETURNING id_exercise, user_id, name, created_at, deleted_at
	`

	err := r.db.QueryRow(
		query,
		exercise.Name,
		exercise.ID,
		exercise.UserID,
	).Scan(
		&exercise.ID,
		&exercise.UserID,
		&exercise.Name,
		&exercise.CreatedAt,
		&exercise.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("exercise not found")
		}
		return err
	}

	return nil
}

// Delete performs a soft delete on an exercise
func (r *exerciseRepository) Delete(id, userID int64) error {
	query := `
		UPDATE exercises 
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id_exercise = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("exercise not found")
	}

	return nil
}
