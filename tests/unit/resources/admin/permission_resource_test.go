package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"reflect"
	"testing"
	"time"
)

func TestTransformPermission(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Permission
		expected resources_admin.PermissionResource
	}{
		"single permission": {
			input: domain.Permission{
				ID:        1,
				Scope:     "view:test",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources_admin.PermissionResource{
				ID:    1,
				Scope: "view:test",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformPermission(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
