package domain

import (
	"time"

	"gorm.io/gorm"
)

type Requirement struct {
	gorm.Model
	ID                uint    `gorm:"primaryKey"`
	OS                string  `gorm:"size:255;not null" validate:"required"`
	DX                string  `gorm:"size:255;not null" validate:"required"`
	CPU               string  `gorm:"size:255;not null" validate:"required"`
	RAM               string  `gorm:"size:255;not null" validate:"required"`
	GPU               string  `gorm:"size:255;not null" validate:"required"`
	ROM               string  `gorm:"size:255;not null" validate:"required"`
	OBS               *string `gorm:"size:255"`
	Network           string  `gorm:"size:255;not null" validate:"required"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	RequirementTypeID uint            `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	RequirementType   RequirementType `gorm:"foreignKey:RequirementTypeID;references:ID"`
	GameID            uint            `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Game              Game            `gorm:"foreignKey:GameID;references:ID"`
}

func (r *Requirement) ValidateRequirement() error {
	Init()

	err := validate.Struct(r)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
