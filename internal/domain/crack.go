package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	CrackedStatus   = "cracked"
	UncrackedStatus = "uncracked"
	CrackedSameDay  = "cracked-oneday"
)

type Crack struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	Status       string `gorm:"size:255;not null;type:enum('cracked','uncracked','cracked-oneday');default:uncracked" validate:"required"`
	CrackedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CrackerID    uint       `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Cracker      Cracker    `gorm:"foreignKey:CrackerID;references:ID;"`
	ProtectionID uint       `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Protection   Protection `gorm:"foreignKey:ProtectionID;references:ID;"`
	GameID       uint       `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Game         Game       `gorm:"foreignKey:GameID;references:ID;"`
}

func (c *Crack) ValidateCrack() error {
	Init()

	err := validate.Struct(c)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
