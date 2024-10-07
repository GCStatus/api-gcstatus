package db

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type GameRepositoryMySQL struct {
	db *gorm.DB
}

func NewGameRepositoryMySQL(db *gorm.DB) ports.GameRepository {
	return &GameRepositoryMySQL{db: db}
}

func (h *GameRepositoryMySQL) FindBySlug(slug string) (domain.Game, error) {
	var game domain.Game
	err := h.db.Preload("Categories.Category").
		Preload("Genres.Genre").
		Preload("Tags.Tag").
		Preload("Platforms.Platform").
		Preload("Languages.Language").
		Preload("Requirements.RequirementType").
		Preload("Crack.Cracker").
		Preload("Crack.Protection").
		Preload("Torrents.TorrentProvider").
		Preload("Publishers.Publisher").
		Preload("Developers.Developer").
		Preload("Support").
		Where("slug = ?", slug).
		First(&game).
		Error

	return game, err
}
