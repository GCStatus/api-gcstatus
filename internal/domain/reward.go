package domain

import (
	"time"

	"gorm.io/gorm"
)

type Reward struct {
	gorm.Model
	ID             uint   `gorm:"primaryKey"`
	SourceableID   uint   `gorm:"not null" validate:"required"`
	SourceableType string `gorm:"not null" validate:"required"`

	RewardableID   uint   `gorm:"not null" validate:"required"`
	RewardableType string `gorm:"not null" validate:"required"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Rewardable any `gorm:"-"`
	Sourceable any `gorm:"-"`
}

func (r *Reward) ValidateReward() error {
	Init()

	err := validate.Struct(r)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
