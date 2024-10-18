package domain

import (
	"time"

	"gorm.io/gorm"
)

type Permissionable struct {
	gorm.Model
	ID                 uint       `gorm:"primaryKey"`
	PermissionableID   uint       `gorm:"index"`
	PermissionableType string     `gorm:"index"`
	PermissionID       uint       `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Permission         Permission `gorm:"foreignKey:PermissionID"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (p *Permissionable) ValidatePermissionable() error {
	Init()

	if err := validate.Struct(p); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
