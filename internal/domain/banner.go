package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	HomeHeaderCarouselBannersComponent = "home-header-carousel"
)

type Banner struct {
	gorm.Model
	ID             uint   `gorm:"primaryKey"`
	Component      string `gorm:"size:255;not null" validate:"required"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	BannerableID   uint   `gorm:"index"`
	BannerableType string `gorm:"index"`

	Bannerable any `gorm:"-"`
}

func (b *Banner) ValidateBanner() error {
	Init()

	err := validate.Struct(b)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
