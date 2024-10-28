package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type CategoryResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformCategory(category domain.Category) CategoryResource {
	return CategoryResource{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: utils.FormatTimestamp(category.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(category.UpdatedAt),
	}
}

func TransformCategories(categories []domain.Category) []CategoryResource {
	var resources []CategoryResource

	for _, category := range categories {
		resources = append(resources, TransformCategory(category))
	}

	return resources
}
