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
	if err := h.db.Preload("Categories.Category").
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
		Preload("Galleries").
		Preload("DLCs.Galleries").
		Preload("DLCs.Platforms.Platform").
		Preload("DLCs.Stores.Store").
		Preload("Comments", "parent_id IS NULL").
		Preload("Comments.User").
		Preload("Comments.Replies.User").
		Preload("Support").
		Preload("Views").
		Preload("Hearts").
		Where("slug = ?", slug).
		First(&game).
		Error; err != nil {
		return game, err
	}

	var view domain.Viewable
	if err := h.db.Where("viewable_id = ? AND viewable_type = ? AND user_id = ?", game.ID, "games", userID).First(&view).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			view = domain.Viewable{
				ViewableID:   game.ID,
				ViewableType: "games",
				UserID:       userID,
			}
			if err := h.db.Create(&view).Error; err != nil {
				return game, err
			}
		} else {
			return game, err
		}
	}

	return game, nil
}
