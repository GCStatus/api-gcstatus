package domain

import (
	"time"

	"gorm.io/gorm"
)

type DLCLanguage struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	Menu       bool `gorm:"not null;default:false" validate:"boolean"`
	Dubs       bool `gorm:"not null;default:false" validate:"boolean"`
	Subtitles  bool `gorm:"not null;default:false" validate:"boolean"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LanguageID uint     `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Language   Language `gorm:"foreignKey:LanguageID;references:ID"`
	DLCID      uint     `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	DLC        DLC      `gorm:"foreignKey:DLCID;references:ID"`
}

func (dl *DLCLanguage) ValidateDLCLanguage() error {
	Init()

	err := validate.Struct(dl)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
