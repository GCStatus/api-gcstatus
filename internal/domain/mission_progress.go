package domain

import (
	"time"

	"gorm.io/gorm"
)

type MissionProgress struct {
	gorm.Model
	ID                   uint `gorm:"primaryKey" json:"id"`
	Progress             int  `json:"progress"`
	Completed            bool `json:"completed"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	UserID               uint                `json:"user_id"`
	MissionRequirementID uint                `json:"Mission_requirement_id"`
	User                 User                `gorm:"foreignKey:UserID;references:ID"`
	MissionRequirement   *MissionRequirement `gorm:"foreignKey:MissionRequirementID;references:ID"`
}

func (mp *MissionProgress) ValidateMissionProgress() error {
	Init()

	err := validate.Struct(mp)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
