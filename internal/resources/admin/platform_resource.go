package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type PlatformResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformPlatform(platform domain.Platform) PlatformResource {
	return PlatformResource{
		ID:        platform.ID,
		Name:      platform.Name,
		Slug:      platform.Slug,
		CreatedAt: utils.FormatTimestamp(platform.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(platform.UpdatedAt),
	}
}

func TransformPlatforms(platforms []domain.Platform) []PlatformResource {
	var resources []PlatformResource

	for _, platform := range platforms {
		resources = append(resources, TransformPlatform(platform))
	}

	return resources
}
