package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type RewardResource struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	SourceableType string `json:"sourceable_type"`

	RewardableType string           `json:"rewardable_type"`
	Title          *TitleResource   `json:"title,omitempty"`
	Mission        *MissionResource `json:"mission,omitempty"`
}

func TransformReward(reward domain.Reward) *RewardResource {
	rewardResource := &RewardResource{
		ID:             reward.ID,
		CreatedAt:      utils.FormatTimestamp(reward.CreatedAt),
		UpdatedAt:      utils.FormatTimestamp(reward.UpdatedAt),
		SourceableType: reward.SourceableType,
		RewardableType: reward.RewardableType,
	}

	// Check for polymorphic Rewardable types
	switch reward.RewardableType {
	case domain.RewardableTypeTitles:
		if title, ok := reward.Rewardable.(*domain.Title); ok {
			rewardResource.Title = transformTitle(title)
		}
	case "missions":
		if mission, ok := reward.Rewardable.(*domain.Mission); ok {
			rewardResource.Mission = transformMission(mission)
		}
	}

	return rewardResource
}

func transformTitle(title *domain.Title) *TitleResource {
	return &TitleResource{
		ID:          title.ID,
		Title:       title.Title,
		Description: title.Description,
		Purchasable: title.Purchasable,
		Cost:        title.Cost,
		Status:      title.Status,
		CreatedAt:   utils.FormatTimestamp(title.CreatedAt),
	}
}

func transformMission(mission *domain.Mission) *MissionResource {
	return &MissionResource{
		ID:          mission.ID,
		Mission:     mission.Mission,
		Description: mission.Description,
		Coins:       mission.Coins,
		Experience:  mission.Experience,
		Status:      mission.Status,
		Frequency:   mission.Frequency,
		ResetTime:   utils.FormatTimestamp(mission.ResetTime),
		CreatedAt:   utils.FormatTimestamp(mission.CreatedAt),
	}
}

func TransformRewards(rewards []domain.Reward) []RewardResource {
	var resources []RewardResource

	for _, reward := range rewards {
		resources = append(resources, *TransformReward(reward))
	}

	return resources
}
