package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
	"time"
)

func TestTransformPlatform(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    domain.Platform
		expected resources.PlatformResource
	}{
		"as null": {
			input:    domain.Platform{},
			expected: resources.PlatformResource{},
		},
		"multiple categories": {
			input: domain.Platform{
				ID:        1,
				Name:      "Platform 1",
				Slug:      "platform-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources.PlatformResource{
				ID:   1,
				Name: "Platform 1",
				Slug: "platform-1",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			platformResource := resources.TransformPlatform(test.input)

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
		expected []resources.PlatformResource
	}{
		"as null": {
			input:    []domain.Platform{},
			expected: []resources.PlatformResource{},
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
			expected: []resources.PlatformResource{
				{
					ID:   1,
					Name: "Platform 1",
					Slug: "platform-1",
				},
				{
					ID:   2,
					Name: "Platform 2",
					Slug: "platform-2",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			platformsResources := resources.TransformPlatforms(test.input)

			if platformsResources == nil {
				platformsResources = []resources.PlatformResource{}
			}

			if !reflect.DeepEqual(platformsResources, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, platformsResources)
			}
		})
	}
}
