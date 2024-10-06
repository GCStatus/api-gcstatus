package domain

import (
	"time"

	"gorm.io/gorm"
)

type Categoriable struct {
	gorm.Model
	ID               uint     `gorm:"primaryKey"`
	CategoriableID   uint     `gorm:"index"`
	CategoriableType string   `gorm:"index"`
	CategoryID       uint     `gorm:"index"`
	Category         Category `gorm:"foreignKey:CategoryID;references:ID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	Categoriable any `gorm:"-"`
}
