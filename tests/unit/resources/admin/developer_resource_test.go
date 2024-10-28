package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformDeveloper(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Developer
		expected resources_admin.DeveloperResource
	}{
		"single Developer": {
			input: domain.Developer{
				ID:        1,
				Name:      "Developer 1",
				Acting:    false,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources_admin.DeveloperResource{
				ID:        1,
				Name:      "Developer 1",
				Acting:    false,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformDeveloper(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestTransformDevelopers(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    []domain.Developer
		expected []resources_admin.DeveloperResource
	}{
		"empty slice": {
			input:    []domain.Developer{},
			expected: []resources_admin.DeveloperResource{},
		},
		"multiple Developers": {
			input: []domain.Developer{
				{
					ID:        1,
					Name:      "Developer 1",
					Acting:    true,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Name:      "Developer 2",
					Acting:    false,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []resources_admin.DeveloperResource{
				{
					ID:        1,
					Name:      "Developer 1",
					Acting:    true,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID:        2,
					Name:      "Developer 2",
					Acting:    false,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformDevelopers(tc.input)

			if result == nil {
				result = []resources_admin.DeveloperResource{}
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
