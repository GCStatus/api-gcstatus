package resources

import "gcstatus/internal/domain"

type DLCStoreResource struct {
	ID         uint          `json:"id"`
	Price      uint          `json:"price"`
	URL        string        `json:"url"`
	Store      StoreResource `json:"store"`
	StoreDLCID string        `json:"store_dlc_id"`
}

func TransformDLCtore(DLCStore domain.DLCStore) DLCStoreResource {
	return DLCStoreResource{
		ID:         DLCStore.ID,
		Price:      DLCStore.Price,
		URL:        DLCStore.URL,
		StoreDLCID: DLCStore.StorDLCID,
		Store:      TransformStore(DLCStore.Store),
	}
}
