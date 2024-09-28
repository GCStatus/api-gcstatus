package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type TitleRequirementResource struct {
	ID            uint                   `json:"id"`
	Task          string                 `json:"task"`
	Description   string                 `json:"description"`
	Goal          int                    `json:"goal"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
	TitleProgress *TitleProgressResource `json:"progress,omitempty"`
}

func TransformTitleRequirement(titleRequirement domain.TitleRequirement) TitleRequirementResource {
	titleRequirementResource := TitleRequirementResource{
		ID:          titleRequirement.ID,
		Task:        titleRequirement.Task,
		Description: titleRequirement.Description,
		Goal:        titleRequirement.Goal,
		CreatedAt:   utils.FormatTimestamp(titleRequirement.CreatedAt),
		UpdatedAt:   utils.FormatTimestamp(titleRequirement.UpdatedAt),
	}

	if titleRequirement.TitleProgress.ID != 0 {
		titleRequirementResource.TitleProgress = TransformTitleProgress(titleRequirement.TitleProgress)
	} else {
		titleRequirementResource.TitleProgress = nil
	}

	return titleRequirementResource
}

func TransformTitleRequirements(titleRequirements []domain.TitleRequirement) []TitleRequirementResource {
	if len(titleRequirements) == 0 {
		return []TitleRequirementResource{}
	}

	var resources []TitleRequirementResource

	for _, titleRequirement := range titleRequirements {
		resources = append(resources, TransformTitleRequirement(titleRequirement))
	}

	return resources
}
