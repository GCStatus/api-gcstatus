package domain

import (
	"time"

	"gorm.io/gorm"
)

type DLC struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:255;not null" validate:"required"`
	Cover       string `gorm:"size:255;not null" validate:"required"`
	ReleaseDate time.Time
	GameID      uint           `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Game        Game           `gorm:"foreignKey:GameID;references:ID"`
	Galleries   []Galleriable  `gorm:"polymorphic:Galleriable"`
	Platforms   []Platformable `gorm:"polymorphic:Platformable;"`
	Stores      []DLCStore     `gorm:"foreignKey:DLCID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (d *DLC) ValidateDLC() error {
	Init()

	if err := validate.Struct(d); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
