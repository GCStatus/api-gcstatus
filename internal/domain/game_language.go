package domain

import (
	"time"

	"gorm.io/gorm"
)

type GameLanguage struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	Menu       bool `gorm:"not null;default:false" validate:"boolean"`
	Dubs       bool `gorm:"not null;default:false" validate:"boolean"`
	Subtitles  bool `gorm:"not null;default:false" validate:"boolean"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LanguageID uint     `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Language   Language `gorm:"foreignKey:LanguageID;references:ID"`
	GameID     uint     `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Game       Game     `gorm:"foreignKey:GameID;references:ID"`
}

func (gl *GameLanguage) ValidateGameLanguage() error {
	Init()

	err := validate.Struct(gl)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
