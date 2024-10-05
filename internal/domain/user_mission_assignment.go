package domain

import "gorm.io/gorm"

type UserMissionAssignment struct {
	gorm.Model
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"`
	MissionID uint    `gorm:"not null"`
	User      User    `gorm:"foreignKey:UserID;references:ID"`
	Mission   Mission `gorm:"foreignKey:MissionID;references:ID"`
}

func (um *UserMissionAssignment) ValidateUserMissionAssignment() error {
	Init()

	err := validate.Struct(um)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
