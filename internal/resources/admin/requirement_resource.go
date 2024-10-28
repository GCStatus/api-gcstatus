package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type RequirementResource struct {
	ID              uint                    `json:"id"`
	OS              string                  `json:"os"`
	DX              string                  `json:"dx"`
	CPU             string                  `json:"cpu"`
	RAM             string                  `json:"ram"`
	GPU             string                  `json:"gpu"`
	ROM             string                  `json:"rom"`
	OBS             *string                 `json:"obs,omitempty"`
	Network         string                  `json:"network"`
	RequirementType RequirementTypeResource `json:"requirement_type"`
	CreatedAt       string                  `json:"created_at"`
	UpdatedAt       string                  `json:"updated_at"`
}

func TransformRequirement(requirement domain.Requirement) RequirementResource {
	return RequirementResource{
		ID:              requirement.ID,
		OS:              requirement.OS,
		DX:              requirement.DX,
		CPU:             requirement.CPU,
		RAM:             requirement.RAM,
		GPU:             requirement.GPU,
		ROM:             requirement.ROM,
		OBS:             requirement.OBS,
		Network:         requirement.Network,
		CreatedAt:       utils.FormatTimestamp(requirement.CreatedAt),
		UpdatedAt:       utils.FormatTimestamp(requirement.UpdatedAt),
		RequirementType: TransformRequirementType(requirement.RequirementType),
	}
}
