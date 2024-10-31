package db

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type HeartRepositryMySQL struct {
	db *gorm.DB
}

func NewHeartRepositoryMySQL(db *gorm.DB) ports.HeartRepositry {
	return &HeartRepositryMySQL{db: db}
}

func (h *HeartRepositryMySQL) Create(heartable *domain.Heartable) error {
	if err := h.db.Create(&heartable).Error; err != nil {
		return err
	}

	return nil
}

func (h *HeartRepositryMySQL) Delete(id uint) error {
	if err := h.db.Unscoped().Delete(&domain.Heartable{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (h *HeartRepositryMySQL) FindForUser(heartableID uint, heartableType string, userID uint) (*domain.Heartable, error) {
	var heart *domain.Heartable

	if err := h.db.Model(&domain.Heartable{}).
		Where("heartable_id = ? AND heartable_type = ? AND user_id = ?", heartableID, heartableType, userID).
		First(&heart).
		Error; err != nil {
		return nil, err
	}

	return heart, nil
}
