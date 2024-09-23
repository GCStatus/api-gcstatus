package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"testing"
	"time"
)

func TestTransformLevel(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		inputLevel domain.Level
		expected   resources.LevelResource
	}{
		"normal level": {
			inputLevel: domain.Level{
				ID:         1,
				Level:      1,
				Experience: 500,
				Coins:      100,
				CreatedAt:  fixedTime,
			},
			expected: resources.LevelResource{
				ID:         1,
				Level:      1,
				Experience: 500,
				Coins:      100,
				CreatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
			},
		},
		"missing level": {
			inputLevel: domain.Level{
				ID:         1,
				Experience: 500,
				Coins:      100,
				CreatedAt:  fixedTime,
			},
			expected: resources.LevelResource{
				ID:         1,
				Experience: 500,
				Coins:      100,
				CreatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
			},
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			levelResource := resources.TransformLevel(&test.inputLevel)

			if levelResource.ID != test.expected.ID {
				t.Errorf("Expected ID %d, got %d", test.expected.ID, levelResource.ID)
			}
			if levelResource.Level != test.expected.Level {
				t.Errorf("Expected Level %d, got %d", test.expected.Level, levelResource.Level)
			}
			if levelResource.Coins != test.expected.Coins {
				t.Errorf("Expected Coins %d, got %d", test.expected.Coins, levelResource.Coins)
			}
			if levelResource.Experience != test.expected.Experience {
				t.Errorf("Expected Experience %d, got %d", test.expected.Experience, levelResource.Experience)
			}
			if levelResource.CreatedAt != test.expected.CreatedAt {
				t.Errorf("Expected CreatedAt %s, got %s", test.expected.CreatedAt, levelResource.CreatedAt)
			}
		})
	}
}

func TestTransformLevels(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		inputLevels []*domain.Level
		expected    []resources.LevelResource
	}{
		"two users": {
			inputLevels: []*domain.Level{
				{
					ID:         1,
					Level:      1,
					Experience: 500,
					Coins:      100,
					CreatedAt:  fixedTime,
				},
				{
					ID:         2,
					Level:      2,
					Experience: 1000,
					Coins:      150,
					CreatedAt:  fixedTime,
				},
			},
			expected: []resources.LevelResource{
				{
					ID:         1,
					Level:      1,
					Experience: 500,
					Coins:      100,
					CreatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				},
				{
					ID:         2,
					Level:      2,
					Experience: 1000,
					Coins:      150,
					CreatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			levelResources := resources.TransformLevels(test.inputLevels)

			if len(levelResources) != len(test.expected) {
				t.Errorf("Expected %d resources, got %d", len(test.expected), len(levelResources))
			}

			for i, level := range test.inputLevels {
				if levelResources[i].ID != level.ID {
					t.Errorf("Expected ID %d, got %d", level.ID, levelResources[i].ID)
				}
				if levelResources[i].Level != level.Level {
					t.Errorf("Expected Level %d, got %d", level.Level, levelResources[i].Level)
				}
				if levelResources[i].Coins != level.Coins {
					t.Errorf("Expected Coins %d, got %d", level.Coins, levelResources[i].Coins)
				}
				if levelResources[i].Experience != level.Experience {
					t.Errorf("Expected Experience %d, got %d", level.Experience, levelResources[i].Experience)
				}
				if levelResources[i].CreatedAt != level.CreatedAt.Format("2006-01-02T15:04:05") {
					t.Errorf("Expected CreatedAt %s, got %s", level.CreatedAt, levelResources[i].CreatedAt)
				}
			}
		})
	}
}
