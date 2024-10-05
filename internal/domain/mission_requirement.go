package domain

import (
	"time"

	"gorm.io/gorm"
)

type MissionRequirement struct {
	gorm.Model
	ID              uint   `gorm:"primaryKey" json:"id"`
	Task            string `gorm:"not null" validate:"required"`
	Key             string `gorm:"not null" validate:"required"`
	Goal            int    `gorm:"not null" validate:"required,numeric"`
	Description     string `gorm:"not null" validate:"required"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	MissionProgress MissionProgress `json:"progress" gorm:"foreignKey:MissionRequirementID"`
	MissionID       uint
	Mission         Mission `gorm:"foreignKey:MissionID;references:ID"`
}

func (tr *MissionRequirement) ValidateMissionRequirement() error {
	Init()

	err := validate.Struct(tr)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
