package domain

import (
	"time"

	"gorm.io/gorm"
)

type Commentable struct {
	gorm.Model
	ID              uint   `gorm:"primaryKey"`
	Comment         string `gorm:"size:255;type:text;not null" validate:"required"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UserID          uint          `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	User            User          `gorm:"foreignKey:UserID"`
	CommentableID   uint          `gorm:"index"`
	CommentableType string        `gorm:"index"`
	ParentID        *uint         `gorm:"index"`
	Replies         []Commentable `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
	Hearts          []Heartable   `gorm:"polymorphic:Heartable"`
}

func (c *Commentable) ValidateCommentable() error {
	Init()

	if err := validate.Struct(c); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
