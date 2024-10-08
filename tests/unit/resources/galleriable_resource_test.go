package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
)

func TestTransformGalleriable(t *testing.T) {
	testCases := map[string]struct {
		input    domain.Galleriable
		expected resources.GalleriableResource
	}{
		"as nil": {
			input:    domain.Galleriable{},
			expected: resources.GalleriableResource{},
		},
		"no s3": {
			input: domain.Galleriable{
				ID:              1,
				S3:              false,
				Path:            "https://placehold.co/600x400/EEE/31343C",
				GalleriableID:   1,
				GalleriableType: "games",
			},
			expected: resources.GalleriableResource{
				ID:   1,
				Path: "https://placehold.co/600x400/EEE/31343C",
			},
		},
		"as s3": {
			input: domain.Galleriable{
				ID:              1,
				S3:              true,
				Path:            "photo-key-1",
				GalleriableID:   1,
				GalleriableType: "games",
			},
			expected: resources.GalleriableResource{
				ID:   1,
				Path: "https://mock-presigned-url.com/photo-key-1",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformGalleriable(tc.input, &MockS3Client{})

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
