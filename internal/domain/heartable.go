package domain

import (
	"time"

	"gorm.io/gorm"
)

type Heartable struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	HeartableID   uint   `gorm:"index"`
	HeartableType string `gorm:"index"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserID        uint `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	User          User `gorm:"foreignKey:UserID"`
}

func (h *Heartable) ValidateHeartable() error {
	Init()

	if err := validate.Struct(h); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
