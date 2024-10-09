package domain

import (
	"time"

	"gorm.io/gorm"
)

type Torrent struct {
	gorm.Model
	ID                uint      `gorm:"primaryKey"`
	URL               string    `gorm:"size:255;not null" validate:"required"`
	PostedAt          time.Time `gorm:"not null" validate:"required"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	TorrentProviderID uint            `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	TorrentProvider   TorrentProvider `gorm:"foreignKey:TorrentProviderID;references:ID"`
	GameID            uint            `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Game              Game            `gorm:"foreignKey:GameID;references:ID"`
}

func (t *Torrent) ValidateTorrent() error {
	Init()

	err := validate.Struct(t)
	if err != nil {
		return FormatValidationError(err)
	}

	return nil
}
