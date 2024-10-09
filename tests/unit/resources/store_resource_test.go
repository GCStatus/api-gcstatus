package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
	"time"
)

func TestTransformStore(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.Store
		expected resources.StoreResource
	}{
		"as nil": {
			input:    domain.Store{},
			expected: resources.StoreResource{},
		},
		"basic transformation": {
			input: domain.Store{
				ID:        1,
				Name:      "Store 1",
				Slug:      "store-1",
				URL:       "https://google.com",
				Logo:      "https://placehold.co/600x400/EEE/31343C",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources.StoreResource{
				ID:   1,
				Name: "Store 1",
				Slug: "store-1",
				URL:  "https://google.com",
				Logo: "https://placehold.co/600x400/EEE/31343C",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformStore(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
