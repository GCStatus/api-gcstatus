package db

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type TitleRepositoryMySQL struct {
	db *gorm.DB
}

func NewTitleRepositoryMySQL(db *gorm.DB) ports.TitleRepository {
	return &TitleRepositoryMySQL{db: db}
}

func (h *TitleRepositoryMySQL) GetAllForUser(userID uint) ([]domain.Title, error) {
	var titles []domain.Title

	err := h.db.Model(&domain.Title{}).Where(
		"status NOT IN (?, ?)",
		domain.TitleUnavailable,
		domain.TitleCanceled,
	).
		Preload("TitleRequirements").
		Preload("TitleRequirements.TitleProgress", "user_id = ?", userID).
		Find(&titles).Error

	return titles, err
}

func (h *TitleRepositoryMySQL) FindById(titleID uint) (domain.Title, error) {
	var title domain.Title
	err := h.db.Model(&domain.Title{}).Where("id = ?", titleID).First(&title).Error
	return title, err
}

func (h *TitleRepositoryMySQL) ToggleEnableTitle(userID, titleID uint) error {
	var title domain.UserTitle

	if err := h.db.Model(&domain.UserTitle{}).Where("user_id = ? AND title_id = ?", userID, titleID).First(&title).Error; err != nil {
		return err
	}

	title.Enabled = !title.Enabled
	if err := h.db.Save(&title).Error; err != nil {
		return err
	}

	if err := h.db.Model(&domain.UserTitle{}).
		Where("user_id = ? AND title_id != ?", userID, titleID).
		Update("enabled", false).Error; err != nil {
		return err
	}

	return nil
}
