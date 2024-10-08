package domain

import (
	"time"

	"gorm.io/gorm"
)

type Critic struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null" validate:"required"`
	URL       string `gorm:"size:255;not null" validate:"required"`
	Logo      string `gorm:"size:255;not null" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Critic) ValidateCritic() error {
	Init()

	if err := validate.Struct(c); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
