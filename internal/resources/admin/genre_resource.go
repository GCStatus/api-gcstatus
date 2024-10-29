package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type GenreResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformGenre(genre domain.Genre) GenreResource {
	return GenreResource{
		ID:        genre.ID,
		Name:      genre.Name,
		Slug:      genre.Slug,
		CreatedAt: utils.FormatTimestamp(genre.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(genre.UpdatedAt),
	}
}

func TransformGenres(genres []domain.Genre) []GenreResource {
	var resources []GenreResource

	for _, genre := range genres {
		resources = append(resources, TransformGenre(genre))
	}

	return resources
}
