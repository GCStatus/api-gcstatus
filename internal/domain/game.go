package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	HotCondition     = "hot"
	SaleCondition    = "sale"
	CommomCondition  = "commom"
	PopularCondition = "popular"
)

type Game struct {
	gorm.Model
	ID               uint      `gorm:"primaryKey"`
	Age              int       `gorm:"not null" validate:"required,numeric"`
	Slug             string    `gorm:"size:255;uniqueIndex;not null" validate:"required"`
	Title            string    `gorm:"size:255;not null" validate:"required"`
	Condition        string    `gorm:"size:255;not null;type:enum('hot','sale','popular','commom');default:commom" validate:"required"`
	Cover            string    `gorm:"size:255" validate:"required"`
	About            string    `gorm:"type:text" validate:"required"`
	Description      string    `gorm:"type:text" validate:"required"`
	ShortDescription string    `gorm:"size:255" validate:"required"`
	Free             bool      `gorm:"not null;default:false" validate:"boolean"`
	Legal            *string   `gorm:"size:255"`
	Website          *string   `gorm:"size:255"`
	ReleaseDate      time.Time `gorm:"size:255" validate:"required"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Categories       []Categoriable `gorm:"polymorphic:Categoriable;"`
	Tags             []Taggable     `gorm:"polymorphic:Taggable;"`
	Genres           []Genreable    `gorm:"polymorphic:Genreable;"`
	Platforms        []Platformable `gorm:"polymorphic:Platformable;"`
}

func (g *Game) ValidateGame() error {
	Init()

	err := validate.Struct(g)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}