package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
	"time"
)

func TestTransformDLCStore(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.DLCStore
		expected resources.DLCStoreResource
	}{
		"as nil": {
			input:    domain.DLCStore{},
			expected: resources.DLCStoreResource{},
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
			expected: resources.DLCStoreResource{
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
			result := resources.TransformDLCtore(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
