package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type TitleResource struct {
	ID                uint                       `json:"id"`
	Title             string                     `json:"title"`
	Description       string                     `json:"description"`
	Cost              *int                       `json:"cost,omitempty"`
	Purchasable       bool                       `json:"purchasable"`
	Status            string                     `json:"status,omitempty"`
	CreatedAt         string                     `json:"created_at,omitempty"`
	TitleRequirements []TitleRequirementResource `json:"requirements"`
}

func TransformTitle(title domain.Title) TitleResource {
	titleResource := TitleResource{
		ID:          title.ID,
		Title:       title.Title,
		Description: title.Description,
		Cost:        title.Cost,
		Purchasable: title.Purchasable,
		Status:      title.Status,
		CreatedAt:   utils.FormatTimestamp(title.CreatedAt),
	}

	if len(title.TitleRequirements) > 0 {
		titleResource.TitleRequirements = TransformTitleRequirements(title.TitleRequirements)
	} else {
		titleResource.TitleRequirements = []TitleRequirementResource{}
	}

	return titleResource
}

func TransformTitles(titles []domain.Title) []TitleResource {
	var resources []TitleResource

	for _, title := range titles {
		resources = append(resources, TransformTitle(title))
	}

	return resources
}
