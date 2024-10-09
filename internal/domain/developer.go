package domain

import (
	"time"

	"gorm.io/gorm"
)

type Developer struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null" validate:"required"`
	Acting    bool   `gorm:"not null" validate:"boolean"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *Developer) ValidateDeveloper() error {
	Init()

	err := validate.Struct(d)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
