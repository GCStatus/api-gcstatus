package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformRequirement(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Requirement
		expected resources_admin.RequirementResource
	}{
		"basic transformation": {
			input: domain.Requirement{
				ID:        1,
				OS:        "Windows 10",
				DX:        "12",
				CPU:       "Intel i5",
				RAM:       "8 GB",
				GPU:       "NVIDIA GTX 1060",
				ROM:       "50 GB",
				OBS:       nil,
				Network:   "Broadband",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				RequirementType: domain.RequirementType{
					ID:        1,
					Potential: domain.MinimumRequirementType,
					OS:        domain.WindowsOSRequirement,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: resources_admin.RequirementResource{
				ID:        1,
				OS:        "Windows 10",
				DX:        "12",
				CPU:       "Intel i5",
				RAM:       "8 GB",
				GPU:       "NVIDIA GTX 1060",
				ROM:       "50 GB",
				OBS:       nil,
				Network:   "Broadband",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				RequirementType: resources_admin.RequirementTypeResource{
					ID:        1,
					Potential: domain.MinimumRequirementType,
					OS:        domain.WindowsOSRequirement,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
		"with OBS field": {
			input: domain.Requirement{
				ID:        2,
				OS:        "Windows 11",
				DX:        "12",
				CPU:       "AMD Ryzen 5",
				RAM:       "16 GB",
				GPU:       "NVIDIA RTX 2060",
				ROM:       "100 GB",
				OBS:       utils.StringPtr("Streamable"),
				Network:   "High-Speed",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				RequirementType: domain.RequirementType{
					ID:        2,
					Potential: domain.RecommendedRequirementType,
					OS:        domain.WindowsOSRequirement,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: resources_admin.RequirementResource{
				ID:        2,
				OS:        "Windows 11",
				DX:        "12",
				CPU:       "AMD Ryzen 5",
				RAM:       "16 GB",
				GPU:       "NVIDIA RTX 2060",
				ROM:       "100 GB",
				OBS:       utils.StringPtr("Streamable"),
				Network:   "High-Speed",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				RequirementType: resources_admin.RequirementTypeResource{
					ID:        2,
					Potential: domain.RecommendedRequirementType,
					OS:        domain.WindowsOSRequirement,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformRequirement(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
