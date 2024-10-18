package domain

import (
	"time"

	"gorm.io/gorm"
)

type Roleable struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	RoleableID   uint   `gorm:"index"`
	RoleableType string `gorm:"index"`
	RoleID       uint   `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Role         Role   `gorm:"foreignKey:RoleID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (r *Roleable) ValidateRoleable() error {
	Init()

	if err := validate.Struct(r); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
