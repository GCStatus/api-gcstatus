package domain

import (
	"time"

	"gorm.io/gorm"
)

type Protection struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Protection) ValidateProtection() error {
	Init()

	err := validate.Struct(p)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
