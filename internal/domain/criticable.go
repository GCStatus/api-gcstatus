package domain

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Criticable struct {
	gorm.Model
	ID             uint            `gorm:"primaryKey"`
	Rate           decimal.Decimal `gorm:"not null;type:decimal(10,2)" validate:"required"`
	URL            string          `gorm:"size:255;not null" validate:"required"`
	PostedAt       time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CriticableID   uint   `gorm:"index"`
	CriticableType string `gorm:"index"`
	CriticID       uint   `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Critic         Critic `gorm:"foreignKey:CriticID"`
}

func (c *Criticable) ValidateCriticable() error {
	Init()

	if err := validate.Struct(c); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
