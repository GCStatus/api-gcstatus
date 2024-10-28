package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
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
		expected *resources_admin.CrackResource
	}{
		"basic transformation": {
			input: &domain.Crack{
				ID:        1,
				Status:    "cracked",
				CrackedAt: &fixedTime,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Cracker: domain.Cracker{
					ID:        1,
					Name:      "Cracker 1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				Protection: domain.Protection{
					ID:        1,
					Name:      "Protection 1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: &resources_admin.CrackResource{
				ID:        1,
				Status:    "cracked",
				CrackedAt: &formattedTime,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				By: &resources_admin.CrackerResource{
					ID:        1,
					Name:      "Cracker 1",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				Protection: &resources_admin.ProtectionResource{
					ID:        1,
					Name:      "Protection 1",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
		"nil CrackedAt field": {
			input: &domain.Crack{
				ID:        2,
				Status:    "uncracked",
				CrackedAt: nil,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Cracker: domain.Cracker{
					ID:        1,
					Name:      "Cracker 2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				Protection: domain.Protection{
					ID:        1,
					Name:      "Protection 2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: &resources_admin.CrackResource{
				ID:        2,
				Status:    "uncracked",
				CrackedAt: nil,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				By: &resources_admin.CrackerResource{
					ID:        1,
					Name:      "Cracker 2",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				Protection: &resources_admin.ProtectionResource{
					ID:        1,
					Name:      "Protection 2",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
		"nil Cracker and Protection": {
			input: &domain.Crack{
				ID:         3,
				Status:     "uncracked",
				CrackedAt:  &fixedTime,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
				Cracker:    domain.Cracker{ID: 0},
				Protection: domain.Protection{ID: 0},
			},
			expected: &resources_admin.CrackResource{
				ID:         3,
				Status:     "uncracked",
				CrackedAt:  &formattedTime,
				CreatedAt:  utils.FormatTimestamp(fixedTime),
				UpdatedAt:  utils.FormatTimestamp(fixedTime),
				By:         nil,
				Protection: nil,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformCrack(tc.input)

			if tc.expected.CrackedAt != nil && result.CrackedAt != nil {
				if *tc.expected.CrackedAt != *result.CrackedAt {
					t.Errorf("Expected CrackedAt %v, got %v", *tc.expected.CrackedAt, *result.CrackedAt)
				}
				result.CrackedAt = nil
				tc.expected.CrackedAt = nil
			}

			result.CreatedAt = utils.FormatTimestamp(fixedTime)
			result.UpdatedAt = utils.FormatTimestamp(fixedTime)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
