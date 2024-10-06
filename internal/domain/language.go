package domain

import (
	"time"

	"gorm.io/gorm"
)

type Language struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null" validate:"required"`
	ISO       string `gorm:"size:10;not null" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (g *Language) ValidateLanguage() error {
	Init()

	err := validate.Struct(g)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
