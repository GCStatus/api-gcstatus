package domain

import (
	"time"

	"gorm.io/gorm"
)

type UserMission struct {
	gorm.Model
	ID              uint `gorm:"primaryKey"`
	Completed       bool `gorm:"not null"`
	LastCompletedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UserID          uint    `gorm:"not null"`
	MissionID       uint    `gorm:"not null"`
	User            User    `gorm:"foreignKey:UserID;references:ID"`
	Mission         Mission `gorm:"foreignKey:MissionID;references:ID"`
}

func (ut *UserMission) ValidateUserMission() error {
	Init()

	err := validate.Struct(ut)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
