package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
)

func TestTransformProtection(t *testing.T) {
	testCases := map[string]struct {
		input    domain.Protection
		expected *resources.ProtectionResource
	}{
		"single protection": {
			input: domain.Protection{
				ID:   1,
				Name: "Protection 1",
			},
			expected: &resources.ProtectionResource{
				ID:   1,
				Name: "Protection 1",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformProtection(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestTransformProtections(t *testing.T) {
	testCases := map[string]struct {
		input    []domain.Protection
		expected []*resources.ProtectionResource
	}{
		"empty slice": {
			input:    []domain.Protection{},
			expected: []*resources.ProtectionResource{},
		},
		"multiple protections": {
			input: []domain.Protection{
				{
					ID:   1,
					Name: "Protection 1",
				},
				{
					ID:   2,
					Name: "Protection 2",
				},
			},
			expected: []*resources.ProtectionResource{
				{
					ID:   1,
					Name: "Protection 1",
				},
				{
					ID:   2,
					Name: "Protection 2",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformProtections(tc.input)

			if result == nil {
				result = []*resources.ProtectionResource{}
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
