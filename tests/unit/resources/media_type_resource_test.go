package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
	"time"
)

func TestTransformMediaType(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    domain.MediaType
		expected resources.MediaTypeResource
	}{
		"as null": {
			input:    domain.MediaType{},
			expected: resources.MediaTypeResource{},
		},
		"valid category": {
			input: domain.MediaType{
				ID:        1,
				Name:      "photo",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources.MediaTypeResource{
				ID:   1,
				Name: "photo",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			MediaTypeResource := resources.TransformMediaType(test.input)

			if !reflect.DeepEqual(MediaTypeResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, MediaTypeResource)
			}
		})
	}
}
