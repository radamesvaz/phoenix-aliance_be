package models

import "time"

// ExerciseHistoryResponse represents the history of sets for an exercise
type ExerciseHistoryResponse struct {
	ExerciseID   int64          `json:"exercise_id"`
	ExerciseName string         `json:"exercise_name"`
	Sets         []*SetResponse `json:"sets"`
	Metrics      *ExerciseMetrics `json:"metrics,omitempty"`
}

// ExerciseMetrics represents aggregated metrics for an exercise
type ExerciseMetrics struct {
	TotalSets        int     `json:"total_sets"`
	TotalVolume      float64 `json:"total_volume"`      // Sum of (weight * reps)
	MaxWeight        float64 `json:"max_weight"`
	MaxReps          int     `json:"max_reps"`
	AverageWeight    float64 `json:"average_weight"`
	AverageReps      float64 `json:"average_reps"`
	AverageRest      *float64 `json:"average_rest,omitempty"`
	AverageRPE       *float64 `json:"average_rpe,omitempty"`
	FirstRecordedAt  *time.Time `json:"first_recorded_at,omitempty"`
	LastRecordedAt   *time.Time `json:"last_recorded_at,omitempty"`
}

// ProgressRange represents the time range for progress queries
type ProgressRange string

const (
	ProgressRangeWeek  ProgressRange = "week"
	ProgressRangeMonth ProgressRange = "month"
	ProgressRangeYear  ProgressRange = "year"
)

// ExerciseProgressResponse represents progress data for a specific time range
type ExerciseProgressResponse struct {
	ExerciseID   int64          `json:"exercise_id"`
	ExerciseName string         `json:"exercise_name"`
	Range        string         `json:"range"`
	StartDate    time.Time      `json:"start_date"`
	EndDate      time.Time      `json:"end_date"`
	DataPoints   []ProgressDataPoint `json:"data_points"`
	Summary      *ExerciseMetrics    `json:"summary,omitempty"`
}

// ProgressDataPoint represents a single data point in progress tracking
type ProgressDataPoint struct {
	Date         time.Time `json:"date"`
	TotalVolume  float64   `json:"total_volume"`
	MaxWeight    float64   `json:"max_weight"`
	TotalSets    int       `json:"total_sets"`
	AverageRPE   *float64  `json:"average_rpe,omitempty"`
}

