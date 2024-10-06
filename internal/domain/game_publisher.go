package domain

import (
	"time"

	"gorm.io/gorm"
)

type GamePublisher struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey"`
	GameID      uint      `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Game        Game      `gorm:"foreignKey:GameID;references:ID"`
	PublisherID uint      `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Publisher   Publisher `gorm:"foreignKey:PublisherID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (gp *GamePublisher) ValidateGamePublisher() error {
	Init()

	err := validate.Struct(gp)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
