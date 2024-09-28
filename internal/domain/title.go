package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	TitleAvailable   = "available"
	TitleUnavailable = "unavailable"
	TitleCanceled    = "canceled"
)

type Title struct {
	gorm.Model
	ID                uint   `gorm:"primaryKey"`
	Title             string `gorm:"not null" validate:"required"`
	Description       string `gorm:"not null" validate:"required"`
	Cost              *int
	Purchasable       bool   `gorm:"not null; default:0" validate:"boolean"`
	Status            string `gorm:"not null; default:available" validate:"required"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	TitleRequirements []TitleRequirement `gorm:"foreignKey:TitleID"`
	Rewards           []Reward           `gorm:"polymorphic:Rewardable;"`
	Users             []UserTitle        `json:"users" gorm:"foreignKey:TitleID"`
}

func (t *Title) ValidateTitle() error {
	Init()

	err := validate.Struct(t)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
