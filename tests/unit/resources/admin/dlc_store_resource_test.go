package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformDLCStore(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.DLCStore
		expected resources_admin.DLCStoreResource
	}{
		"as nil": {
			input: domain.DLCStore{},
			expected: resources_admin.DLCStoreResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
				Store: resources_admin.StoreResource{
					CreatedAt: utils.FormatTimestamp(zeroTime),
					UpdatedAt: utils.FormatTimestamp(zeroTime),
				},
			},
		},
		"basic transformation": {
			input: domain.DLCStore{
				ID:        1,
				Price:     22999,
				URL:       "https://google.com",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				DLC: domain.DLC{
					ID:          1,
					Name:        "DLC 1",
					Cover:       "photo-key-1",
					ReleaseDate: fixedTime,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
					Galleries:   []domain.Galleriable{},
					Platforms:   []domain.Platformable{},
					Stores:      []domain.DLCStore{},
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
			expected: resources_admin.DLCStoreResource{
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
			result := resources_admin.TransformDLCtore(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
