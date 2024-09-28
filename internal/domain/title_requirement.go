package domain

import (
	"time"

	"gorm.io/gorm"
)

type TitleRequirement struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey" json:"id"`
	Task          string `gorm:"not null" validate:"required"`
	Key           string `gorm:"not null" validate:"required"`
	Goal          int    `gorm:"not null" validate:"required,numeric"`
	Description   string `gorm:"not null" validate:"required"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TitleProgress TitleProgress `json:"progress" gorm:"foreignKey:TitleRequirementID"`
	TitleID       uint
	Title         Title `gorm:"foreignKey:TitleID;references:ID"`
}

func (tr *TitleRequirement) ValidateTitleRequirement() error {
	Init()

	err := validate.Struct(tr)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
