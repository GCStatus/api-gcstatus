package db

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"gcstatus/pkg/utils"

	"gorm.io/gorm"
)

type ProfileRepositoryMySQL struct {
	db *gorm.DB
}

func NewProfileRepositoryMySQL(db *gorm.DB) ports.ProfileRepository {
	return &ProfileRepositoryMySQL{db: db}
}

func (h *ProfileRepositoryMySQL) UpdateSocials(profileID uint, request ports.UpdateSocialsRequest) error {
	updateFields := map[string]interface{}{
		"share":     request.Share,
		"phone":     utils.NullString(request.Phone),
		"github":    utils.NullString(request.Github),
		"twitch":    utils.NullString(request.Twitch),
		"twitter":   utils.NullString(request.Twitter),
		"youtube":   utils.NullString(request.Youtube),
		"facebook":  utils.NullString(request.Facebook),
		"instagram": utils.NullString(request.Instagram),
	}

	if err := h.db.Model(&domain.Profile{}).Where("id = ?", profileID).Updates(updateFields).Error; err != nil {
		return fmt.Errorf("failed to update socials: %w", err)
	}

	return nil
}

func (h *ProfileRepositoryMySQL) UpdatePicture(profileID uint, path string) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Profile{}).Where("id = ?", profileID).Update("photo", path).Error; err != nil {
			return err
		}

		return nil
	})
}
