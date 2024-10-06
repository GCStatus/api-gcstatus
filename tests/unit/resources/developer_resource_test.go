package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
)

func TestTransformDeveloper(t *testing.T) {
	testCases := map[string]struct {
		input    domain.Developer
		expected resources.DeveloperResource
	}{
		"single Developer": {
			input: domain.Developer{
				ID:     1,
				Name:   "Developer 1",
				Acting: false,
			},
			expected: resources.DeveloperResource{
				ID:     1,
				Name:   "Developer 1",
				Acting: false,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformDeveloper(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestTransformDevelopers(t *testing.T) {
	testCases := map[string]struct {
		input    []domain.Developer
		expected []resources.DeveloperResource
	}{
		"empty slice": {
			input:    []domain.Developer{},
			expected: []resources.DeveloperResource{},
		},
		"multiple Developers": {
			input: []domain.Developer{
				{
					ID:     1,
					Name:   "Developer 1",
					Acting: true,
				},
				{
					ID:     2,
					Name:   "Developer 2",
					Acting: false,
				},
			},
			expected: []resources.DeveloperResource{
				{
					ID:     1,
					Name:   "Developer 1",
					Acting: true,
				},
				{
					ID:     2,
					Name:   "Developer 2",
					Acting: false,
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformDevelopers(tc.input)

			if result == nil {
				result = []resources.DeveloperResource{}
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
