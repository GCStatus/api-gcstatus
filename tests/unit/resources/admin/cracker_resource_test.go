package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformCracker(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.Cracker
		expected *resources_admin.CrackerResource
	}{
		"as nil": {
			input: domain.Cracker{},
			expected: &resources_admin.CrackerResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"single Cracker": {
			input: domain.Cracker{
				ID:        1,
				Name:      "Cracker 1",
				Acting:    false,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: &resources_admin.CrackerResource{
				ID:        1,
				Name:      "Cracker 1",
				Acting:    false,
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformCracker(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestTransformCrackers(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    []domain.Cracker
		expected []*resources_admin.CrackerResource
	}{
		"empty slice": {
			input:    []domain.Cracker{},
			expected: []*resources_admin.CrackerResource{},
		},
		"multiple Crackers": {
			input: []domain.Cracker{
				{
					ID:        1,
					Name:      "Cracker 1",
					Acting:    true,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Name:      "Cracker 2",
					Acting:    false,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []*resources_admin.CrackerResource{
				{
					ID:        1,
					Name:      "Cracker 1",
					Acting:    true,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID:        2,
					Name:      "Cracker 2",
					Acting:    false,
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformCrackers(tc.input)

			if result == nil {
				result = []*resources_admin.CrackerResource{}
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
