package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformDLC(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		input    domain.DLC
		expected resources.DLCResource
	}{
		"empty relations": {
			input: domain.DLC{
				ID:          1,
				Name:        "DLC 1",
				Cover:       "photo-key-1",
				ReleaseDate: fixedTime,
				Galleries:   []domain.Galleriable{},
				Platforms:   []domain.Platformable{},
				Stores:      []domain.DLCStore{},
			},
			expected: resources.DLCResource{
				ID:          1,
				Name:        "DLC 1",
				Cover:       "photo-key-1",
				ReleaseDate: utils.FormatTimestamp(fixedTime),
				Galleries:   []resources.GalleriableResource{},
				Platforms:   []resources.PlatformResource{},
				Stores:      []resources.DLCStoreResource{},
			},
		},
		"fully relations": {
			input: domain.DLC{
				ID:          1,
				Name:        "DLC 1",
				Cover:       "photo-key-1",
				ReleaseDate: fixedTime,
				Galleries: []domain.Galleriable{
					{
						ID:              1,
						S3:              false,
						Path:            "https://google.com",
						GalleriableID:   1,
						GalleriableType: "dlcs",
					},
				},
				Platforms: []domain.Platformable{
					{
						ID:               1,
						PlatformableID:   1,
						PlatformableType: "dlcs",
						PlatformID:       1,
						Platform: domain.Platform{
							ID:   1,
							Name: "Platform 1",
						},
					},
				},
				Stores: []domain.DLCStore{
					{
						ID:      1,
						Price:   2200,
						URL:     "https://google.com",
						DLCID:   1,
						StoreID: 1,
						Store: domain.Store{
							ID:   1,
							Name: "Store 1",
							URL:  "https://google.com",
							Slug: "store-1",
							Logo: "https://google.com",
						},
					},
				},
			},
			expected: resources.DLCResource{
				ID:          1,
				Name:        "DLC 1",
				Cover:       "photo-key-1",
				ReleaseDate: utils.FormatTimestamp(fixedTime),
				Galleries: []resources.GalleriableResource{
					{
						ID:   1,
						Path: "https://google.com",
					},
				},
				Platforms: []resources.PlatformResource{
					{
						ID:   1,
						Name: "Platform 1",
					},
				},
				Stores: []resources.DLCStoreResource{
					{
						ID:    1,
						Price: 2200,
						URL:   "https://google.com",
						Store: resources.StoreResource{
							ID:   1,
							Name: "Store 1",
							URL:  "https://google.com",
							Slug: "store-1",
							Logo: "https://google.com",
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformDLC(tc.input, &MockS3Client{})

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
