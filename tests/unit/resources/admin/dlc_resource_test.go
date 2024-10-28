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

func TestTransformDLC(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.DLC
		expected resources_admin.DLCResource
	}{
		"empty relations": {
			input: domain.DLC{
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
			expected: resources_admin.DLCResource{
				ID:          1,
				Name:        "DLC 1",
				Cover:       "photo-key-1",
				ReleaseDate: utils.FormatTimestamp(fixedTime),
				CreatedAt:   utils.FormatTimestamp(fixedTime),
				UpdatedAt:   utils.FormatTimestamp(fixedTime),
				Galleries:   []resources_admin.GalleriableResource{},
				Platforms:   []resources_admin.PlatformResource{},
				Stores:      []resources_admin.DLCStoreResource{},
			},
		},
		"fully relations": {
			input: domain.DLC{
				ID:          1,
				Name:        "DLC 1",
				Cover:       "photo-key-1",
				ReleaseDate: fixedTime,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				Galleries: []domain.Galleriable{
					{
						ID:              1,
						S3:              false,
						Path:            "https://google.com",
						GalleriableID:   1,
						GalleriableType: "dlcs",
						CreatedAt:       fixedTime,
						UpdatedAt:       fixedTime,
					},
				},
				Platforms: []domain.Platformable{
					{
						ID:               1,
						PlatformableID:   1,
						PlatformableType: "dlcs",
						PlatformID:       1,
						CreatedAt:        fixedTime,
						UpdatedAt:        fixedTime,
						Platform: domain.Platform{
							ID:        1,
							Name:      "Platform 1",
							Slug:      "platform-1",
							CreatedAt: fixedTime,
							UpdatedAt: fixedTime,
						},
					},
				},
				Stores: []domain.DLCStore{
					{
						ID:        1,
						Price:     2200,
						URL:       "https://google.com",
						DLCID:     1,
						StoreID:   1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						Store: domain.Store{
							ID:        1,
							Name:      "Store 1",
							Slug:      "store-1",
							URL:       "https://google.com",
							Logo:      "https://google.com",
							CreatedAt: fixedTime,
							UpdatedAt: fixedTime,
						},
					},
				},
			},
			expected: resources_admin.DLCResource{
				ID:          1,
				Name:        "DLC 1",
				Cover:       "photo-key-1",
				ReleaseDate: utils.FormatTimestamp(fixedTime),
				CreatedAt:   utils.FormatTimestamp(fixedTime),
				UpdatedAt:   utils.FormatTimestamp(fixedTime),
				Galleries: []resources_admin.GalleriableResource{
					{
						ID:        1,
						Path:      "https://google.com",
						CreatedAt: utils.FormatTimestamp(fixedTime),
						UpdatedAt: utils.FormatTimestamp(fixedTime),
						MediaType: resources_admin.MediaTypeResource{
							ID:        0,
							Name:      "",
							CreatedAt: "0001-01-01T00:00:00",
							UpdatedAt: "0001-01-01T00:00:00",
						},
					},
				},
				Platforms: []resources_admin.PlatformResource{
					{
						ID:        1,
						Name:      "Platform 1",
						Slug:      "platform-1",
						CreatedAt: utils.FormatTimestamp(fixedTime),
						UpdatedAt: utils.FormatTimestamp(fixedTime),
					},
				},
				Stores: []resources_admin.DLCStoreResource{
					{
						ID:        1,
						Price:     2200,
						URL:       "https://google.com",
						CreatedAt: utils.FormatTimestamp(fixedTime),
						UpdatedAt: utils.FormatTimestamp(fixedTime),
						Store: resources_admin.StoreResource{
							ID:        1,
							Name:      "Store 1",
							Slug:      "store-1",
							URL:       "https://google.com",
							Logo:      "https://google.com",
							CreatedAt: utils.FormatTimestamp(fixedTime),
							UpdatedAt: utils.FormatTimestamp(fixedTime),
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformDLC(tc.input, &testutils.MockS3Client{})

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
