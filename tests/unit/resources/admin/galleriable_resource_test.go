package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	testutils "gcstatus/tests/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformGalleriable(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.Galleriable
		expected resources_admin.GalleriableResource
	}{
		"as nil": {
			input: domain.Galleriable{},
			expected: resources_admin.GalleriableResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
				MediaType: resources_admin.MediaTypeResource{
					CreatedAt: utils.FormatTimestamp(zeroTime),
					UpdatedAt: utils.FormatTimestamp(zeroTime),
				},
			},
		},
		"no s3": {
			input: domain.Galleriable{
				ID:              1,
				S3:              false,
				Path:            "https://placehold.co/600x400/EEE/31343C",
				GalleriableID:   1,
				GalleriableType: "games",
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
				MediaType: domain.MediaType{
					ID:        1,
					Name:      "photo",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: resources_admin.GalleriableResource{
				ID:        1,
				Path:      "https://placehold.co/600x400/EEE/31343C",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				MediaType: resources_admin.MediaTypeResource{
					ID:        1,
					Name:      "photo",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
		"as s3": {
			input: domain.Galleriable{
				ID:              1,
				S3:              true,
				Path:            "photo-key-1",
				GalleriableID:   1,
				GalleriableType: "games",
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
				MediaType: domain.MediaType{
					ID:        2,
					Name:      "video",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: resources_admin.GalleriableResource{
				ID:        1,
				Path:      "https://mock-presigned-url.com/photo-key-1",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				MediaType: resources_admin.MediaTypeResource{
					ID:        2,
					Name:      "video",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformGalleriable(tc.input, &testutils.MockS3Client{})

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
