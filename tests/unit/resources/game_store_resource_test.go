package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
	"time"
)

func TestTransformGameStore(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.GameStore
		expected resources.GameStoreResource
	}{
		"as nil": {
			input:    domain.GameStore{},
			expected: resources.GameStoreResource{},
		},
		"basic transformation": {
			input: domain.GameStore{
				ID:        1,
				Price:     22999,
				URL:       "https://google.com",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Game: domain.Game{
					Slug:             "valid",
					Age:              18,
					Title:            "Game Test",
					Condition:        domain.CommomCondition,
					Cover:            "https://placehold.co/600x400/EEE/31343C",
					About:            "About game",
					Description:      "Description",
					ShortDescription: "Short description",
					Free:             false,
					ReleaseDate:      fixedTime,
					CreatedAt:        fixedTime,
					UpdatedAt:        fixedTime,
					Views: []domain.Viewable{
						{
							UserID:       10,
							ViewableID:   1,
							ViewableType: "games",
						},
					},
				},
				Store: domain.Store{
					ID:        1,
					Name:      "Store 1",
					Slug:      "store-1",
					URL:       "https://google.com",
					Logo:      "https://placehold.co/600x400/EEE/31343C",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: resources.GameStoreResource{
				ID:    1,
				Price: 22999,
				URL:   "https://google.com",
				Store: resources.StoreResource{
					ID:   1,
					Name: "Store 1",
					Slug: "store-1",
					URL:  "https://google.com",
					Logo: "https://placehold.co/600x400/EEE/31343C",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformGameStore(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
