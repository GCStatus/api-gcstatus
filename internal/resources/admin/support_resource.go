package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type SupportResource struct {
	ID        uint    `json:"id"`
	URL       *string `json:"url"`
	Email     *string `json:"email"`
	Contact   *string `json:"contact"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func TransformSupport(support *domain.GameSupport) *SupportResource {
	return &SupportResource{
		ID:        support.ID,
		URL:       support.URL,
		Email:     support.Email,
		Contact:   support.Contact,
		CreatedAt: utils.FormatTimestamp(support.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(support.UpdatedAt),
	}
}
