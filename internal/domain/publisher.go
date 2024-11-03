package domain

import (
	"time"

	"gorm.io/gorm"
)

type Publisher struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null" validate:"required"`
	Slug      string `gorm:"size:255;not null" validate:"required"`
	Acting    bool   `gorm:"not null" validate:"boolean"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Publisher) ValidatePublisher() error {
	Init()

	err := validate.Struct(p)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
