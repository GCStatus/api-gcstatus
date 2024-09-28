package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type RewardResource struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	SourceableType string `json:"sourceable_type"`

	RewardableType string         `json:"rewardable_type"`
	Title          *TitleResource `json:"title,omitempty"`
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
	case "titles":
		if title, ok := reward.Rewardable.(*domain.Title); ok {
			rewardResource.Title = transformTitle(title)
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

func TransformRewards(rewards []domain.Reward) []RewardResource {
	var resources []RewardResource

	for _, reward := range rewards {
		resources = append(resources, *TransformReward(reward))
	}

	return resources
}
