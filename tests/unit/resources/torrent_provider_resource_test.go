package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"reflect"
	"testing"
)

func TestTransformTorrentProvider(t *testing.T) {
	testCases := map[string]struct {
		input    domain.TorrentProvider
		expected resources.TorrentProviderResource
	}{
		"as null": {
			input:    domain.TorrentProvider{},
			expected: resources.TorrentProviderResource{},
		},
		"valid torrent provider": {
			input: domain.TorrentProvider{
				ID:   1,
				Name: "Google",
				URL:  "https://google.com",
			},
			expected: resources.TorrentProviderResource{
				ID:   1,
				Name: "Google",
				URL:  "https://google.com",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := resources.TransformTorrentProvider(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
