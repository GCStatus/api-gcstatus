package domain

import (
	"time"

	"gorm.io/gorm"
)

type Platformable struct {
	gorm.Model
	ID               uint     `gorm:"primaryKey"`
	PlatformableID   uint     `gorm:"index"`
	PlatformableType string   `gorm:"index"`
	PlatformID       uint     `gorm:"index"`
	Platform         Platform `gorm:"foreignKey:PlatformID;references:ID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	Platformable any `gorm:"-"`
}
