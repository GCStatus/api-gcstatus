package domain

import (
	"time"

	"gorm.io/gorm"
)

type GameSupport struct {
	gorm.Model
	ID        uint    `gorm:"primaryKey"`
	URL       *string `gorm:"size:255"`
	Email     *string `gorm:"size:255" validate:"email"`
	Contact   *string `gorm:"size:255"`
	GameID    uint    `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Game      Game    `gorm:"foreignKey:GameID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (gs *GameSupport) ValidateGameSupport() error {
	Init()

	err := validate.Struct(gs)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
