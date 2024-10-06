package domain

import (
	"time"

	"gorm.io/gorm"
)

type Cracker struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null" validate:"required"`
	Acting    bool   `gorm:"not null" validate:"boolean"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Cracker) ValidateCracker() error {
	Init()

	err := validate.Struct(c)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
