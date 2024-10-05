package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/pkg/utils"
	"testing"
	"time"
)

func TestTransformMissionRequirement(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    domain.MissionRequirement
		expected resources.MissionRequirementResource
	}{
		"normal requirement": {
			input: domain.MissionRequirement{
				ID:          1,
				Task:        "Complete Task",
				Description: "Complete the task description.",
				Goal:        10,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				MissionProgress: domain.MissionProgress{
					ID:        1,
					Progress:  5,
					Completed: false,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: resources.MissionRequirementResource{
				ID:          1,
				Task:        "Complete Task",
				Description: "Complete the task description.",
				Goal:        10,
				CreatedAt:   utils.FormatTimestamp(fixedTime),
				UpdatedAt:   utils.FormatTimestamp(fixedTime),
				MissionProgress: &resources.MissionProgressResource{
					ID:        1,
					Progress:  5,
					Completed: false,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
		"requirement without progress": {
			input: domain.MissionRequirement{
				ID:          2,
				Task:        "Another Task",
				Description: "Another task description.",
				Goal:        5,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				MissionProgress: domain.MissionProgress{
					ID: 0,
				},
			},
			expected: resources.MissionRequirementResource{
				ID:              2,
				Task:            "Another Task",
				Description:     "Another task description.",
				Goal:            5,
				CreatedAt:       utils.FormatTimestamp(fixedTime),
				UpdatedAt:       utils.FormatTimestamp(fixedTime),
				MissionProgress: nil,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformMissionRequirement(test.input)

			if test.expected.MissionProgress != nil {
				if result.MissionProgress == nil || *result.MissionProgress != *test.expected.MissionProgress {
					t.Errorf("Expected MissionProgress %+v, got %+v", test.expected.MissionProgress, result.MissionProgress)
				}
			} else if result.MissionProgress != nil {
				t.Errorf("Expected MissionProgress nil, got %+v", result.MissionProgress)
			}

			if result.ID != test.expected.ID ||
				result.Task != test.expected.Task ||
				result.Description != test.expected.Description ||
				result.Goal != test.expected.Goal ||
				result.CreatedAt != test.expected.CreatedAt ||
				result.UpdatedAt != test.expected.UpdatedAt {
				t.Errorf("Expected %+v, got %+v", test.expected, result)
			}
		})
	}
}

func TestTransformMissionRequirements(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    []domain.MissionRequirement
		expected []resources.MissionRequirementResource
	}{
		"multiple requirements": {
			input: []domain.MissionRequirement{
				{
					ID:          1,
					Task:        "Task 1",
					Description: "Description 1",
					Goal:        10,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				},
				{
					ID:          2,
					Task:        "Task 2",
					Description: "Description 2",
					Goal:        5,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				},
			},
			expected: []resources.MissionRequirementResource{
				{
					ID:          1,
					Task:        "Task 1",
					Description: "Description 1",
					Goal:        10,
					CreatedAt:   utils.FormatTimestamp(fixedTime),
					UpdatedAt:   utils.FormatTimestamp(fixedTime),
				},
				{
					ID:          2,
					Task:        "Task 2",
					Description: "Description 2",
					Goal:        5,
					CreatedAt:   utils.FormatTimestamp(fixedTime),
					UpdatedAt:   utils.FormatTimestamp(fixedTime),
				},
			},
		},
		"empty requirements": {
			input:    []domain.MissionRequirement{},
			expected: []resources.MissionRequirementResource{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformMissionRequirements(test.input)

			if len(result) != len(test.expected) {
				t.Errorf("Expected %d requirements, got %d", len(test.expected), len(result))
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
