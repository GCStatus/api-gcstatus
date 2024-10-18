package resources_admin

import "gcstatus/internal/domain"

type RoleResource struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Permissions []PermissionResource `json:"permissions"`
}

func TransformRole(role domain.Role) RoleResource {
	resource := RoleResource{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: []PermissionResource{},
	}

	if len(role.Permissions) > 0 {
		for _, permissionable := range role.Permissions {
			resource.Permissions = append(resource.Permissions, TransformPermission(permissionable.Permission))
		}
	}

	return resource
}
