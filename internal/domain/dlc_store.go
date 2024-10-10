package domain

import (
	"time"

	"gorm.io/gorm"
)

type DLCStore struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Price     uint   `gorm:"not null" validate:"required"`
	URL       string `gorm:"size:255;not null" validate:"required"`
	DLCID     uint   `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	DLC       DLC    `gorm:"foreignKey:DLCID;references:ID"`
	StoreID   uint   `gorm:"constraint:OnDelete:SET NULL,OnUpdate:CASCADE;"`
	Store     Store  `gorm:"foreignKey:StoreID;references:ID"`
	StorDLCID string `gorm:"not null;" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ds *DLCStore) ValidateDLCStore() error {
	Init()

	if err := validate.Struct(ds); err != nil {
		return FormatValidationError(err)
	}

	return nil
}
