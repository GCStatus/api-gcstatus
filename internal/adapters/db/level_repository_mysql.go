package db

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type LevelRepositoryMySQL struct {
	db *gorm.DB
}

func NewLevelRepositoryMySQL(db *gorm.DB) ports.LevelRepository {
	return &LevelRepositoryMySQL{db: db}
}

func (h *LevelRepositoryMySQL) GetAll() ([]*domain.Level, error) {
	var levels []*domain.Level

	err := h.db.Preload("Rewards").Find(&levels).Error
	if err != nil {
		return levels, err
	}

	var titleIDs []uint
	for _, level := range levels {
		for _, reward := range level.Rewards {
			switch reward.RewardableType {
			case "titles":
				titleIDs = append(titleIDs, reward.RewardableID)
			}
		}
	}

	var titles []domain.Title
	if len(titleIDs) > 0 {
		err = h.db.Where("id IN (?)", titleIDs).Find(&titles).Error
		if err != nil {
			return levels, err
		}
	}

	titleMap := make(map[uint]domain.Title)
	for _, title := range titles {
		titleMap[title.ID] = title
	}

	for _, level := range levels {
		for i := range level.Rewards {
			reward := &level.Rewards[i]
			switch reward.RewardableType {
			case "titles":
				if title, ok := titleMap[reward.RewardableID]; ok {
					reward.Rewardable = &title
				}
			}
		}
	}

	return levels, nil
}

func (h *LevelRepositoryMySQL) FindById(id uint) (*domain.Level, error) {
	var level domain.Level
	err := h.db.Preload("Rewards").First(&level, id).Error
	return &level, err
}

func (h *LevelRepositoryMySQL) FindByLevel(lvl uint) (*domain.Level, error) {
	var level domain.Level
	err := h.db.Where("level = ?", lvl).First(&level).Error
	return &level, err
}
