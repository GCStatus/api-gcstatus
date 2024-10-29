package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type TagResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformTag(tag domain.Tag) TagResource {
	return TagResource{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: utils.FormatTimestamp(tag.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(tag.UpdatedAt),
	}
}

func TransformTags(tags []domain.Tag) []TagResource {
	var resources []TagResource

	for _, tag := range tags {
		resources = append(resources, TransformTag(tag))
	}

	return resources
}
