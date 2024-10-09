package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/pkg/utils"
	"reflect"
	"testing"
)

func TestTransformRequirement(t *testing.T) {
	testCases := map[string]struct {
		input    domain.Requirement
		expected resources.RequirementResource
	}{
		"basic transformation": {
			input: domain.Requirement{
				ID:      1,
				OS:      "Windows 10",
				DX:      "12",
				CPU:     "Intel i5",
				RAM:     "8 GB",
				GPU:     "NVIDIA GTX 1060",
				ROM:     "50 GB",
				OBS:     nil,
				Network: "Broadband",
				RequirementType: domain.RequirementType{
					ID:        1,
					Potential: domain.MinimumRequirementType,
					OS:        domain.WindowsOSRequirement,
				},
			},
			expected: resources.RequirementResource{
				ID:      1,
				OS:      "Windows 10",
				DX:      "12",
				CPU:     "Intel i5",
				RAM:     "8 GB",
				GPU:     "NVIDIA GTX 1060",
				ROM:     "50 GB",
				OBS:     nil,
				Network: "Broadband",
				RequirementType: resources.RequirementTypeResource{
					ID:        1,
					Potential: domain.MinimumRequirementType,
					OS:        domain.WindowsOSRequirement,
				},
			},
		},
		"with OBS field": {
			input: domain.Requirement{
				ID:      2,
				OS:      "Windows 11",
				DX:      "12",
				CPU:     "AMD Ryzen 5",
				RAM:     "16 GB",
				GPU:     "NVIDIA RTX 2060",
				ROM:     "100 GB",
				OBS:     utils.StringPtr("Streamable"),
				Network: "High-Speed",
				RequirementType: domain.RequirementType{
					ID:        2,
					Potential: domain.RecommendedRequirementType,
					OS:        domain.WindowsOSRequirement,
				},
			},
			expected: resources.RequirementResource{
				ID:      2,
				OS:      "Windows 11",
				DX:      "12",
				CPU:     "AMD Ryzen 5",
				RAM:     "16 GB",
				GPU:     "NVIDIA RTX 2060",
				ROM:     "100 GB",
				OBS:     utils.StringPtr("Streamable"),
				Network: "High-Speed",
				RequirementType: resources.RequirementTypeResource{
					ID:        2,
					Potential: domain.RecommendedRequirementType,
					OS:        domain.WindowsOSRequirement,
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformRequirement(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
