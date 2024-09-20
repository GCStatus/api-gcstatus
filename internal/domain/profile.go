package domain

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	Share     bool      `gorm:"not null" json:"share"`
	Photo     string    `json:"photo"`
	Phone     string    `json:"phone"`
	Facebook  string    `json:"facebook"`
	Instagram string    `json:"instagram"`
	Twitter   string    `json:"twitter"`
	Youtube   string    `json:"youtube"`
	Twitch    string    `json:"twitch"`
	Github    string    `json:"github"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint
}
