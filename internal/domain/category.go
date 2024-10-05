package domain

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null;unique" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Category) ValidateCategory() error {
	Init()

	err := validate.Struct(c)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
