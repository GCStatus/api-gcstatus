package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type TitleProgressResource struct {
	ID        uint   `json:"id"`
	Progress  uint   `json:"progress"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformTitleProgress(titleProgress domain.TitleProgress) *TitleProgressResource {
	return &TitleProgressResource{
		ID:        titleProgress.ID,
		Progress:  uint(titleProgress.Progress),
		Completed: titleProgress.Completed,
		CreatedAt: utils.FormatTimestamp(titleProgress.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(titleProgress.UpdatedAt),
	}
}

func TransformTitleProgresses(titleProgresses []domain.TitleProgress) []TitleProgressResource {
	var resources []TitleProgressResource

	for _, titleProgress := range titleProgresses {
		resources = append(resources, *TransformTitleProgress(titleProgress))
	}

	return resources
}
