package domain

import (
	"time"

	"gorm.io/gorm"
)

type TorrentProvider struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	URL       string `gorm:"size:255;not null" validate:"required"`
	Name      string `gorm:"size:255;not null" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (tp *TorrentProvider) ValidateTorrentProvider() error {
	Init()

	err := validate.Struct(tp)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
