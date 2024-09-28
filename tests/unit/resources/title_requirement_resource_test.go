package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/pkg/utils"
	"testing"
	"time"
)

func TestTransformTitleRequirement(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    domain.TitleRequirement
		expected resources.TitleRequirementResource
	}{
		"normal requirement": {
			input: domain.TitleRequirement{
				ID:          1,
				Task:        "Complete Task",
				Description: "Complete the task description.",
				Goal:        10,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				TitleProgress: domain.TitleProgress{
					ID:        1,
					Progress:  5,
					Completed: false,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: resources.TitleRequirementResource{
				ID:          1,
				Task:        "Complete Task",
				Description: "Complete the task description.",
				Goal:        10,
				CreatedAt:   utils.FormatTimestamp(fixedTime),
				UpdatedAt:   utils.FormatTimestamp(fixedTime),
				TitleProgress: &resources.TitleProgressResource{
					ID:        1,
					Progress:  5,
					Completed: false,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
		"requirement without progress": {
			input: domain.TitleRequirement{
				ID:          2,
				Task:        "Another Task",
				Description: "Another task description.",
				Goal:        5,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				TitleProgress: domain.TitleProgress{
					ID: 0,
				},
			},
			expected: resources.TitleRequirementResource{
				ID:            2,
				Task:          "Another Task",
				Description:   "Another task description.",
				Goal:          5,
				CreatedAt:     utils.FormatTimestamp(fixedTime),
				UpdatedAt:     utils.FormatTimestamp(fixedTime),
				TitleProgress: nil,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformTitleRequirement(test.input)

			if test.expected.TitleProgress != nil {
				if result.TitleProgress == nil || *result.TitleProgress != *test.expected.TitleProgress {
					t.Errorf("Expected TitleProgress %+v, got %+v", test.expected.TitleProgress, result.TitleProgress)
				}
			} else if result.TitleProgress != nil {
				t.Errorf("Expected TitleProgress nil, got %+v", result.TitleProgress)
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

func TestTransformTitleRequirements(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    []domain.TitleRequirement
		expected []resources.TitleRequirementResource
	}{
		"multiple requirements": {
			input: []domain.TitleRequirement{
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
			expected: []resources.TitleRequirementResource{
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
			input:    []domain.TitleRequirement{},
			expected: []resources.TitleRequirementResource{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformTitleRequirements(test.input)

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
