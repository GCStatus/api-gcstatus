package domain

import (
	"time"

	"gorm.io/gorm"
)

type Genreable struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	GenreableID   uint   `gorm:"index"`
	GenreableType string `gorm:"index"`
	GenreID       uint   `gorm:"index"`
	Genre         Genre  `gorm:"foreignKey:GenreID;references:ID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Genreable any `gorm:"-"`
}
