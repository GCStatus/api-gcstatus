package resources

import "gcstatus/internal/domain"

type CategoryResource struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func TransformCategory(category domain.Category) CategoryResource {
	return CategoryResource{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}
}

func TransformCategories(categories []domain.Category) []CategoryResource {
	var resources []CategoryResource

	for _, category := range categories {
		resources = append(resources, TransformCategory(category))
	}

	return resources
}
