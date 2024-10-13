package domain

import (
	"time"

	"gorm.io/gorm"
)

type DLCDeveloper struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey"`
	DLCID       uint      `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	DLC         DLC       `gorm:"foreignKey:DLCID;references:ID"`
	DeveloperID uint      `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Developer   Developer `gorm:"foreignKey:DeveloperID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (dd *DLCDeveloper) ValidateDLCDeveloper() error {
	Init()

	err := validate.Struct(dd)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
