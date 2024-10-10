package domain

import (
	"time"

	"gorm.io/gorm"
)

type Galleriable struct {
	gorm.Model
	ID              uint      `gorm:"primaryKey"`
	S3              bool      `gorm:"not null;default:false" validate:"boolean"`
	Path            string    `gorm:"size:255;not null" validate:"required"`
	GalleriableID   uint      `gorm:"index"`
	GalleriableType string    `gorm:"index"`
	MediaTypeID     uint      `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	MediaType       MediaType `gorm:"foreignKey:MediaTypeID;references:ID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (g *Galleriable) ValidateGalleriable() error {
	Init()

	if err := validate.Struct(g); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
