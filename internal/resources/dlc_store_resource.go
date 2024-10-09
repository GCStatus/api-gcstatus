package resources

import "gcstatus/internal/domain"

type DLCStoreResource struct {
	ID    uint          `json:"id"`
	Price uint          `json:"price"`
	URL   string        `json:"url"`
	Store StoreResource `json:"store"`
}

func TransformDLCtore(DLCStore domain.DLCStore) DLCStoreResource {
	return DLCStoreResource{
		ID:    DLCStore.ID,
		Price: DLCStore.Price,
		URL:   DLCStore.URL,
		Store: TransformStore(DLCStore.Store),
	}
}
