package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"reflect"
	"testing"
	"time"
)

func TestTransformPermissionable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Permissionable
		expected resources_admin.PermissionableResource
	}{
		"single permissionable": {
			input: domain.Permissionable{
				ID:                 1,
				PermissionableID:   1,
				PermissionableType: "users",
				CreatedAt:          fixedTime,
				UpdatedAt:          fixedTime,
				Permission: domain.Permission{
					ID:    1,
					Scope: "test:permission",
				},
			},
			expected: resources_admin.PermissionableResource{
				ID: 1,
				Permission: resources_admin.PermissionResource{
					ID:    1,
					Scope: "test:permission",
				},
			},
		},
	}

	for Scope, tc := range testCases {
		t.Run(Scope, func(t *testing.T) {
			result := resources_admin.TransformPermissionable(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestTransformPermissionables(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    []domain.Permissionable
		expected []resources_admin.PermissionableResource
	}{
		"single permissionable": {
			input: []domain.Permissionable{
				{
					ID:                 1,
					PermissionableID:   1,
					PermissionableType: "users",
					CreatedAt:          fixedTime,
					UpdatedAt:          fixedTime,
					Permission: domain.Permission{
						ID:    1,
						Scope: "test:permission",
					},
				},
			},
			expected: []resources_admin.PermissionableResource{
				{
					ID: 1,
					Permission: resources_admin.PermissionResource{
						ID:    1,
						Scope: "test:permission",
					},
				},
			},
		},
		"multiple permissionable": {
			input: []domain.Permissionable{
				{
					ID:                 1,
					PermissionableID:   1,
					PermissionableType: "users",
					CreatedAt:          fixedTime,
					UpdatedAt:          fixedTime,
					Permission: domain.Permission{
						ID:    1,
						Scope: "test:permission",
					},
				},
				{
					ID:                 2,
					PermissionableID:   2,
					PermissionableType: "users",
					CreatedAt:          fixedTime,
					UpdatedAt:          fixedTime,
					Permission: domain.Permission{
						ID:    1,
						Scope: "test:permission",
					},
				},
			},
			expected: []resources_admin.PermissionableResource{
				{
					ID: 1,
					Permission: resources_admin.PermissionResource{
						ID:    1,
						Scope: "test:permission",
					},
				},
				{
					ID: 2,
					Permission: resources_admin.PermissionResource{
						ID:    1,
						Scope: "test:permission",
					},
				},
			},
		},
	}

	for Scope, tc := range testCases {
		t.Run(Scope, func(t *testing.T) {
			result := resources_admin.TransformPermissionables(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
