package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/pkg/utils"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestTransformCriticable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Criticable
		expected resources.CriticableResource
	}{
		"as nil": {
			input:    domain.Criticable{},
			expected: resources.CriticableResource{},
		},
		"basic transformation": {
			input: domain.Criticable{
				ID:             1,
				URL:            "https://google.com",
				Rate:           decimal.NewFromFloat32(5.8),
				PostedAt:       fixedTime,
				CriticableID:   1,
				CriticableType: "games",
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
				Critic: domain.Critic{
					ID:   1,
					Name: "Criticable 1",
					URL:  "https://google.com",
					Logo: "https://placehold.co/600x400/EEE/31343C",
				},
			},
			expected: resources.CriticableResource{
				ID:       1,
				Rate:     decimal.NewFromFloat32(5.8),
				URL:      "https://google.com",
				PostedAt: utils.FormatTimestamp(fixedTime),
				Critic: resources.CriticResource{
					ID:   1,
					Name: "Criticable 1",
					URL:  "https://google.com",
					Logo: "https://placehold.co/600x400/EEE/31343C",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformCriticable(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
