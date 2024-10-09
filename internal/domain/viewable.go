package domain

import (
	"time"

	"gorm.io/gorm"
)

type Viewable struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	ViewableID   uint   `gorm:"index"`
	ViewableType string `gorm:"index"`
	UserID       uint   `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	User         User   `gorm:"foreignKey:UserID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (v *Viewable) ValidateViewable() error {
	Init()

	err := validate.Struct(v)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
