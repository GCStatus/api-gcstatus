package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	testutils "gcstatus/tests/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformUser(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.User
		expected resources_admin.UserResource
	}{
		"with roles and permissions": {
			input: domain.User{
				ID:         1,
				Name:       "John Doe",
				Email:      "john@example.com",
				Nickname:   "Johnny",
				Blocked:    false,
				Birthdate:  fixedTime,
				Experience: 500,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
				Profile:    domain.Profile{ID: 1, Photo: "key-1", CreatedAt: fixedTime, UpdatedAt: fixedTime},
				Roles: []domain.Roleable{
					{
						ID:           1,
						RoleableID:   1,
						RoleableType: "users",
						Role: domain.Role{
							ID:   1,
							Name: "Technology",
							Permissions: []domain.Permissionable{
								{
									ID:                 2,
									PermissionableID:   1,
									PermissionableType: "roles",
									Permission: domain.Permission{
										ID:    2,
										Scope: "test:permission-2",
									},
								},
							},
						},
					},
				},
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
			expected: resources_admin.UserResource{
				ID:        1,
				Name:      "John Doe",
				Email:     "john@example.com",
				Nickname:  "Johnny",
				Birthdate: utils.FormatTimestamp(fixedTime),
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				Profile: &resources.ProfileResource{
					ID:        1,
					Share:     false,
					Photo:     "https://mock-presigned-url.com/key-1",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				Roles: []resources_admin.RoleResource{
					{
						ID:   1,
						Name: "Technology",
						Permissions: []resources_admin.PermissionResource{
							{
								ID:    2,
								Scope: "test:permission-2",
							},
						},
					},
				},
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
			mockS3Client := &testutils.MockS3Client{}
			result := resources_admin.TransformUser(tc.input, mockS3Client)

			if !reflect.DeepEqual(result.ID, tc.expected.ID) ||
				!reflect.DeepEqual(result.Name, tc.expected.Name) ||
				!reflect.DeepEqual(result.Email, tc.expected.Email) ||
				!reflect.DeepEqual(result.Nickname, tc.expected.Nickname) ||
				!reflect.DeepEqual(result.Birthdate, tc.expected.Birthdate) ||
				!reflect.DeepEqual(result.CreatedAt, tc.expected.CreatedAt) ||
				!reflect.DeepEqual(result.UpdatedAt, tc.expected.UpdatedAt) ||
				!reflect.DeepEqual(result.Roles, tc.expected.Roles) ||
				!reflect.DeepEqual(result.Permissions, tc.expected.Permissions) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}

			if !reflect.DeepEqual(*result.Profile, *tc.expected.Profile) {
				t.Errorf("Expected profile %+v, got %+v", *tc.expected.Profile, *result.Profile)
			}
		})
	}
}
