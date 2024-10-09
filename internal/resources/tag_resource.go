package resources

import "gcstatus/internal/domain"

type TagResource struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func TransformTag(tag domain.Tag) TagResource {
	return TagResource{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func TransformTags(tags []domain.Tag) []TagResource {
	var resources []TagResource

	for _, tag := range tags {
		resources = append(resources, TransformTag(tag))
	}

	return resources
}
