package db

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type BannerRepositoryMySQL struct {
	db *gorm.DB
}

func NewBannerRepositoryMySQL(db *gorm.DB) ports.BannerRepository {
	return &BannerRepositoryMySQL{db: db}
}

func (h *BannerRepositoryMySQL) GetBannersForHome() ([]domain.Banner, error) {
	var banners []domain.Banner

	if err := h.db.Model(&domain.Banner{}).
		Where("component = ?", domain.HomeHeaderCarouselBannersComponent).
		Find(&banners).Error; err != nil {
		return banners, err
	}

	for i, banner := range banners {
		if banner.BannerableType == "games" {
			var game domain.Game
			if err := h.db.Where("id = ?", banner.BannerableID).
				Preload("Genres.Genre").
				Preload("Platforms.Platform").
				Preload("Crack").
				First(&game).Error; err == nil {
				banners[i].Bannerable = game
			}
		}
	}

	return banners, nil
}
