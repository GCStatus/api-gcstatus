package domain

import (
	"time"

	"gorm.io/gorm"
)

type GameDeveloper struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey"`
	GameID      uint      `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Game        Game      `gorm:"foreignKey:GameID;references:ID"`
	DeveloperID uint      `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Developer   Developer `gorm:"foreignKey:DeveloperID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (gp *GameDeveloper) ValidateGameDeveloper() error {
	Init()

	err := validate.Struct(gp)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
