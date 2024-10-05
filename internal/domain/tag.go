package domain

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null;unique" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Tag) ValidateTag() error {
	Init()

	err := validate.Struct(t)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
