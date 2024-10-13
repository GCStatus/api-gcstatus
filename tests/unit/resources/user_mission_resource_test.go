package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/utils"
	"testing"
	"time"
)

func TestTransformUserMission(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    domain.UserMission
		expected *resources.UserMissionResource
	}{
		"normal progress": {
			input: domain.UserMission{
				ID:              1,
				Completed:       false,
				LastCompletedAt: fixedTime,
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
			},
			expected: &resources.UserMissionResource{
				Completed:       false,
				LastCompletedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformUserMission(test.input)

			if *result != *test.expected {
				t.Errorf("Expected %+v, got %+v", test.expected, result)
			}
		})
	}
}

func TestTransformUserMissiones(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    []domain.UserMission
		expected []resources.UserMissionResource
	}{
		"multiple progresses": {
			input: []domain.UserMission{
				{
					ID:              1,
					Completed:       false,
					LastCompletedAt: fixedTime,
					CreatedAt:       fixedTime,
					UpdatedAt:       fixedTime,
				},
				{
					ID:              2,
					Completed:       true,
					LastCompletedAt: fixedTime,
					CreatedAt:       fixedTime,
					UpdatedAt:       fixedTime,
				},
			},
			expected: []resources.UserMissionResource{
				{
					Completed:       false,
					LastCompletedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					Completed:       true,
					LastCompletedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
		"empty progresses": {
			input:    []domain.UserMission{},
			expected: []resources.UserMissionResource{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformUserMissions(test.input)

			if len(result) != len(test.expected) {
				t.Errorf("Expected %d progresses, got %d", len(test.expected), len(result))
				return
			}

			for i := range result {
				if *result[i] != test.expected[i] {
					t.Errorf("Expected %+v, got %+v", test.expected[i], result[i])
				}
			}
		})
	}
}
