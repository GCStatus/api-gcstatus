package db_admin

import (
	"gcstatus/internal/domain"

	"gorm.io/gorm"
)

type AdminGameRepositoryMySQL struct {
	db *gorm.DB
}

func NewAdminGameRepositoryMySQL(db *gorm.DB) *AdminGameRepositoryMySQL {
	return &AdminGameRepositoryMySQL{db: db}
}

func (h *AdminGameRepositoryMySQL) GetAll() ([]domain.Game, error) {
	var games []domain.Game
	err := h.db.Model(&domain.Game{}).
		Preload("Platforms.Platform").
		Preload("Categories.Category").
		Preload("Genres.Genre").
		Preload("Tags.Tag").
		Preload("Hearts").
		Preload("Views").
		Preload("Crack.Cracker").
		Preload("Crack.Protection").
		Order("created_at DESC").
		Find(&games).
		Error

	return games, err
}

func (h *AdminGameRepositoryMySQL) FindByID(id uint) (domain.Game, error) {
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
		Preload("Galleries.MediaType").
		Preload("DLCs.Galleries.MediaType").
		Preload("DLCs.Platforms.Platform").
		Preload("DLCs.Stores.Store").
		Preload("Comments", "parent_id IS NULL").
		Preload("Comments.Hearts").
		Preload("Comments.User").
		Preload("Comments.Replies.User").
		Preload("Support").
		Preload("Views").
		Preload("Hearts").
		Where("id = ?", id).
		First(&game).
		Error; err != nil {
		return game, err
	}

	return game, nil
}
