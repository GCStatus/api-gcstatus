package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type TorrentProviderResource struct {
	ID        uint   `json:"id"`
	URL       string `json:"url"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformTorrentProvider(torrentProvider domain.TorrentProvider) TorrentProviderResource {
	return TorrentProviderResource{
		ID:        torrentProvider.ID,
		URL:       torrentProvider.URL,
		Name:      torrentProvider.Name,
		CreatedAt: utils.FormatTimestamp(torrentProvider.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(torrentProvider.UpdatedAt),
	}
}
