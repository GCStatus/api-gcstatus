package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type DeveloperResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Acting    bool   `json:"acting"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformDeveloper(developer domain.Developer) DeveloperResource {
	return DeveloperResource{
		ID:        developer.ID,
		Name:      developer.Name,
		Slug:      developer.Slug,
		Acting:    developer.Acting,
		CreatedAt: utils.FormatTimestamp(developer.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(developer.UpdatedAt),
	}
}

func TransformDevelopers(developers []domain.Developer) []DeveloperResource {
	var resources []DeveloperResource
	for _, developer := range developers {
		resources = append(resources, TransformDeveloper(developer))
	}
	return resources
}
