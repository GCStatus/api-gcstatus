package domain

import (
	"time"

	"gorm.io/gorm"
)

type UserTitle struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Enabled   bool `gorm:"not null; default:false" validate:"boolean"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint  `gorm:"not null"`
	TitleID   uint  `gorm:"not null"`
	User      User  `gorm:"foreignKey:UserID;references:ID"`
	Title     Title `gorm:"foreignKey:TitleID;references:ID"`
}

func (ut *UserTitle) ValidateUserTitle() error {
	Init()

	err := validate.Struct(ut)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
