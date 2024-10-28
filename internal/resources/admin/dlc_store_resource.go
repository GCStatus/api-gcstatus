package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type DLCStoreResource struct {
	ID         uint          `json:"id"`
	Price      uint          `json:"price"`
	URL        string        `json:"url"`
	StoreDLCID string        `json:"store_dlc_id"`
	CreatedAt  string        `json:"created_at"`
	UpdatedAt  string        `json:"updated_at"`
	Store      StoreResource `json:"store"`
}

func TransformDLCtore(DLCStore domain.DLCStore) DLCStoreResource {
	return DLCStoreResource{
		ID:         DLCStore.ID,
		Price:      DLCStore.Price,
		URL:        DLCStore.URL,
		StoreDLCID: DLCStore.StorDLCID,
		CreatedAt:  utils.FormatTimestamp(DLCStore.CreatedAt),
		UpdatedAt:  utils.FormatTimestamp(DLCStore.UpdatedAt),
		Store:      TransformStore(DLCStore.Store),
	}
}
