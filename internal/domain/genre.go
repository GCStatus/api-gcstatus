package domain

import (
	"time"

	"gorm.io/gorm"
)

type Genre struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null;unique" validate:"required"`
	Slug      string `gorm:"size:255;not null;unique" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (g *Genre) ValidateGenre() error {
	Init()

	err := validate.Struct(g)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
