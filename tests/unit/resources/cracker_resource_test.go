package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
)

func TestTransformCracker(t *testing.T) {
	testCases := map[string]struct {
		input    domain.Cracker
		expected *resources.CrackerResource
	}{
		"single Cracker": {
			input: domain.Cracker{
				ID:   1,
				Name: "Cracker 1",
			},
			expected: &resources.CrackerResource{
				ID:   1,
				Name: "Cracker 1",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformCracker(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestTransformCrackers(t *testing.T) {
	testCases := map[string]struct {
		input    []domain.Cracker
		expected []*resources.CrackerResource
	}{
		"empty slice": {
			input:    []domain.Cracker{},
			expected: []*resources.CrackerResource{},
		},
		"multiple Crackers": {
			input: []domain.Cracker{
				{
					ID:   1,
					Name: "Cracker 1",
				},
				{
					ID:   2,
					Name: "Cracker 2",
				},
			},
			expected: []*resources.CrackerResource{
				{
					ID:   1,
					Name: "Cracker 1",
				},
				{
					ID:   2,
					Name: "Cracker 2",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformCrackers(tc.input)

			if result == nil {
				result = []*resources.CrackerResource{}
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
