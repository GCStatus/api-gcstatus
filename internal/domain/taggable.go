package domain

import (
	"time"

	"gorm.io/gorm"
)

type Taggable struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	TaggableID   uint   `gorm:"index"`
	TaggableType string `gorm:"index"`
	TagID        uint   `gorm:"index"`
	Tag          Tag    `gorm:"foreignKey:TagID;references:ID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Taggable any `gorm:"-"`
}
