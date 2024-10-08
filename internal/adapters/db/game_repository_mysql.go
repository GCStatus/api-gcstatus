package db

import (
	"errors"
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

func (h *GameRepositoryMySQL) FindBySlug(slug string, userID uint) (domain.Game, error) {
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
		Preload("Critics.Critic").
		Preload("Stores.Store").
		Preload("Support").
		Preload("Views").
		Preload("Hearts").
		Where("slug = ?", slug).
		First(&game).
		Error

	if err != nil {
		return game, err
	}

	var view domain.Viewable
	err = h.db.Where("viewable_id = ? AND viewable_type = ? AND user_id = ?", game.ID, "games", userID).First(&view).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		view = domain.Viewable{
			ViewableID:   game.ID,
			ViewableType: "games",
			UserID:       userID,
		}
		if err := h.db.Create(&view).Error; err != nil {
			return game, err
		}
	} else if err != nil {
		return game, err
	}

	return game, nil
}
