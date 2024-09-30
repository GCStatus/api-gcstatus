package domain

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Type      string `gorm:"not null" validate:"required"`
	Data      string `gorm:"not null" validate:"required"`
	ReadAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
}

type NotificationData struct {
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	ActionUrl string `json:"actionUrl"`
}

func (n *Notification) ValidateNotification() error {
	Init()

	err := validate.Struct(n)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
