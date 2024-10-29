package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformGameStore(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.GameStore
		expected resources_admin.GameStoreResource
	}{
		"as nil": {
			input: domain.GameStore{},
			expected: resources_admin.GameStoreResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
				Store: resources_admin.StoreResource{
					CreatedAt: utils.FormatTimestamp(zeroTime),
					UpdatedAt: utils.FormatTimestamp(zeroTime),
				},
			},
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
			expected: resources_admin.GameStoreResource{
				ID:        1,
				Price:     22999,
				URL:       "https://google.com",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				Store: resources_admin.StoreResource{
					ID:        1,
					Name:      "Store 1",
					Slug:      "store-1",
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
			result := resources_admin.TransformGameStore(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
