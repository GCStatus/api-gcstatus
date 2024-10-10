package domain

import (
	"time"

	"gorm.io/gorm"
)

type GameStore struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Price       uint   `gorm:"not null" validate:"required"`
	URL         string `gorm:"size:255;not null" validate:"required"`
	GameID      uint   `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Game        Game   `gorm:"foreignKey:GameID;references:ID"`
	StoreID     uint   `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Store       Store  `gorm:"foreignKey:StoreID;references:ID"`
	StoreGameID string `gorm:"not null;" validate:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (gs *GameStore) ValidateGameStore() error {
	Init()

	if err := validate.Struct(gs); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
