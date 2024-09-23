package resources

import "gcstatus/internal/domain"

type LevelResource struct {
	ID         uint   `json:"id"`
	Level      uint   `json:"level"`
	Coins      uint   `json:"coins"`
	Experience uint   `json:"experience"`
	CreatedAt  string `json:"created_at"`
}

func TransformLevel(level *domain.Level) *LevelResource {
	return &LevelResource{
		ID:         level.ID,
		Level:      level.Level,
		Coins:      level.Coins,
		Experience: level.Experience,
		CreatedAt:  level.CreatedAt.Format("2006-01-02T15:04:05"),
	}
}

func TransformLevels(levels []*domain.Level) []LevelResource {
	var resources []LevelResource

	for _, level := range levels {
		resources = append(resources, *TransformLevel(level))
	}

	return resources
}
