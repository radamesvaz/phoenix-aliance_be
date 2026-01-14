package service

import (
	"errors"
	"testing"
	"time"

	"phoenix-alliance-be/internal/models"
)

// mockWorkoutRepository is a mock implementation of WorkoutRepository
type mockWorkoutRepository struct {
	createFunc          func(workout *models.Workout) error
	getByIDFunc         func(id int64) (*models.Workout, error)
	getByIDAndUserIDFunc func(id, userID int64) (*models.Workout, error)
	getByUserIDFunc     func(userID int64) ([]*models.Workout, error)
	deleteFunc          func(id, userID int64) error
}

func (m *mockWorkoutRepository) Create(workout *models.Workout) error {
	if m.createFunc != nil {
		return m.createFunc(workout)
	}
	return nil
}

func (m *mockWorkoutRepository) GetByID(id int64) (*models.Workout, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

func (m *mockWorkoutRepository) GetByIDAndUserID(id, userID int64) (*models.Workout, error) {
	if m.getByIDAndUserIDFunc != nil {
		return m.getByIDAndUserIDFunc(id, userID)
	}
	return nil, nil
}

func (m *mockWorkoutRepository) GetByUserID(userID int64) ([]*models.Workout, error) {
	if m.getByUserIDFunc != nil {
		return m.getByUserIDFunc(userID)
	}
	return nil, nil
}

func (m *mockWorkoutRepository) Delete(id, userID int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id, userID)
	}
	return nil
}

func TestCreateWorkout(t *testing.T) {
	userID := int64(1)
	name := "Leg Day"

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockWorkoutRepository{
			createFunc: func(workout *models.Workout) error {
				if workout.UserID != userID {
					t.Errorf("expected userID %d, got %d", userID, workout.UserID)
				}
				if workout.Name != name {
					t.Errorf("expected name %s, got %s", name, workout.Name)
				}
				if workout.CreatedAt.IsZero() {
					t.Error("expected CreatedAt to be set")
				}
				workout.ID = 10
				return nil
			},
		}

		svc := NewWorkoutService(mockRepo)
		res, err := svc.CreateWorkout(userID, &models.WorkoutCreateRequest{Name: name})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res.ID != 10 {
			t.Errorf("expected ID 10, got %d", res.ID)
		}
		if res.UserID != userID {
			t.Errorf("expected userID %d, got %d", userID, res.UserID)
		}
		if res.Name != name {
			t.Errorf("expected name %s, got %s", name, res.Name)
		}
	})

	t.Run("create error", func(t *testing.T) {
		mockRepo := &mockWorkoutRepository{
			createFunc: func(workout *models.Workout) error {
				return errors.New("db error")
			},
		}

		svc := NewWorkoutService(mockRepo)
		res, err := svc.CreateWorkout(userID, &models.WorkoutCreateRequest{Name: name})

		if err == nil || err.Error() != "failed to create workout" {
			t.Fatalf("expected failed to create workout error, got %v", err)
		}
		if res != nil {
			t.Fatalf("expected nil result, got %+v", res)
		}
	})
}

func TestGetWorkoutByID(t *testing.T) {
	userID := int64(1)
	workoutID := int64(20)

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockWorkoutRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Workout, error) {
				if id != workoutID || uid != userID {
					t.Fatalf("expected ids (%d,%d) got (%d,%d)", workoutID, userID, id, uid)
				}
				return &models.Workout{
					ID:        workoutID,
					UserID:    userID,
					Name:      "Push",
					CreatedAt: time.Now(),
				}, nil
			},
		}

		svc := NewWorkoutService(mockRepo)
		res, err := svc.GetWorkoutByID(userID, workoutID)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res.ID != workoutID {
			t.Errorf("expected id %d, got %d", workoutID, res.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := &mockWorkoutRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Workout, error) {
				return nil, errors.New("workout not found")
			},
		}
		svc := NewWorkoutService(mockRepo)
		res, err := svc.GetWorkoutByID(userID, workoutID)

		if err == nil || err.Error() != "workout not found" {
			t.Fatalf("expected workout not found error, got %v", err)
		}
		if res != nil {
			t.Fatalf("expected nil result, got %+v", res)
		}
	})
}

func TestGetWorkouts(t *testing.T) {
	userID := int64(1)

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockWorkoutRepository{
			getByUserIDFunc: func(uid int64) ([]*models.Workout, error) {
				if uid != userID {
					t.Fatalf("expected user %d, got %d", userID, uid)
				}
				return []*models.Workout{
					{ID: 1, UserID: uid, Name: "A", CreatedAt: time.Now()},
					{ID: 2, UserID: uid, Name: "B", CreatedAt: time.Now()},
				}, nil
			},
		}

		svc := NewWorkoutService(mockRepo)
		res, err := svc.GetWorkouts(userID)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(res) != 2 {
			t.Fatalf("expected 2 workouts, got %d", len(res))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := &mockWorkoutRepository{
			getByUserIDFunc: func(uid int64) ([]*models.Workout, error) {
				return []*models.Workout{}, nil
			},
		}
		svc := NewWorkoutService(mockRepo)
		res, err := svc.GetWorkouts(userID)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res == nil {
			t.Fatalf("expected empty slice, got nil")
		}
		if len(res) != 0 {
			t.Fatalf("expected 0 workouts, got %d", len(res))
		}
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := &mockWorkoutRepository{
			getByUserIDFunc: func(uid int64) ([]*models.Workout, error) {
				return nil, errors.New("db down")
			},
		}
		svc := NewWorkoutService(mockRepo)
		res, err := svc.GetWorkouts(userID)

		if err == nil || err.Error() != "failed to retrieve workouts" {
			t.Fatalf("expected failed to retrieve workouts, got %v", err)
		}
		if res != nil {
			t.Fatalf("expected nil result, got %+v", res)
		}
	})
}

func TestDeleteWorkout(t *testing.T) {
	userID := int64(1)
	workoutID := int64(33)

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockWorkoutRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Workout, error) {
				return &models.Workout{ID: id, UserID: uid, Name: "Test", CreatedAt: time.Now()}, nil
			},
			deleteFunc: func(id, uid int64) error {
				if id != workoutID || uid != userID {
					t.Fatalf("unexpected ids: %d %d", id, uid)
				}
				return nil
			},
		}

		svc := NewWorkoutService(mockRepo)
		if err := svc.DeleteWorkout(userID, workoutID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("not found on precheck", func(t *testing.T) {
		mockRepo := &mockWorkoutRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Workout, error) {
				return nil, errors.New("workout not found")
			},
		}
		svc := NewWorkoutService(mockRepo)
		err := svc.DeleteWorkout(userID, workoutID)

		if err == nil || err.Error() != "workout not found" {
			t.Fatalf("expected workout not found, got %v", err)
		}
	})

	t.Run("not found on delete", func(t *testing.T) {
		mockRepo := &mockWorkoutRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Workout, error) {
				return &models.Workout{ID: id, UserID: uid}, nil
			},
			deleteFunc: func(id, uid int64) error {
				return errors.New("workout not found")
			},
		}
		svc := NewWorkoutService(mockRepo)
		err := svc.DeleteWorkout(userID, workoutID)

		if err == nil || err.Error() != "workout not found" {
			t.Fatalf("expected workout not found, got %v", err)
		}
	})

	t.Run("delete fails", func(t *testing.T) {
		mockRepo := &mockWorkoutRepository{
			getByIDAndUserIDFunc: func(id, uid int64) (*models.Workout, error) {
				return &models.Workout{ID: id, UserID: uid}, nil
			},
			deleteFunc: func(id, uid int64) error {
				return errors.New("db down")
			},
		}
		svc := NewWorkoutService(mockRepo)
		err := svc.DeleteWorkout(userID, workoutID)

		if err == nil || err.Error() != "failed to delete workout" {
			t.Fatalf("expected failed to delete workout, got %v", err)
		}
	})
}
