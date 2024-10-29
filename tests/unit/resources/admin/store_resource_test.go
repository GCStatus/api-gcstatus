package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformStore(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.Store
		expected resources_admin.StoreResource
	}{
		"as nil": {
			input: domain.Store{},
			expected: resources_admin.StoreResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
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
			expected: resources_admin.StoreResource{
				ID:        1,
				Name:      "Store 1",
				Slug:      "store-1",
				URL:       "https://google.com",
				Logo:      "https://placehold.co/600x400/EEE/31343C",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformStore(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
