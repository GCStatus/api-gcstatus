package resources

import (
	"gcstatus/internal/domain"
)

type SupportResource struct {
	ID      uint    `json:"id"`
	URL     *string `json:"url"`
	Email   *string `json:"email"`
	Contact *string `json:"contact"`
}

func TransformSupport(support *domain.GameSupport) *SupportResource {
	return &SupportResource{
		ID:      support.ID,
		URL:     support.URL,
		Email:   support.Email,
		Contact: support.Contact,
	}
}
