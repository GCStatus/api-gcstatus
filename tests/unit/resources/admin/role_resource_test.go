package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"reflect"
	"testing"
	"time"
)

func TestTransformRole(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Role
		expected resources_admin.RoleResource
	}{
		"single role": {
			input: domain.Role{
				ID:          1,
				Name:        "Technology",
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				Permissions: []domain.Permissionable{},
			},
			expected: resources_admin.RoleResource{
				ID:          1,
				Name:        "Technology",
				Permissions: []resources_admin.PermissionResource{},
			},
		},
		"role with permissions": {
			input: domain.Role{
				ID:        1,
				Name:      "Technology",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Permissions: []domain.Permissionable{
					{
						ID:                 1,
						PermissionableID:   1,
						PermissionableType: "users",
						Permission: domain.Permission{
							ID:    1,
							Scope: "test:permission",
						},
					},
				},
			},
			expected: resources_admin.RoleResource{
				ID:   1,
				Name: "Technology",
				Permissions: []resources_admin.PermissionResource{
					{
						ID:    1,
						Scope: "test:permission",
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformRole(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
