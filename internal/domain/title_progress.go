package domain

import (
	"time"

	"gorm.io/gorm"
)

type TitleProgress struct {
	gorm.Model
	ID                 uint `gorm:"primaryKey" json:"id"`
	Progress           int  `json:"progress"`
	Completed          bool `json:"completed"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	UserID             uint              `json:"user_id"`
	TitleRequirementID uint              `json:"title_requirement_id"`
	User               User              `gorm:"foreignKey:UserID;references:ID"`
	TitleRequirement   *TitleRequirement `gorm:"foreignKey:TitleRequirementID;references:ID"`
}

func (tp *TitleProgress) ValidateTitleProgress() error {
	Init()

	err := validate.Struct(tp)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
