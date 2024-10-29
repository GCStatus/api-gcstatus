package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformTorrentProvider(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	testCases := map[string]struct {
		input    domain.TorrentProvider
		expected resources_admin.TorrentProviderResource
	}{
		"as null": {
			input: domain.TorrentProvider{},
			expected: resources_admin.TorrentProviderResource{
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
			},
		},
		"valid torrent provider": {
			input: domain.TorrentProvider{
				ID:        1,
				Name:      "Google",
				URL:       "https://google.com",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expected: resources_admin.TorrentProviderResource{
				ID:        1,
				Name:      "Google",
				URL:       "https://google.com",
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources_admin.TransformTorrentProvider(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
