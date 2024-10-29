package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type CriticResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Logo      string `json:"logo"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformCritic(critic domain.Critic) CriticResource {
	return CriticResource{
		ID:        critic.ID,
		Name:      critic.Name,
		URL:       critic.URL,
		Logo:      critic.Logo,
		CreatedAt: utils.FormatTimestamp(critic.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(critic.UpdatedAt),
	}
}
