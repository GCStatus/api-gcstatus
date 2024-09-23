package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
)

type LevelService struct {
	repo ports.LevelRepository
}

func NewLevelService(repo ports.LevelRepository) *LevelService {
	return &LevelService{repo: repo}
}

func (h *LevelService) GetAll() ([]*domain.Level, error) {
	return h.repo.GetAll()
}

func (h *LevelService) FindById(id uint) (*domain.Level, error) {
	return h.repo.FindById(id)
}

func (h *LevelService) FindByLevel(level uint) (*domain.Level, error) {
	return h.repo.FindByLevel(level)
}
