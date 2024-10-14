package domain

import (
	"time"

	"gorm.io/gorm"
)

type DLCPublisher struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey"`
	DLCID       uint      `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	DLC         DLC       `gorm:"foreignKey:DLCID;references:ID"`
	PublisherID uint      `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Publisher   Publisher `gorm:"foreignKey:PublisherID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (dp *DLCPublisher) ValidateDLCPublisher() error {
	Init()

	err := validate.Struct(dp)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
