package repository

import (
	"database/sql"
	"errors"

	"phoenix-alliance-be/internal/models"
)

// WorkoutRepository defines the interface for workout data operations
type WorkoutRepository interface {
	Create(workout *models.Workout) error
	GetByID(id int64) (*models.Workout, error)
	GetByIDAndUserID(id, userID int64) (*models.Workout, error)
	GetByUserID(userID int64) ([]*models.Workout, error)
	Update(workout *models.Workout) error
	Delete(id, userID int64) error
}

type workoutRepository struct {
	db *sql.DB
}

// NewWorkoutRepository creates a new workout repository
func NewWorkoutRepository(db *sql.DB) WorkoutRepository {
	return &workoutRepository{db: db}
}

// Create creates a new workout
func (r *workoutRepository) Create(workout *models.Workout) error {
	query := `
		INSERT INTO workouts (user_id, name, created_at)
		VALUES ($1, $2, $3)
		RETURNING id_workout, user_id, name, created_at, deleted_at
	`

	err := r.db.QueryRow(
		query,
		workout.UserID,
		workout.Name,
		workout.CreatedAt,
	).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.CreatedAt,
		&workout.DeletedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a workout by ID
func (r *workoutRepository) GetByID(id int64) (*models.Workout, error) {
	workout := &models.Workout{}
	query := `SELECT id_workout, user_id, name, created_at, deleted_at FROM workouts WHERE id_workout = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.CreatedAt,
		&workout.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("workout not found")
		}
		return nil, err
	}

	return workout, nil
}

// GetByIDAndUserID retrieves a workout by ID and ensures it belongs to the user
func (r *workoutRepository) GetByIDAndUserID(id, userID int64) (*models.Workout, error) {
	workout := &models.Workout{}
	query := `SELECT id_workout, user_id, name, created_at, deleted_at FROM workouts WHERE id_workout = $1 AND user_id = $2 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id, userID).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.CreatedAt,
		&workout.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("workout not found")
		}
		return nil, err
	}

	return workout, nil
}

// GetByUserID retrieves all workouts for a user
func (r *workoutRepository) GetByUserID(userID int64) ([]*models.Workout, error) {
	query := `
		SELECT id_workout, user_id, name, created_at, deleted_at
		FROM workouts
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []*models.Workout
	for rows.Next() {
		workout := &models.Workout{}
		err := rows.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.Name,
			&workout.CreatedAt,
			&workout.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, workout)
	}

	return workouts, rows.Err()
}

// Update updates an existing workout (only non-deleted)
func (r *workoutRepository) Update(workout *models.Workout) error {
	query := `
		UPDATE workouts
		SET name = $1
		WHERE id_workout = $2 AND user_id = $3 AND deleted_at IS NULL
		RETURNING id_workout, user_id, name, created_at, deleted_at
	`

	err := r.db.QueryRow(
		query,
		workout.Name,
		workout.ID,
		workout.UserID,
	).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.CreatedAt,
		&workout.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("workout not found")
		}
		return err
	}

	return nil
}

// Delete performs a soft delete on a workout for a user
func (r *workoutRepository) Delete(id, userID int64) error {
	query := `
		UPDATE workouts
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id_workout = $1 AND user_id = $2 AND deleted_at IS NULL
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
		return errors.New("workout not found")
	}

	return nil
}
