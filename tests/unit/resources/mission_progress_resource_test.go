package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/pkg/utils"
	"testing"
	"time"
)

func TestTransformMissionProgress(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    domain.MissionProgress
		expected *resources.MissionProgressResource
	}{
		"normal progress": {
			input: domain.MissionProgress{
				ID:        1,
				Progress:  5,
				Completed: false,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: &resources.MissionProgressResource{
				ID:        1,
				Progress:  5,
				Completed: false,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformMissionProgress(test.input)

			if *result != *test.expected {
				t.Errorf("Expected %+v, got %+v", test.expected, result)
			}
		})
	}
}

func TestTransformMissionProgresses(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    []domain.MissionProgress
		expected []resources.MissionProgressResource
	}{
		"multiple progresses": {
			input: []domain.MissionProgress{
				{
					ID:        1,
					Progress:  5,
					Completed: false,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Progress:  10,
					Completed: true,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []resources.MissionProgressResource{
				{
					ID:        1,
					Progress:  5,
					Completed: false,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID:        2,
					Progress:  10,
					Completed: true,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
		"empty progresses": {
			input:    []domain.MissionProgress{},
			expected: []resources.MissionProgressResource{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformMissionProgresses(test.input)

			if len(result) != len(test.expected) {
				t.Errorf("Expected %d progresses, got %d", len(test.expected), len(result))
				return
			}

			for i := range result {
				if result[i] != test.expected[i] {
					t.Errorf("Expected %+v, got %+v", test.expected[i], result[i])
				}
			}
		})
	}
}
