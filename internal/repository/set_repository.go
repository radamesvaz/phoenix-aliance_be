package repository

import (
	"database/sql"
	"errors"
	"time"

	"phoenix-alliance-be/internal/models"
)

// SetRepository defines the interface for set data operations
type SetRepository interface {
	Create(set *models.Set) error
	GetByID(id int64) (*models.Set, error)
	GetByWorkoutID(workoutID int64) ([]*models.Set, error)
	GetByExerciseID(exerciseID int64) ([]*models.Set, error)
	GetByExerciseIDAndUserID(exerciseID, userID int64) ([]*models.Set, error)
	GetByExerciseIDAndDateRange(exerciseID int64, startDate, endDate time.Time) ([]*models.Set, error)
}

type setRepository struct {
	db *sql.DB
}

// NewSetRepository creates a new set repository
func NewSetRepository(db *sql.DB) SetRepository {
	return &setRepository{db: db}
}

// Create creates a new set
func (r *setRepository) Create(set *models.Set) error {
	query := `
		INSERT INTO sets (workout_id, exercise_id, weight, reps, rest_seconds, notes, rpe, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id_set, workout_id, exercise_id, weight, reps, rest_seconds, notes, rpe, created_at
	`

	err := r.db.QueryRow(
		query,
		set.WorkoutID,
		set.ExerciseID,
		set.Weight,
		set.Reps,
		set.RestSeconds,
		set.Notes,
		set.RPE,
		set.CreatedAt,
	).Scan(
		&set.ID,
		&set.WorkoutID,
		&set.ExerciseID,
		&set.Weight,
		&set.Reps,
		&set.RestSeconds,
		&set.Notes,
		&set.RPE,
		&set.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a set by ID
func (r *setRepository) GetByID(id int64) (*models.Set, error) {
	set := &models.Set{}
	query := `
		SELECT id_set, workout_id, exercise_id, weight, reps, rest_seconds, notes, rpe, created_at 
		FROM sets 
		WHERE id_set = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&set.ID,
		&set.WorkoutID,
		&set.ExerciseID,
		&set.Weight,
		&set.Reps,
		&set.RestSeconds,
		&set.Notes,
		&set.RPE,
		&set.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("set not found")
		}
		return nil, err
	}

	return set, nil
}

// GetByWorkoutID retrieves all sets for a workout
func (r *setRepository) GetByWorkoutID(workoutID int64) ([]*models.Set, error) {
	query := `
		SELECT id_set, workout_id, exercise_id, weight, reps, rest_seconds, notes, rpe, created_at 
		FROM sets 
		WHERE workout_id = $1 
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []*models.Set
	for rows.Next() {
		set := &models.Set{}
		err := rows.Scan(
			&set.ID,
			&set.WorkoutID,
			&set.ExerciseID,
			&set.Weight,
			&set.Reps,
			&set.RestSeconds,
			&set.Notes,
			&set.RPE,
			&set.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}

	return sets, rows.Err()
}

// GetByExerciseID retrieves all sets for an exercise
func (r *setRepository) GetByExerciseID(exerciseID int64) ([]*models.Set, error) {
	query := `
		SELECT id_set, workout_id, exercise_id, weight, reps, rest_seconds, notes, rpe, created_at 
		FROM sets 
		WHERE exercise_id = $1 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []*models.Set
	for rows.Next() {
		set := &models.Set{}
		err := rows.Scan(
			&set.ID,
			&set.WorkoutID,
			&set.ExerciseID,
			&set.Weight,
			&set.Reps,
			&set.RestSeconds,
			&set.Notes,
			&set.RPE,
			&set.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}

	return sets, rows.Err()
}

// GetByExerciseIDAndUserID retrieves all sets for an exercise that belong to a user
func (r *setRepository) GetByExerciseIDAndUserID(exerciseID, userID int64) ([]*models.Set, error) {
	query := `
		SELECT s.id_set, s.workout_id, s.exercise_id, s.weight, s.reps, s.rest_seconds, s.notes, s.rpe, s.created_at 
		FROM sets s
		INNER JOIN exercises e ON s.exercise_id = e.id_exercise
		WHERE s.exercise_id = $1 AND e.user_id = $2
		ORDER BY s.created_at DESC
	`

	rows, err := r.db.Query(query, exerciseID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []*models.Set
	for rows.Next() {
		set := &models.Set{}
		err := rows.Scan(
			&set.ID,
			&set.WorkoutID,
			&set.ExerciseID,
			&set.Weight,
			&set.Reps,
			&set.RestSeconds,
			&set.Notes,
			&set.RPE,
			&set.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}

	return sets, rows.Err()
}

// GetByExerciseIDAndDateRange retrieves sets for an exercise within a date range
func (r *setRepository) GetByExerciseIDAndDateRange(exerciseID int64, startDate, endDate time.Time) ([]*models.Set, error) {
	query := `
		SELECT s.id_set, s.workout_id, s.exercise_id, s.weight, s.reps, s.rest_seconds, s.notes, s.rpe, s.created_at 
		FROM sets s
		WHERE s.exercise_id = $1 AND s.created_at >= $2 AND s.created_at <= $3
		ORDER BY s.created_at ASC
	`

	rows, err := r.db.Query(query, exerciseID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []*models.Set
	for rows.Next() {
		set := &models.Set{}
		err := rows.Scan(
			&set.ID,
			&set.WorkoutID,
			&set.ExerciseID,
			&set.Weight,
			&set.Reps,
			&set.RestSeconds,
			&set.Notes,
			&set.RPE,
			&set.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}

	return sets, rows.Err()
}

