package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type StoreResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	URL       string `json:"url"`
	Logo      string `json:"logo"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformStore(store domain.Store) StoreResource {
	return StoreResource{
		ID:        store.ID,
		Name:      store.Name,
		Slug:      store.Slug,
		URL:       store.URL,
		Logo:      store.Logo,
		CreatedAt: utils.FormatTimestamp(store.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(store.UpdatedAt),
	}
}
