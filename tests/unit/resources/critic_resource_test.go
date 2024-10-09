package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
	"time"
)

func TestTransformCritic(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Critic
		expected resources.CriticResource
	}{
		"as nil": {
			input:    domain.Critic{},
			expected: resources.CriticResource{},
		},
		"basic transformation": {
			input: domain.Critic{
				ID:        1,
				Name:      "Critic 1",
				URL:       "https://google.com",
				Logo:      "https://placehold.co/600x400/EEE/31343C",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources.CriticResource{
				ID:   1,
				Name: "Critic 1",
				URL:  "https://google.com",
				Logo: "https://placehold.co/600x400/EEE/31343C",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformCritic(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
