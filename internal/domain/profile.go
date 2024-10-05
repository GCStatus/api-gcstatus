package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	ProfilePictureTitleRequirementKey = "update_profile_picture"
)

type Profile struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Share     bool `gorm:"not null" validate:"required,boolean"`
	Photo     string
	Phone     string
	Facebook  string
	Instagram string
	Twitter   string
	Youtube   string
	Twitch    string
	Github    string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint `gorm:"unique;constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
}

func (p *Profile) ValidateProfile() error {
	Init()

	err := validate.Struct(p)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
