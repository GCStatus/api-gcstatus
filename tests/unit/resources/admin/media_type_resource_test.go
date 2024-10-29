package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformMediaType(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	tests := map[string]struct {
		input    domain.MediaType
		expected resources_admin.MediaTypeResource
	}{
		"as null": {
			input: domain.MediaType{},
			expected: resources_admin.MediaTypeResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"valid category": {
			input: domain.MediaType{
				ID:        1,
				Name:      "photo",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources_admin.MediaTypeResource{
				ID:        1,
				Name:      "photo",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			MediaTypeResource := resources_admin.TransformMediaType(test.input)

			if !reflect.DeepEqual(MediaTypeResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, MediaTypeResource)
			}
		})
	}
}
