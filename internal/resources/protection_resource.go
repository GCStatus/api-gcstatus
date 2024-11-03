package resources

import "gcstatus/internal/domain"

type ProtectionResource struct {
	ID   uint   `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func TransformProtection(protection domain.Protection) *ProtectionResource {
	return &ProtectionResource{
		ID:   protection.ID,
		Name: protection.Name,
		Slug: protection.Slug,
	}
}

func TransformProtections(protections []domain.Protection) []*ProtectionResource {
	var resources []*ProtectionResource

	for _, protection := range protections {
		resources = append(resources, TransformProtection(protection))
	}

	return resources
}
