package service

import (
	"testing"
	"time"

	"phoenix-alliance-be/internal/models"
)

func TestCalculateMetrics(t *testing.T) {
	now := time.Now()
	sets := []*models.Set{
		{
			ID:         1,
			Weight:     60.0,
			Reps:       10,
			RestSeconds: intPtr(120),
			RPE:         intPtr(7),
			CreatedAt:   now,
		},
		{
			ID:         2,
			Weight:     65.0,
			Reps:       8,
			RestSeconds: intPtr(120),
			RPE:         intPtr(8),
			CreatedAt:   now.Add(time.Minute),
		},
		{
			ID:         3,
			Weight:     70.0,
			Reps:       6,
			RestSeconds: nil,
			RPE:         intPtr(9),
			CreatedAt:   now.Add(2 * time.Minute),
		},
	}

	metrics := calculateMetrics(sets)

	if metrics == nil {
		t.Fatal("Expected metrics, got nil")
	}

	if metrics.TotalSets != 3 {
		t.Errorf("Expected TotalSets to be 3, got %d", metrics.TotalSets)
	}

	expectedVolume := 60.0*10 + 65.0*8 + 70.0*6
	if metrics.TotalVolume != expectedVolume {
		t.Errorf("Expected TotalVolume to be %.2f, got %.2f", expectedVolume, metrics.TotalVolume)
	}

	if metrics.MaxWeight != 70.0 {
		t.Errorf("Expected MaxWeight to be 70.0, got %.2f", metrics.MaxWeight)
	}

	if metrics.MaxReps != 10 {
		t.Errorf("Expected MaxReps to be 10, got %d", metrics.MaxReps)
	}

	expectedAvgWeight := (60.0 + 65.0 + 70.0) / 3.0
	if metrics.AverageWeight != expectedAvgWeight {
		t.Errorf("Expected AverageWeight to be %.2f, got %.2f", expectedAvgWeight, metrics.AverageWeight)
	}

	expectedAvgReps := float64(10+8+6) / 3.0
	if metrics.AverageReps != expectedAvgReps {
		t.Errorf("Expected AverageReps to be %.2f, got %.2f", expectedAvgReps, metrics.AverageReps)
	}

	if metrics.AverageRest == nil {
		t.Error("Expected AverageRest to be set")
	} else {
		expectedAvgRest := 120.0
		if *metrics.AverageRest != expectedAvgRest {
			t.Errorf("Expected AverageRest to be %.2f, got %.2f", expectedAvgRest, *metrics.AverageRest)
		}
	}

	if metrics.AverageRPE == nil {
		t.Error("Expected AverageRPE to be set")
	} else {
		expectedAvgRPE := (7.0 + 8.0 + 9.0) / 3.0
		if *metrics.AverageRPE != expectedAvgRPE {
			t.Errorf("Expected AverageRPE to be %.2f, got %.2f", expectedAvgRPE, *metrics.AverageRPE)
		}
	}
}

func intPtr(i int) *int {
	return &i
}

