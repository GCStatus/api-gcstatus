package resources_admin

import "gcstatus/internal/domain"

type RoleableResource struct {
	ID   uint         `json:"id"`
	Role RoleResource `json:"role"`
}

func TransformRoleable(roleable domain.Roleable) RoleableResource {
	return RoleableResource{
		ID:   roleable.ID,
		Role: TransformRole(roleable.Role),
	}
}

func TransformRoleables(roleables []domain.Roleable) []RoleableResource {
	var resources []RoleableResource
	for _, roleable := range roleables {
		resources = append(resources, TransformRoleable(roleable))
	}
	return resources
}
