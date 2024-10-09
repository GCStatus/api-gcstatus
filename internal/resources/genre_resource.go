package resources

import "gcstatus/internal/domain"

type GenreResource struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func TransformGenre(genre domain.Genre) GenreResource {
	return GenreResource{
		ID:   genre.ID,
		Name: genre.Name,
	}
}

func TransformGenres(genres []domain.Genre) []GenreResource {
	var resources []GenreResource

	for _, genre := range genres {
		resources = append(resources, TransformGenre(genre))
	}

	return resources
}
