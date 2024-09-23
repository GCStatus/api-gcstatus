package domain

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Amount    int  `gorm:"not null; index; default:0" validate:"required,gte=0,numeric"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint `gorm:"unique;constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
}

func (w *Wallet) ValidateWallet() error {
	Init()

	err := validate.Struct(w)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
