package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type LevelResource struct {
	ID         uint             `json:"id"`
	Level      uint             `json:"level"`
	Coins      uint             `json:"coins"`
	Experience uint             `json:"experience"`
	CreatedAt  string           `json:"created_at"`
	Rewards    []RewardResource `json:"rewards"`
}

func TransformLevel(level *domain.Level) *LevelResource {
	resource := LevelResource{
		ID:         level.ID,
		Level:      level.Level,
		Coins:      level.Coins,
		Experience: level.Experience,
		CreatedAt:  utils.FormatTimestamp(level.CreatedAt),
	}

	if len(level.Rewards) > 0 {
		resource.Rewards = TransformRewards(level.Rewards)
	} else {
		resource.Rewards = []RewardResource{}
	}

	return &resource
}

func TransformLevels(levels []*domain.Level) []LevelResource {
	resources := make([]LevelResource, 0, len(levels))

	for _, level := range levels {
		resources = append(resources, *TransformLevel(level))
	}

	return resources
}
