package resources

import "gcstatus/internal/domain"

type MediaTypeResource struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func TransformMediaType(mediaType domain.MediaType) MediaTypeResource {
	return MediaTypeResource{
		ID:   mediaType.ID,
		Name: mediaType.Name,
	}
}
