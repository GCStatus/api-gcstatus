package domain

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Scope     string `gorm:"size:255;not null;unique" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Permission) ValidatePermission() error {
	Init()

	if err := validate.Struct(p); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
