package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type CrackerResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Acting    bool   `json:"acting"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformCracker(cracker domain.Cracker) *CrackerResource {
	return &CrackerResource{
		ID:        cracker.ID,
		Name:      cracker.Name,
		Slug:      cracker.Slug,
		Acting:    cracker.Acting,
		CreatedAt: utils.FormatTimestamp(cracker.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(cracker.UpdatedAt),
	}
}

func TransformCrackers(crackers []domain.Cracker) []*CrackerResource {
	var resources []*CrackerResource
	for _, cracker := range crackers {
		resources = append(resources, TransformCracker(cracker))
	}
	return resources
}
