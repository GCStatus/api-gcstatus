package resources_admin

import "gcstatus/internal/domain"

type PermissionableResource struct {
	ID         uint               `json:"id"`
	Permission PermissionResource `json:"permission"`
}

func TransformPermissionable(permissionable domain.Permissionable) PermissionableResource {
	return PermissionableResource{
		ID:         permissionable.ID,
		Permission: TransformPermission(permissionable.Permission),
	}
}

func TransformPermissionables(permissionables []domain.Permissionable) []PermissionableResource {
	var resources []PermissionableResource
	for _, permissionable := range permissionables {
		resources = append(resources, TransformPermissionable(permissionable))
	}
	return resources
}
