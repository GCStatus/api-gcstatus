package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type GameStoreResource struct {
	ID          uint          `json:"id"`
	Price       uint          `json:"price"`
	URL         string        `json:"url"`
	StoreGameID string        `json:"store_game_id"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
	Store       StoreResource `json:"store"`
}

func TransformGameStore(gameStore domain.GameStore) GameStoreResource {
	return GameStoreResource{
		ID:          gameStore.ID,
		Price:       gameStore.Price,
		URL:         gameStore.URL,
		StoreGameID: gameStore.StoreGameID,
		CreatedAt:   utils.FormatTimestamp(gameStore.CreatedAt),
		UpdatedAt:   utils.FormatTimestamp(gameStore.UpdatedAt),
		Store:       TransformStore(gameStore.Store),
	}
}
