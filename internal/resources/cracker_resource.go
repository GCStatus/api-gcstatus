package resources

import "gcstatus/internal/domain"

type CrackerResource struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Acting bool   `json:"acting"`
}

func TransformCracker(cracker domain.Cracker) *CrackerResource {
	return &CrackerResource{
		ID:     cracker.ID,
		Name:   cracker.Name,
		Acting: cracker.Acting,
	}
}

func TransformCrackers(crackers []domain.Cracker) []*CrackerResource {
	var resources []*CrackerResource

	for _, cracker := range crackers {
		resources = append(resources, TransformCracker(cracker))
	}

	return resources
}
