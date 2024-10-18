package resources_admin

import "gcstatus/internal/domain"

type PermissionResource struct {
	ID    uint   `json:"id"`
	Scope string `json:"scope"`
}

func TransformPermission(permission domain.Permission) PermissionResource {
	return PermissionResource{
		ID:    permission.ID,
		Scope: permission.Scope,
	}
}
