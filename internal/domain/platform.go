package domain

import (
	"time"

	"gorm.io/gorm"
)

type Platform struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null;unique" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Platform) ValidatePlatform() error {
	Init()

	err := validate.Struct(p)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
