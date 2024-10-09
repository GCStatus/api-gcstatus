package resources

import "gcstatus/internal/domain"

type StoreResource struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	URL  string `json:"url"`
	Logo string `json:"logo"`
}

func TransformStore(store domain.Store) StoreResource {
	return StoreResource{
		ID:   store.ID,
		Name: store.Name,
		Slug: store.Slug,
		URL:  store.URL,
		Logo: store.Logo,
	}
}
