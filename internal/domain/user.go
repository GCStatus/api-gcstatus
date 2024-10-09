package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                uint      `gorm:"primaryKey"`
	Name              string    `gorm:"size:100;not null" validate:"required"`
	Email             string    `gorm:"unique;not null" validate:"required,email"`
	Nickname          string    `gorm:"unique;not null" validate:"required"`
	Experience        uint      `gorm:"not null; default:0"`
	Blocked           bool      `gorm:"not null; default:false"`
	Birthdate         time.Time `gorm:"not null" validate:"required"`
	Password          string    `gorm:"not null" validate:"required,min=8"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Profile           Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LevelID           uint    `gorm:"default:1"`
	Level             Level
	Wallet            Wallet
	Titles            []UserTitle             `json:"titles" gorm:"foreignKey:UserID"`
	TitleProgresses   []TitleProgress         `json:"title_progresses" gorm:"foreignKey:UserID"`
	MissionProgresses []MissionProgress       `json:"mission_progresses" gorm:"foreignKey:UserID"`
	Transactions      []Transaction           `json:"transactions" gorm:"foreignKey:UserID"`
	Notifications     []Notification          `json:"notifications" gorm:"foreignKey:UserID"`
	Missions          []UserMission           `gorm:"foreignKey:UserID"`
	MyMissions        []UserMissionAssignment `gorm:"foreignKey:UserID"`
	Reviews           []Reviewable            `gorm:"foreignKey:UserID"`
}

func (u *User) ValidateUser() error {
	Init()

	err := validate.Struct(u)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
