package domain

import (
	"time"

	"gorm.io/gorm"
)

type Reviewable struct {
	gorm.Model
	ID             uint   `gorm:"primaryKey"`
	Rate           uint   `gorm:"not nul;" validate:"required,numeric"`
	Review         string `gorm:"size:255;not null" validate:"required"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ReviewableID   uint   `gorm:"index"`
	ReviewableType string `gorm:"index"`
	UserID         uint   `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	User           User   `gorm:"foreignKey:UserID"`
}
