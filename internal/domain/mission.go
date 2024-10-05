package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	MissionAvailable   = "available"
	MissionUnavailable = "unavailable"
	MissionCanceled    = "canceled"
	OneTimeMission     = "one-time"
	DailyMission       = "daily"
	WeeklyMission      = "weekly"
	MonthlyMission     = "monthly"
)

type Mission struct {
	gorm.Model
	ID                  uint      `gorm:"primaryKey"`
	Mission             string    `gorm:"not null" validate:"required"`
	Description         string    `gorm:"not null" validate:"required"`
	Status              string    `gorm:"not null;default:available" validate:"required"`
	ForAll              bool      `gorm:"not null;default:true" validate:"boolean"`
	Coins               uint      `gorm:"not null" validate:"required,numeric"`
	Experience          uint      `gorm:"not null" validate:"required,numeric"`
	Frequency           string    `gorm:"not null;default:one-time" validate:"required"`
	ResetTime           time.Time `gorm:"not null;autoCreateTime"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	MissionRequirements []MissionRequirement `gorm:"foreignKey:MissionID"`
	Rewards             []Reward             `gorm:"polymorphic:Sourceable;"`
	UserMission         []UserMission        `json:"user" gorm:"foreignKey:MissionID"`
}

func (t *Mission) ValidateMission() error {
	Init()

	err := validate.Struct(t)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
