package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type ProtectionResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformProtection(protection domain.Protection) *ProtectionResource {
	return &ProtectionResource{
		ID:        protection.ID,
		Name:      protection.Name,
		CreatedAt: utils.FormatTimestamp(protection.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(protection.UpdatedAt),
	}
}

func TransformProtections(protections []domain.Protection) []*ProtectionResource {
	var resources []*ProtectionResource

	for _, protection := range protections {
		resources = append(resources, TransformProtection(protection))
	}

	return resources
}
