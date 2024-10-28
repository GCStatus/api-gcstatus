package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestTransformTorrent(t *testing.T) {
	fixedTime := time.Now()
	zeroTime := time.Time{}

	tests := map[string]struct {
		input    domain.Torrent
		expected resources_admin.TorrentResource
	}{
		"as nil": {
			input: domain.Torrent{},
			expected: resources_admin.TorrentResource{
				PostedAt:  utils.FormatTimestamp(zeroTime),
				CreatedAt: utils.FormatTimestamp(zeroTime),
				UpdatedAt: utils.FormatTimestamp(zeroTime),
				Provider: resources_admin.TorrentProviderResource{
					CreatedAt: utils.FormatTimestamp(zeroTime),
					UpdatedAt: utils.FormatTimestamp(zeroTime),
				},
			},
		},
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
			expected: resources_admin.TorrentResource{
				ID:        1,
				URL:       "https://google.com",
				PostedAt:  utils.FormatTimestamp(fixedTime),
				CreatedAt: utils.FormatTimestamp(fixedTime),
				UpdatedAt: utils.FormatTimestamp(fixedTime),
				Provider: resources_admin.TorrentProviderResource{
					ID:        1,
					URL:       "https://google.com",
					Name:      "Google",
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			categoryResource := resources_admin.TransformTorrent(test.input)

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
		expected []resources_admin.TorrentResource
	}{
		"as null": {
			input:    []domain.Torrent{},
			expected: []resources_admin.TorrentResource{},
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
			expected: []resources_admin.TorrentResource{
				{
					ID:        1,
					URL:       "https://google.com",
					PostedAt:  utils.FormatTimestamp(fixedTime),
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
					Provider: resources_admin.TorrentProviderResource{
						ID:        1,
						URL:       "https://google.com",
						Name:      "Google",
						CreatedAt: utils.FormatTimestamp(fixedTime),
						UpdatedAt: utils.FormatTimestamp(fixedTime),
					},
				},
				{
					ID:        2,
					URL:       "https://google2.com",
					PostedAt:  utils.FormatTimestamp(fixedTime),
					CreatedAt: utils.FormatTimestamp(fixedTime),
					UpdatedAt: utils.FormatTimestamp(fixedTime),
					Provider: resources_admin.TorrentProviderResource{
						ID:        2,
						URL:       "https://google2.com",
						Name:      "Google2",
						CreatedAt: utils.FormatTimestamp(fixedTime),
						UpdatedAt: utils.FormatTimestamp(fixedTime),
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			torrentsResources_admin := resources_admin.TransformTorrents(test.input)

			if torrentsResources_admin == nil {
				torrentsResources_admin = []resources_admin.TorrentResource{}
			}

			if !reflect.DeepEqual(torrentsResources_admin, test.expected) {
				t.Errorf("Expected %+v, got %+v", test.expected, torrentsResources_admin)
			}
		})
	}
}
