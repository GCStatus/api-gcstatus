package resources

import "gcstatus/internal/domain"

type TagResource struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func TransformTag(tag domain.Tag) TagResource {
	return TagResource{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}
}

func TransformTags(tags []domain.Tag) []TagResource {
	var resources []TagResource

	for _, tag := range tags {
		resources = append(resources, TransformTag(tag))
	}

	return resources
}
