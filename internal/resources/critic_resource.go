package resources

import "gcstatus/internal/domain"

type CriticResource struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Logo string `json:"logo"`
}

func TransformCritic(critic domain.Critic) CriticResource {
	return CriticResource{
		ID:   critic.ID,
		Name: critic.Name,
		URL:  critic.URL,
		Logo: critic.Logo,
	}
}
