package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformPlatform(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	tests := map[string]struct {
		input    domain.Platform
		expected resources_admin.PlatformResource
	}{
		"as null": {
			input: domain.Platform{},
			expected: resources_admin.PlatformResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"multiple categories": {
			input: domain.Platform{
				ID:        1,
				Name:      "Platform 1",
				Slug:      "platform-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources_admin.PlatformResource{
				ID:        1,
				Name:      "Platform 1",
				Slug:      "platform-1",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			platformResource := resources_admin.TransformPlatform(test.input)

			if !reflect.DeepEqual(platformResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, platformResource)
			}
		})
	}
}

func TestTransformPlatforms(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    []domain.Platform
		expected []resources_admin.PlatformResource
	}{
		"as null": {
			input:    []domain.Platform{},
			expected: []resources_admin.PlatformResource{},
		},
		"multiple categories": {
			input: []domain.Platform{
				{
					ID:        1,
					Name:      "Platform 1",
					Slug:      "platform-1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Name:      "Platform 2",
					Slug:      "platform-2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []resources_admin.PlatformResource{
				{
					ID:        1,
					Name:      "Platform 1",
					Slug:      "platform-1",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID:        2,
					Name:      "Platform 2",
					Slug:      "platform-2",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			platformsResources_admin := resources_admin.TransformPlatforms(test.input)

			if platformsResources_admin == nil {
				platformsResources_admin = []resources_admin.PlatformResource{}
			}

			if !reflect.DeepEqual(platformsResources_admin, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, platformsResources_admin)
			}
		})
	}
}
