package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/utils"
	"testing"
	"time"
)

func TestTransformReward(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    domain.Reward
		expected *resources.RewardResource
	}{
		"normal reward with title": {
			input: domain.Reward{
				ID:             1,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
				SourceableType: "source_type",
				RewardableType: "titles",
				Rewardable: &domain.Title{
					ID:          1,
					Title:       "Title 1",
					Description: "Description for Title 1",
					Purchasable: true,
					Cost:        nil,
					Status:      "available",
					CreatedAt:   fixedTime,
				},
			},
			expected: &resources.RewardResource{
				ID:             1,
				CreatedAt:      utils.FormatTimestamp(fixedTime),
				UpdatedAt:      utils.FormatTimestamp(fixedTime),
				SourceableType: "source_type",
				RewardableType: "titles",
				Title: &resources.TitleResource{
					ID:          1,
					Title:       "Title 1",
					Description: "Description for Title 1",
					Purchasable: true,
					Cost:        nil,
					Status:      "available",
					CreatedAt:   utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformReward(test.input)

			if result.ID != test.expected.ID ||
				result.CreatedAt != test.expected.CreatedAt ||
				result.UpdatedAt != test.expected.UpdatedAt ||
				result.SourceableType != test.expected.SourceableType ||
				result.RewardableType != test.expected.RewardableType {
				t.Errorf("Expected %+v, got %+v", test.expected, result)
			}

			if !titlesEqual(result.Title, test.expected.Title) {
				t.Errorf("Expected Title %+v, got %+v", test.expected.Title, result.Title)
			}
		})
	}
}

func TestTransformRewards(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    []domain.Reward
		expected []resources.RewardResource
	}{
		"multiple rewards": {
			input: []domain.Reward{
				{
					ID:             1,
					CreatedAt:      fixedTime,
					UpdatedAt:      fixedTime,
					SourceableType: "source_type",
					RewardableType: "titles",
					Rewardable: &domain.Title{
						ID:          1,
						Title:       "Title 1",
						Description: "Description for Title 1",
						Purchasable: true,
						Cost:        nil,
						Status:      "available",
						CreatedAt:   fixedTime,
					},
				},
				{
					ID:             2,
					CreatedAt:      fixedTime,
					UpdatedAt:      fixedTime,
					SourceableType: "source_type",
					RewardableType: "other_type",
					Rewardable:     nil,
				},
			},
			expected: []resources.RewardResource{
				{
					ID:             1,
					CreatedAt:      utils.FormatTimestamp(fixedTime),
					UpdatedAt:      utils.FormatTimestamp(fixedTime),
					SourceableType: "source_type",
					RewardableType: "titles",
					Title: &resources.TitleResource{
						ID:          1,
						Title:       "Title 1",
						Description: "Description for Title 1",
						Purchasable: true,
						Cost:        nil,
						Status:      "available",
						CreatedAt:   utils.FormatTimestamp(fixedTime),
					},
				},
				{
					ID:             2,
					CreatedAt:      utils.FormatTimestamp(fixedTime),
					UpdatedAt:      utils.FormatTimestamp(fixedTime),
					SourceableType: "source_type",
					RewardableType: "other_type",
					Title:          nil,
				},
			},
		},
		"empty rewards": {
			input:    []domain.Reward{},
			expected: []resources.RewardResource{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := resources.TransformRewards(test.input)

			if len(result) != len(test.expected) {
				t.Errorf("Expected %d rewards, got %d", len(test.expected), len(result))
				return
			}

			for i := range result {
				if result[i].ID != test.expected[i].ID ||
					result[i].CreatedAt != test.expected[i].CreatedAt ||
					result[i].UpdatedAt != test.expected[i].UpdatedAt ||
					result[i].SourceableType != test.expected[i].SourceableType ||
					result[i].RewardableType != test.expected[i].RewardableType {
					t.Errorf("Expected %+v, got %+v", test.expected[i], result[i])
				}

				if !titlesEqual(result[i].Title, test.expected[i].Title) {
					t.Errorf("Expected Title %+v, got %+v", test.expected[i].Title, result[i].Title)
				}
			}
		})
	}
}

func titlesEqual(a, b *resources.TitleResource) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if a.ID != b.ID ||
		a.Title != b.Title ||
		a.Description != b.Description ||
		a.Purchasable != b.Purchasable ||
		a.Cost != b.Cost ||
		a.Status != b.Status ||
		a.CreatedAt != b.CreatedAt {
		return false
	}

	return true
}
