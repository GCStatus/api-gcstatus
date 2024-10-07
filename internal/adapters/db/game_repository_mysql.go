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
		Preload("Reviews.User.Profile").
		Preload("Support").
		Preload("View").
		Where("slug = ?", slug).
		First(&game).
		Error

	if err != nil {
		return game, err
	}

	if game.View.ID == 0 {
		view := domain.Viewable{
			Count:        1,
			ViewableID:   game.ID,
			ViewableType: "games",
		}
		if err := h.db.Create(&view).Error; err != nil {
			return game, err
		}
		game.View = view
	} else {
		game.View.Count += 1
		if err := h.db.Save(&game.View).Error; err != nil {
			return game, err
		}
	}

	return game, nil
}
