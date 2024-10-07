package domain

import (
	"time"

	"gorm.io/gorm"
)

type Viewable struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	Count        uint   `gorm:"not null;default:0" validate:"required,numeric"`
	ViewableID   uint   `gorm:"index"`
	ViewableType string `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (v *Viewable) ValidateViewable() error {
	Init()

	err := validate.Struct(v)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
