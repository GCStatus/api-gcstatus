package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformProtection(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.Protection
		expected *resources_admin.ProtectionResource
	}{
		"as nil": {
			input: domain.Protection{},
			expected: &resources_admin.ProtectionResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"single protection": {
			input: domain.Protection{
				ID:        1,
				Name:      "Protection 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: &resources_admin.ProtectionResource{
				ID:        1,
				Name:      "Protection 1",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformProtection(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestTransformProtections(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    []domain.Protection
		expected []*resources_admin.ProtectionResource
	}{
		"empty slice": {
			input:    []domain.Protection{},
			expected: []*resources_admin.ProtectionResource{},
		},
		"multiple protections": {
			input: []domain.Protection{
				{
					ID:        1,
					Name:      "Protection 1",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				{
					ID:        2,
					Name:      "Protection 2",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: []*resources_admin.ProtectionResource{
				{
					ID:        1,
					Name:      "Protection 1",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
				{
					ID:        2,
					Name:      "Protection 2",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformProtections(tc.input)

			if result == nil {
				result = []*resources_admin.ProtectionResource{}
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
