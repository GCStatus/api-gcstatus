package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformRequirementType(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	tests := map[string]struct {
		input    domain.RequirementType
		expected resources_admin.RequirementTypeResource
	}{
		"as null": {
			input: domain.RequirementType{},
			expected: resources_admin.RequirementTypeResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"multiple categories": {
			input: domain.RequirementType{
				ID:        1,
				Potential: domain.MinimumRequirementType,
				OS:        domain.WindowsOSRequirement,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources_admin.RequirementTypeResource{
				ID:        1,
				Potential: domain.MinimumRequirementType,
				OS:        domain.WindowsOSRequirement,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			RequirementTypeResource := resources_admin.TransformRequirementType(test.input)

			if !reflect.DeepEqual(RequirementTypeResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, RequirementTypeResource)
			}
		})
	}
}

func TestTransformRequirementTypes(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    []domain.RequirementType
		expected []resources_admin.RequirementTypeResource
	}{
		"as null": {
			input:    []domain.RequirementType{},
			expected: []resources_admin.RequirementTypeResource{},
		},
		"multiple categories": {
			input: []domain.RequirementType{
				{
					ID:        1,
					Potential: domain.MinimumRequirementType,
					OS:        domain.WindowsOSRequirement,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Potential: domain.RecommendedRequirementType,
					OS:        domain.WindowsOSRequirement,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []resources_admin.RequirementTypeResource{
				{
					ID:        1,
					Potential: domain.MinimumRequirementType,
					OS:        domain.WindowsOSRequirement,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID:        2,
					Potential: domain.RecommendedRequirementType,
					OS:        domain.WindowsOSRequirement,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			RequirementTypesResources_admin := resources_admin.TransformRequirementTypes(test.input)

			if RequirementTypesResources_admin == nil {
				RequirementTypesResources_admin = []resources_admin.RequirementTypeResource{}
			}

			if !reflect.DeepEqual(RequirementTypesResources_admin, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, RequirementTypesResources_admin)
			}
		})
	}
}
