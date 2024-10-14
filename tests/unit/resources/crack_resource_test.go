package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformCrack(t *testing.T) {
	fixedTime := time.Now()
	formattedTime := utils.FormatTimestamp(fixedTime)

	testCases := map[string]struct {
		input    *domain.Crack
		expected *resources.CrackResource
	}{
		"basic transformation": {
			input: &domain.Crack{
				ID:        1,
				Status:    "cracked",
				CrackedAt: &fixedTime,
				Cracker: domain.Cracker{
					ID:   1,
					Name: "Cracker 1",
				},
				Protection: domain.Protection{
					ID:   1,
					Name: "Protection 1",
				},
			},
			expected: &resources.CrackResource{
				ID:         1,
				Status:     "cracked",
				CrackedAt:  &formattedTime,
				By:         &resources.CrackerResource{ID: 1, Name: "Cracker 1"},
				Protection: &resources.ProtectionResource{ID: 1, Name: "Protection 1"},
			},
		},
		"nil CrackedAt field": {
			input: &domain.Crack{
				ID:         2,
				Status:     "uncracked",
				CrackedAt:  nil,
				Cracker:    domain.Cracker{ID: 1, Name: "Cracker 2"},
				Protection: domain.Protection{ID: 1, Name: "Protection 2"},
			},
			expected: &resources.CrackResource{
				ID:         2,
				Status:     "uncracked",
				CrackedAt:  nil,
				By:         &resources.CrackerResource{ID: 1, Name: "Cracker 2"},
				Protection: &resources.ProtectionResource{ID: 1, Name: "Protection 2"},
			},
		},
		"nil Cracker and Protection": {
			input: &domain.Crack{
				ID:         3,
				Status:     "uncracked",
				CrackedAt:  &fixedTime,
				Cracker:    domain.Cracker{ID: 0},
				Protection: domain.Protection{ID: 0},
			},
			expected: &resources.CrackResource{
				ID:         3,
				Status:     "uncracked",
				CrackedAt:  &formattedTime,
				By:         nil,
				Protection: nil,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformCrack(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
