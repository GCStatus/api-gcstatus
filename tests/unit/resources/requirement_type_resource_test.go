package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
)

func TestTransformRequirementType(t *testing.T) {
	tests := map[string]struct {
		input    domain.RequirementType
		expected resources.RequirementTypeResource
	}{
		"as null": {
			input:    domain.RequirementType{},
			expected: resources.RequirementTypeResource{},
		},
		"multiple categories": {
			input: domain.RequirementType{
				ID:        1,
				Potential: domain.MinimumRequirementType,
				OS:        domain.WindowsOSRequirement,
			},
			expected: resources.RequirementTypeResource{
				ID:        1,
				Potential: domain.MinimumRequirementType,
				OS:        domain.WindowsOSRequirement,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			RequirementTypeResource := resources.TransformRequirementType(test.input)

			if !reflect.DeepEqual(RequirementTypeResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, RequirementTypeResource)
			}
		})
	}
}

func TestTransformRequirementTypes(t *testing.T) {
	tests := map[string]struct {
		input    []domain.RequirementType
		expected []resources.RequirementTypeResource
	}{
		"as null": {
			input:    []domain.RequirementType{},
			expected: []resources.RequirementTypeResource{},
		},
		"multiple categories": {
			input: []domain.RequirementType{
				{
					ID:        1,
					Potential: domain.MinimumRequirementType,
					OS:        domain.WindowsOSRequirement,
				},
				{
					ID:        2,
					Potential: domain.RecommendedRequirementType,
					OS:        domain.WindowsOSRequirement,
				},
			},
			expected: []resources.RequirementTypeResource{
				{
					ID:        1,
					Potential: domain.MinimumRequirementType,
					OS:        domain.WindowsOSRequirement,
				},
				{
					ID:        2,
					Potential: domain.RecommendedRequirementType,
					OS:        domain.WindowsOSRequirement,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			RequirementTypesResources := resources.TransformRequirementTypes(test.input)

			if RequirementTypesResources == nil {
				RequirementTypesResources = []resources.RequirementTypeResource{}
			}

			if !reflect.DeepEqual(RequirementTypesResources, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, RequirementTypesResources)
			}
		})
	}
}
