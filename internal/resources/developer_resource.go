package resources

import "gcstatus/internal/domain"

type DeveloperResource struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Acting bool   `json:"acting"`
}

func TransformDeveloper(developer domain.Developer) DeveloperResource {
	return DeveloperResource{
		ID:     developer.ID,
		Name:   developer.Name,
		Acting: developer.Acting,
	}
}

func TransformDevelopers(developers []domain.Developer) []DeveloperResource {
	var resources []DeveloperResource
	for _, developer := range developers {
		resources = append(resources, TransformDeveloper(developer))
	}
	return resources
}
