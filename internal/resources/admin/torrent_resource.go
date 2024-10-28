package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type TorrentResource struct {
	ID        uint                    `json:"id"`
	URL       string                  `json:"url"`
	PostedAt  string                  `json:"posted_in"`
	CreatedAt string                  `json:"created_at"`
	UpdatedAt string                  `json:"updated_at"`
	Provider  TorrentProviderResource `json:"provider"`
}

func TransformTorrent(torrent domain.Torrent) TorrentResource {
	return TorrentResource{
		ID:        torrent.ID,
		URL:       torrent.URL,
		PostedAt:  utils.FormatTimestamp(torrent.PostedAt),
		CreatedAt: utils.FormatTimestamp(torrent.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(torrent.UpdatedAt),
		Provider:  TransformTorrentProvider(torrent.TorrentProvider),
	}
}

func TransformTorrents(torrents []domain.Torrent) []TorrentResource {
	var resources []TorrentResource
	for _, torrent := range torrents {
		resources = append(resources, TransformTorrent(torrent))
	}
	return resources
}
