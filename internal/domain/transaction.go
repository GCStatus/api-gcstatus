package domain

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID                uint   `gorm:"primaryKey"`
	Amount            uint   `gorm:"not null" validate:"required,numeric"`
	Description       string `gorm:"not null" validate:"required"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	UserID            uint
	User              User `gorm:"foreignKey:UserID"`
	TransactionTypeID uint
	TransactionType   TransactionType `gorm:"foreignKey:TransactionTypeID"`
}

func (t *Transaction) ValidateTransaction() error {
	Init()

	err := validate.Struct(t)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
