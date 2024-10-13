package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformTorrent(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    domain.Torrent
		expected resources.TorrentResource
	}{
		"valid torrent": {
			input: domain.Torrent{
				ID:        1,
				URL:       "https://google.com",
				PostedAt:  fixedTime,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				TorrentProvider: domain.TorrentProvider{
					ID:        1,
					URL:       "https://google.com",
					Name:      "Google",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			expected: resources.TorrentResource{
				ID:       1,
				URL:      "https://google.com",
				PostedAt: utils.FormatTimestamp(fixedTime),
				Provider: resources.TorrentProviderResource{
					ID:   1,
					URL:  "https://google.com",
					Name: "Google",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			categoryResource := resources.TransformTorrent(test.input)

			if !reflect.DeepEqual(categoryResource, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, categoryResource)
			}
		})
	}
}

func TestTransformTorrents(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		input    []domain.Torrent
		expected []resources.TorrentResource
	}{
		"as null": {
			input:    []domain.Torrent{},
			expected: []resources.TorrentResource{},
		},
		"multiple categories": {
			input: []domain.Torrent{
				{
					ID:        1,
					URL:       "https://google.com",
					PostedAt:  fixedTime,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					TorrentProvider: domain.TorrentProvider{
						ID:        1,
						URL:       "https://google.com",
						Name:      "Google",
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					},
				},
				{
					ID:        2,
					URL:       "https://google2.com",
					PostedAt:  fixedTime,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					TorrentProvider: domain.TorrentProvider{
						ID:        2,
						URL:       "https://google2.com",
						Name:      "Google2",
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
					},
				},
			},
			expected: []resources.TorrentResource{
				{
					ID:       1,
					URL:      "https://google.com",
					PostedAt: utils.FormatTimestamp(fixedTime),
					Provider: resources.TorrentProviderResource{
						ID:   1,
						URL:  "https://google.com",
						Name: "Google",
					},
				},
				{
					ID:       2,
					URL:      "https://google2.com",
					PostedAt: utils.FormatTimestamp(fixedTime),
					Provider: resources.TorrentProviderResource{
						ID:   2,
						URL:  "https://google2.com",
						Name: "Google2",
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			torrentsResources := resources.TransformTorrents(test.input)

			if torrentsResources == nil {
				torrentsResources = []resources.TorrentResource{}
			}

			if !reflect.DeepEqual(torrentsResources, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, torrentsResources)
			}
		})
	}
}
