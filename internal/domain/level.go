package domain

import (
	"time"

	"gorm.io/gorm"
)

type Level struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	Level      uint `gorm:"not null" validate:"required,numeric"`
	Experience uint `gorm:"not null" validate:"required,numeric"`
	Coins      uint `gorm:"not null" validate:"required,numeric"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Rewards    []Reward `gorm:"polymorphic:Sourceable;"`
}

func (l *Level) ValidateLevel() error {
	Init()

	err := validate.Struct(l)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
