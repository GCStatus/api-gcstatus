package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"reflect"
	"testing"
	"time"
)

func TestTransformRoleable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Roleable
		expected resources_admin.RoleableResource
	}{
		"single roleable": {
			input: domain.Roleable{
				ID:           1,
				RoleableID:   1,
				RoleableType: "users",
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
				Role: domain.Role{
					ID:   1,
					Name: "Technology",
				},
			},
			expected: resources_admin.RoleableResource{
				ID: 1,
				Role: resources_admin.RoleResource{
					ID:          1,
					Name:        "Technology",
					Permissions: []resources_admin.PermissionResource{},
				},
			},
		},
		"roleable with permissions": {
			input: domain.Roleable{
				ID:           1,
				RoleableID:   1,
				RoleableType: "users",
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
				Role: domain.Role{
					ID:   1,
					Name: "Technology",
					Permissions: []domain.Permissionable{
						{
							ID:                 1,
							PermissionableID:   1,
							PermissionableType: "roles",
							Permission: domain.Permission{
								ID:    1,
								Scope: "test:permission",
							},
						},
					},
				},
			},
			expected: resources_admin.RoleableResource{
				ID: 1,
				Role: resources_admin.RoleResource{
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
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformRoleable(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestTransformRoleables(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    []domain.Roleable
		expected []resources_admin.RoleableResource
	}{
		"single roleable": {
			input: []domain.Roleable{
				{
					ID:           1,
					RoleableID:   1,
					RoleableType: "users",
					CreatedAt:    fixedTime,
					UpdatedAt:    fixedTime,
					Role: domain.Role{
						ID:   1,
						Name: "Technology",
					},
				},
			},
			expected: []resources_admin.RoleableResource{
				{
					ID: 1,
					Role: resources_admin.RoleResource{
						ID:          1,
						Name:        "Technology",
						Permissions: []resources_admin.PermissionResource{},
					},
				},
			},
		},
		"roleable with permissions": {
			input: []domain.Roleable{
				{
					ID:           1,
					RoleableID:   1,
					RoleableType: "users",
					CreatedAt:    fixedTime,
					UpdatedAt:    fixedTime,
					Role: domain.Role{
						ID:   1,
						Name: "Technology",
						Permissions: []domain.Permissionable{
							{
								ID:                 1,
								PermissionableID:   1,
								PermissionableType: "roles",
								Permission: domain.Permission{
									ID:    1,
									Scope: "test:permission",
								},
							},
						},
					},
				},
			},
			expected: []resources_admin.RoleableResource{
				{
					ID: 1,
					Role: resources_admin.RoleResource{
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
			},
		},
		"multiple roleables": {
			input: []domain.Roleable{
				{
					ID:           1,
					RoleableID:   1,
					RoleableType: "users",
					CreatedAt:    fixedTime,
					UpdatedAt:    fixedTime,
					Role: domain.Role{
						ID:   1,
						Name: "Technology",
					},
				},
				{
					ID:           2,
					RoleableID:   2,
					RoleableType: "users",
					CreatedAt:    fixedTime,
					UpdatedAt:    fixedTime,
					Role: domain.Role{
						ID:   1,
						Name: "Technology",
					},
				},
			},
			expected: []resources_admin.RoleableResource{
				{
					ID: 1,
					Role: resources_admin.RoleResource{
						ID:          1,
						Name:        "Technology",
						Permissions: []resources_admin.PermissionResource{},
					},
				},
				{
					ID: 2,
					Role: resources_admin.RoleResource{
						ID:          1,
						Name:        "Technology",
						Permissions: []resources_admin.PermissionResource{},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformRoleables(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
