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
		RETURNING id_workout, user_id, name, created_at
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
	)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a workout by ID
func (r *workoutRepository) GetByID(id int64) (*models.Workout, error) {
	workout := &models.Workout{}
	query := `SELECT id_workout, user_id, name, created_at FROM workouts WHERE id_workout = $1`

	err := r.db.QueryRow(query, id).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.CreatedAt,
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
	query := `SELECT id_workout, user_id, name, created_at FROM workouts WHERE id_workout = $1 AND user_id = $2`

	err := r.db.QueryRow(query, id, userID).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.CreatedAt,
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
		SELECT id_workout, user_id, name, created_at 
		FROM workouts 
		WHERE user_id = $1 
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
		)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, workout)
	}

	return workouts, rows.Err()
}
