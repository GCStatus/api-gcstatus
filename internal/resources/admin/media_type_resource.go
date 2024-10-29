package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type MediaTypeResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformMediaType(mediaType domain.MediaType) MediaTypeResource {
	return MediaTypeResource{
		ID:        mediaType.ID,
		Name:      mediaType.Name,
		CreatedAt: utils.FormatTimestamp(mediaType.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(mediaType.UpdatedAt),
	}
}
