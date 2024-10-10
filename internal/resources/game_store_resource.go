package resources

import "gcstatus/internal/domain"

type GameStoreResource struct {
	ID          uint          `json:"id"`
	Price       uint          `json:"price"`
	URL         string        `json:"url"`
	Store       StoreResource `json:"store"`
	StoreGameID string        `json:"store_game_id"`
}

func TransformGameStore(gameStore domain.GameStore) GameStoreResource {
	return GameStoreResource{
		ID:          gameStore.ID,
		Price:       gameStore.Price,
		URL:         gameStore.URL,
		StoreGameID: gameStore.StoreGameID,
		Store:       TransformStore(gameStore.Store),
	}
}
