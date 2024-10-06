package resources

import "gcstatus/internal/domain"

type RequirementTypeResource struct {
	ID        uint   `json:"id"`
	OS        string `json:"os"`
	Potential string `json:"potential"`
}

func TransformRequirementType(requirementType domain.RequirementType) RequirementTypeResource {
	return RequirementTypeResource{
		ID:        requirementType.ID,
		OS:        requirementType.OS,
		Potential: requirementType.Potential,
	}
}

func TransformRequirementTypes(requirementTypes []domain.RequirementType) []RequirementTypeResource {
	var resources []RequirementTypeResource

	for _, requirementType := range requirementTypes {
		resources = append(resources, TransformRequirementType(requirementType))
	}

	return resources
}
