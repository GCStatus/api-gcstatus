package domain

import (
	"time"

	"gorm.io/gorm"
)

type PasswordReset struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"not null; index"`
	Token     string    `gorm:"not null; unique"`
	ExpiresAt time.Time `gorm:"not null; index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
