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
	err := h.db.Find(&levels).Error
	return levels, err
}

func (h *LevelRepositoryMySQL) FindById(id uint) (*domain.Level, error) {
	var level domain.Level
	err := h.db.First(&level, id).Error
	return &level, err
}

func (h *LevelRepositoryMySQL) FindByLevel(lvl uint) (*domain.Level, error) {
	var level domain.Level
	err := h.db.Where("level = ?", lvl).First(&level).Error
	return &level, err
}
