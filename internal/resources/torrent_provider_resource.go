package resources

import "gcstatus/internal/domain"

type TorrentProviderResource struct {
	ID   uint   `json:"id"`
	URL  string `json:"url"`
	Name string `json:"name"`
}

func TransformTorrentProvider(torrentProvider domain.TorrentProvider) TorrentProviderResource {
	return TorrentProviderResource{
		ID:   torrentProvider.ID,
		URL:  torrentProvider.URL,
		Name: torrentProvider.Name,
	}
}
