package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	AdditionTransactionType      = "addition"
	SubtractionTransactionType   = "subtraction"
	AdditionTransactionTypeID    = 1
	SubtractionTransactionTypeID = 2
)

type TransactionType struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Type      string `gorm:"not null" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (tt *TransactionType) ValidateTransactionType() error {
	Init()

	err := validate.Struct(tt)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
