package domain

import (
	"time"

	"gorm.io/gorm"
)

type PasswordReset struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"not null; index" validate:"required,email"`
	Token     string    `gorm:"not null; unique" validate:"required"`
	ExpiresAt time.Time `gorm:"not null; index" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *PasswordReset) ValidatePasswordReset() error {
	Init()

	err := validate.Struct(p)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
