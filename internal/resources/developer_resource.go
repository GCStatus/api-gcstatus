package resources

import "gcstatus/internal/domain"

type DeveloperResource struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Acting bool   `json:"acting"`
}

func TransformDeveloper(Developer domain.Developer) DeveloperResource {
	return DeveloperResource{
		ID:     Developer.ID,
		Name:   Developer.Name,
		Acting: Developer.Acting,
	}
}

func TransformDevelopers(Developers []domain.Developer) []DeveloperResource {
	var resources []DeveloperResource

	for _, Developer := range Developers {
		resources = append(resources, TransformDeveloper(Developer))
	}

	return resources
}
