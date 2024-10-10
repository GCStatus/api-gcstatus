package domain

import (
	"time"

	"gorm.io/gorm"
)

type MediaType struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null;unique" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (mt *MediaType) ValidateMediaType() error {
	Init()

	if err := validate.Struct(mt); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
