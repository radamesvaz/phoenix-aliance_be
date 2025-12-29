package service

import (
	"errors"
	"time"

	"phoenix-alliance-be/internal/models"
	"phoenix-alliance-be/internal/repository"
)

// SetService defines the interface for set business logic
type SetService interface {
	CreateSet(userID, workoutID int64, req *models.SetCreateRequest) (*models.SetResponse, error)
	GetExerciseHistory(userID, exerciseID int64) (*models.ExerciseHistoryResponse, error)
	GetExerciseProgress(userID, exerciseID int64, rangeType models.ProgressRange) (*models.ExerciseProgressResponse, error)
	GetWorkoutSets(workoutID int64) ([]*models.SetResponse, error)
}

type setService struct {
	setRepo      repository.SetRepository
	exerciseRepo repository.ExerciseRepository
	workoutRepo  repository.WorkoutRepository
}

// NewSetService creates a new set service
func NewSetService(
	setRepo repository.SetRepository,
	exerciseRepo repository.ExerciseRepository,
	workoutRepo repository.WorkoutRepository,
) SetService {
	return &setService{
		setRepo:      setRepo,
		exerciseRepo: exerciseRepo,
		workoutRepo:  workoutRepo,
	}
}

// CreateSet creates a new set for a workout
func (s *setService) CreateSet(userID, workoutID int64, req *models.SetCreateRequest) (*models.SetResponse, error) {
	// Verify workout belongs to user
	_, err := s.workoutRepo.GetByIDAndUserID(workoutID, userID)
	if err != nil {
		return nil, errors.New("workout not found")
	}

	// Verify exercise belongs to user
	_, err = s.exerciseRepo.GetByIDAndUserID(req.ExerciseID, userID)
	if err != nil {
		return nil, errors.New("exercise not found")
	}

	// Create set
	set := &models.Set{
		WorkoutID:   workoutID,
		ExerciseID:  req.ExerciseID,
		Weight:      req.Weight,
		Reps:        req.Reps,
		RestSeconds: req.RestSeconds,
		Notes:       req.Notes,
		RPE:         req.RPE,
		CreatedAt:   time.Now(),
	}

	if err := s.setRepo.Create(set); err != nil {
		return nil, errors.New("failed to create set")
	}

	return set.ToResponse(), nil
}

// GetExerciseHistory retrieves all sets for an exercise with metrics
func (s *setService) GetExerciseHistory(userID, exerciseID int64) (*models.ExerciseHistoryResponse, error) {
	// Verify exercise belongs to user
	exercise, err := s.exerciseRepo.GetByIDAndUserID(exerciseID, userID)
	if err != nil {
		return nil, errors.New("exercise not found")
	}

	// Get all sets for this exercise
	sets, err := s.setRepo.GetByExerciseIDAndUserID(exerciseID, userID)
	if err != nil {
		return nil, errors.New("failed to retrieve sets")
	}

	// Convert to responses
	setResponses := make([]*models.SetResponse, len(sets))
	for i, set := range sets {
		setResponses[i] = set.ToResponse()
	}

	// Calculate metrics
	metrics := calculateMetrics(sets)

	response := &models.ExerciseHistoryResponse{
		ExerciseID:   exerciseID,
		ExerciseName: exercise.Name,
		Sets:         setResponses,
		Metrics:      metrics,
	}

	return response, nil
}

// GetExerciseProgress retrieves progress data for an exercise within a time range
func (s *setService) GetExerciseProgress(userID, exerciseID int64, rangeType models.ProgressRange) (*models.ExerciseProgressResponse, error) {
	// Verify exercise belongs to user
	exercise, err := s.exerciseRepo.GetByIDAndUserID(exerciseID, userID)
	if err != nil {
		return nil, errors.New("exercise not found")
	}

	// Calculate date range
	now := time.Now()
	var startDate, endDate time.Time

	switch rangeType {
	case models.ProgressRangeWeek:
		startDate = now.AddDate(0, 0, -7)
		endDate = now
	case models.ProgressRangeMonth:
		startDate = now.AddDate(0, -1, 0)
		endDate = now
	case models.ProgressRangeYear:
		startDate = now.AddDate(-1, 0, 0)
		endDate = now
	default:
		return nil, errors.New("invalid range type")
	}

	// Get sets within date range
	sets, err := s.setRepo.GetByExerciseIDAndDateRange(exerciseID, startDate, endDate)
	if err != nil {
		return nil, errors.New("failed to retrieve sets")
	}

	// Group sets by date and calculate data points
	dataPoints := groupSetsByDate(sets)

	// Calculate summary metrics
	summary := calculateMetrics(sets)

	response := &models.ExerciseProgressResponse{
		ExerciseID:   exerciseID,
		ExerciseName: exercise.Name,
		Range:        string(rangeType),
		StartDate:    startDate,
		EndDate:      endDate,
		DataPoints:   dataPoints,
		Summary:      summary,
	}

	return response, nil
}

// calculateMetrics calculates aggregated metrics from sets
func calculateMetrics(sets []*models.Set) *models.ExerciseMetrics {
	if len(sets) == 0 {
		return nil
	}

	metrics := &models.ExerciseMetrics{
		TotalSets: len(sets),
	}

	var totalVolume float64
	var totalWeight float64
	var totalReps int
	var totalRest int
	var restCount int
	var totalRPE int
	var rpeCount int
	var maxWeight float64
	var maxReps int
	var firstDate *time.Time
	var lastDate *time.Time

	for _, set := range sets {
		volume := set.Weight * float64(set.Reps)
		totalVolume += volume
		totalWeight += set.Weight
		totalReps += set.Reps

		if set.Weight > maxWeight {
			maxWeight = set.Weight
		}
		if set.Reps > maxReps {
			maxReps = set.Reps
		}

		if set.RestSeconds != nil {
			totalRest += *set.RestSeconds
			restCount++
		}

		if set.RPE != nil {
			totalRPE += *set.RPE
			rpeCount++
		}

		if firstDate == nil || set.CreatedAt.Before(*firstDate) {
			firstDate = &set.CreatedAt
		}
		if lastDate == nil || set.CreatedAt.After(*lastDate) {
			lastDate = &set.CreatedAt
		}
	}

	metrics.TotalVolume = totalVolume
	metrics.MaxWeight = maxWeight
	metrics.MaxReps = maxReps
	metrics.AverageWeight = totalWeight / float64(len(sets))
	metrics.AverageReps = float64(totalReps) / float64(len(sets))

	if restCount > 0 {
		avgRest := float64(totalRest) / float64(restCount)
		metrics.AverageRest = &avgRest
	}

	if rpeCount > 0 {
		avgRPE := float64(totalRPE) / float64(rpeCount)
		metrics.AverageRPE = &avgRPE
	}

	metrics.FirstRecordedAt = firstDate
	metrics.LastRecordedAt = lastDate

	return metrics
}

// groupSetsByDate groups sets by date and creates data points
func groupSetsByDate(sets []*models.Set) []models.ProgressDataPoint {
	if len(sets) == 0 {
		return []models.ProgressDataPoint{}
	}

	// Group sets by date
	dateMap := make(map[string][]*models.Set)
	for _, set := range sets {
		dateKey := set.CreatedAt.Format("2006-01-02")
		dateMap[dateKey] = append(dateMap[dateKey], set)
	}

	// Create data points
	var dataPoints []models.ProgressDataPoint
	for dateKey, daySets := range dateMap {
		date, _ := time.Parse("2006-01-02", dateKey)

		var totalVolume float64
		var maxWeight float64
		var totalRPE int
		var rpeCount int

		for _, set := range daySets {
			volume := set.Weight * float64(set.Reps)
			totalVolume += volume
			if set.Weight > maxWeight {
				maxWeight = set.Weight
			}
			if set.RPE != nil {
				totalRPE += *set.RPE
				rpeCount++
			}
		}

		dp := models.ProgressDataPoint{
			Date:        date,
			TotalVolume: totalVolume,
			MaxWeight:   maxWeight,
			TotalSets:   len(daySets),
		}

		if rpeCount > 0 {
			avgRPE := float64(totalRPE) / float64(rpeCount)
			dp.AverageRPE = &avgRPE
		}

		dataPoints = append(dataPoints, dp)
	}

	return dataPoints
}

// GetWorkoutSets retrieves all sets for a workout
func (s *setService) GetWorkoutSets(workoutID int64) ([]*models.SetResponse, error) {
	sets, err := s.setRepo.GetByWorkoutID(workoutID)
	if err != nil {
		return nil, errors.New("failed to retrieve sets")
	}

	responses := make([]*models.SetResponse, len(sets))
	for i, set := range sets {
		responses[i] = set.ToResponse()
	}

	return responses, nil
}
