package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Nickname  string    `gorm:"unique;not null" json:"nickname"`
	Blocked   bool      `gorm:"not null" json:"blocked"`
	Birthdate time.Time `gorm:"not null" json:"birthdate"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Profile   Profile   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"profile"`
}
