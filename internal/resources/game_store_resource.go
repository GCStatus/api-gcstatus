package resources

import "gcstatus/internal/domain"

type GameStoreResource struct {
	ID    uint          `json:"id"`
	Price uint          `json:"price"`
	URL   string        `json:"url"`
	Store StoreResource `json:"store"`
}

func TransformGameStore(gameStore domain.GameStore) GameStoreResource {
	return GameStoreResource{
		ID:    gameStore.ID,
		Price: gameStore.Price,
		URL:   gameStore.URL,
		Store: TransformStore(gameStore.Store),
	}
}
