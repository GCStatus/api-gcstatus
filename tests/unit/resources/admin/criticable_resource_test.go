package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestTransformCriticable(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.Criticable
		expected resources_admin.CriticableResource
	}{
		"as nil": {
			input: domain.Criticable{},
			expected: resources_admin.CriticableResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
				Critic: resources_admin.CriticResource{
					CreatedAt: utils.FormatTimestamp(zeroTime),
					UpdatedAt: utils.FormatTimestamp(zeroTime),
				},
			},
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
					ID:        1,
					Name:      "Criticable 1",
					URL:       "https://google.com",
					Logo:      "https://placehold.co/600x400/EEE/31343C",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: resources_admin.CriticableResource{
				ID:        1,
				Rate:      decimal.NewFromFloat32(5.8),
				URL:       "https://google.com",
				PostedAt:  utils.FormatTimestamp(fixedTime),
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				Critic: resources_admin.CriticResource{
					ID:        1,
					Name:      "Criticable 1",
					URL:       "https://google.com",
					Logo:      "https://placehold.co/600x400/EEE/31343C",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformCriticable(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
