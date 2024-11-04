package db

import (
	"errors"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"time"

	"gorm.io/gorm"
)

type GameRepositoryMySQL struct {
	db *gorm.DB
}

func NewGameRepositoryMySQL(db *gorm.DB) ports.GameRepository {
	return &GameRepositoryMySQL{db: db}
}

func (h *GameRepositoryMySQL) HomeGames() (
	[]domain.Game,
	[]domain.Game,
	[]domain.Game,
	*domain.Game,
	[]domain.Game,
	error,
) {
	var hotGames, popularGames []domain.Game

	hotQuery := h.db.Model(&domain.Game{}).
		Preload("Platforms.Platform").
		Preload("Categories.Category").
		Preload("Genres.Genre").
		Preload("Tags.Tag").
		Preload("Hearts").
		Preload("Views").
		Preload("Crack.Cracker").
		Preload("Crack.Protection").
		Where("`condition` = ?", "hot").
		Order("created_at DESC").
		Limit(9)

	if err := hotQuery.Find(&hotGames).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	popularQuery := h.db.Model(&domain.Game{}).
		Preload("Platforms.Platform").
		Preload("Categories.Category").
		Preload("Genres.Genre").
		Preload("Tags.Tag").
		Preload("Hearts").
		Preload("Views").
		Preload("Crack.Cracker").
		Preload("Crack.Protection").
		Where("`condition` = ?", "popular").
		Order("created_at DESC").
		Limit(9)

	if err := popularQuery.Find(&popularGames).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	var mostHeartedGames []domain.Game
	subQuery := h.db.Table("heartables").
		Select("heartable_id, COUNT(*) AS heart_count").
		Where("heartable_type = ?", "games").
		Group("heartable_id").
		Order("heart_count DESC").
		Limit(9)

	if err := h.db.Model(&domain.Game{}).
		Joins("JOIN (?) AS h ON games.id = h.heartable_id", subQuery).
		Preload("Platforms.Platform").
		Preload("Categories.Category").
		Preload("Genres.Genre").
		Preload("Tags.Tag").
		Preload("Hearts").
		Preload("Views").
		Preload("Crack.Cracker").
		Preload("Crack.Protection").
		Find(&mostHeartedGames).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	var nextGreatReleaseGame *domain.Game
	err := h.db.Model(&domain.Game{}).
		Preload("Platforms.Platform").
		Preload("Genres.Genre").
		Where("`great_release` = ?", true).
		Where("`release_date` >= ?", time.Now()).
		First(&nextGreatReleaseGame).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil, nil, nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		nextGreatReleaseGame = nil
	}

	var upcomingGames []domain.Game
	if err := h.db.Model(&domain.Game{}).
		Preload("Platforms.Platform").
		Where("`release_date` > ?", time.Now()).
		Find(&upcomingGames).
		Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return hotGames, popularGames, mostHeartedGames, nextGreatReleaseGame, upcomingGames, nil
}

func (h *GameRepositoryMySQL) FindGamesByCondition(condition string, limit *uint) ([]domain.Game, error) {
	var games []domain.Game

	query := h.db.Model(&domain.Game{}).
		Preload("Platforms.Platform").
		Preload("Categories.Category").
		Preload("Genres.Genre").
		Preload("Tags.Tag").
		Preload("Hearts").
		Preload("Views").
		Preload("Crack.Cracker").
		Preload("Crack.Protection").
		Where("`condition` = ?", condition).
		Order("created_at DESC")

	if limit != nil {
		query = query.Limit(int(*limit))
	}

	if err := query.Find(&games).Error; err != nil {
		return games, err
	}

	return games, nil
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
		Preload("Galleries.MediaType").
		Preload("DLCs.Galleries.MediaType").
		Preload("DLCs.Platforms.Platform").
		Preload("DLCs.Stores.Store").
		Preload("Comments", "parent_id IS NULL").
		Preload("Comments.Hearts").
		Preload("Comments.User").
		Preload("Comments.Replies.User").
		Preload("Comments.Replies.Hearts").
		Preload("Support").
		Preload("Views").
		Preload("Hearts").
		Where("slug = ?", slug).
		First(&game).
		Error; err != nil {
		return game, err
	}

	if userID != 0 {
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
	}

	return game, nil
}

func (h *GameRepositoryMySQL) ExistsForStore(storeID uint, appID uint) (bool, error) {
	var count int64
	err := h.db.Model(&domain.GameStore{}).
		Where("store_id = ? AND store_game_id = ?", storeID, appID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (h *GameRepositoryMySQL) Search(input string) ([]domain.Game, error) {
	var games []domain.Game

	likeInput := "%" + input + "%"

	query := h.db.Model(&domain.Game{}).
		Preload("Platforms.Platform").
		Preload("Categories.Category").
		Preload("Genres.Genre").
		Preload("Tags.Tag").
		Preload("Hearts").
		Preload("Views").
		Preload("Crack.Cracker").
		Preload("Crack.Protection")

	columns := []string{"title", "description", "about", "short_description"}

	for _, column := range columns {
		query = query.Or(column+" LIKE ?", likeInput)
	}

	if err := query.Limit(100).Find(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (h *GameRepositoryMySQL) FindByClassification(classification string, filterable string) ([]domain.Game, error) {
	var games []domain.Game
	query := h.db.Model(&domain.Game{}).
		Preload("Hearts").
		Preload("Views").
		Preload("Crack.Cracker").
		Preload("Crack.Protection").
		Preload("Platforms.Platform").
		Preload("Categories.Category").
		Preload("Tags.Tag").
		Preload("Genres.Genre")

	switch classification {
	case "categories":
		query = query.Joins("JOIN categoriables ON categoriables.categoriable_id = games.id AND categoriables.categoriable_type = 'games'").
			Joins("JOIN categories ON categories.id = categoriables.category_id").
			Where("categories.slug = ?", filterable)
	case "platforms":
		query = query.Joins("JOIN platformables ON platformables.platformable_id = games.id AND platformables.platformable_type = 'games'").
			Joins("JOIN platforms ON platforms.id = platformables.platform_id").
			Where("platforms.slug = ?", filterable)
	case "tags":
		query = query.Joins("JOIN taggables ON taggables.taggable_id = games.id AND taggables.taggable_type = 'games'").
			Joins("JOIN tags ON tags.id = taggables.tag_id").
			Where("tags.slug = ?", filterable)
	case "genres":
		query = query.Joins("JOIN genreables ON genreables.genreable_id = games.id AND genreables.genreable_type = 'games'").
			Joins("JOIN genres ON genres.id = genreables.genre_id").
			Where("genres.slug = ?", filterable)
	case "crackers":
		query = query.Joins("JOIN cracks ON cracks.game_id = games.id").
			Joins("JOIN crackers ON crackers.id = cracks.cracker_id").
			Where("crackers.slug = ?", filterable)
	case "publishers":
		query = query.Joins("JOIN game_publishers ON game_publishers.game_id = games.id").
			Joins("JOIN publishers ON publishers.id = game_publishers.publisher_id").
			Where("publishers.slug = ?", filterable)
	case "developers":
		query = query.Joins("JOIN game_developers ON game_developers.game_id = games.id").
			Joins("JOIN developers ON developers.id = game_developers.developer_id").
			Where("developers.slug = ?", filterable)
	case "protections":
		query = query.Joins("JOIN cracks ON cracks.game_id = games.id").
			Joins("JOIN protections ON protections.id = cracks.protection_id").
			Where("protections.slug = ?", filterable)
	case "cracks":
		query = query.Joins("JOIN cracks ON cracks.game_id = games.id").
			Where("cracks.status = ?", filterable)
	default:
		return []domain.Game{}, nil
	}

	err := query.Limit(100).Find(&games).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return games, err
	}

	return games, nil
}

func (h *GameRepositoryMySQL) CalendarGames() ([]domain.Game, error) {
	var games []domain.Game

	now := time.Now()

	oneMonthAgo := now.AddDate(0, -1, 0)

	err := h.db.Model(&domain.Game{}).
		Preload("Crack").
		Where("release_date >= ?", oneMonthAgo).
		Limit(100).
		Find(&games).Error

	return games, err
}
