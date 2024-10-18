package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:255;not null" validate:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Permissions []Permissionable `gorm:"polymorphic:Permissionable"`
}

func (r *Role) ValidateRole() error {
	Init()

	if err := validate.Struct(r); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
