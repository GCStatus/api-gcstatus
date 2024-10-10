package resources

import "gcstatus/internal/domain"

type PlatformResource struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func TransformPlatform(platform domain.Platform) PlatformResource {
	return PlatformResource{
		ID:   platform.ID,
		Name: platform.Name,
		Slug: platform.Slug,
	}
}

func TransformPlatforms(platforms []domain.Platform) []PlatformResource {
	var resources []PlatformResource

	for _, platform := range platforms {
		resources = append(resources, TransformPlatform(platform))
	}

	return resources
}
