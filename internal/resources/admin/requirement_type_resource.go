package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type RequirementTypeResource struct {
	ID        uint   `json:"id"`
	OS        string `json:"os"`
	Potential string `json:"potential"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformRequirementType(requirementType domain.RequirementType) RequirementTypeResource {
	return RequirementTypeResource{
		ID:        requirementType.ID,
		OS:        requirementType.OS,
		Potential: requirementType.Potential,
		CreatedAt: utils.FormatTimestamp(requirementType.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(requirementType.UpdatedAt),
	}
}

func TransformRequirementTypes(requirementTypes []domain.RequirementType) []RequirementTypeResource {
	var resources []RequirementTypeResource

	for _, requirementType := range requirementTypes {
		resources = append(resources, TransformRequirementType(requirementType))
	}

	return resources
}
