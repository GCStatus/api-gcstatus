package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
)

type GameService struct {
	repo ports.GameRepository
}

func NewGameService(repo ports.GameRepository) *GameService {
	return &GameService{repo: repo}
}

func (h *GameService) FindBySlug(slug string) (domain.Game, error) {
	return h.repo.FindBySlug(slug)
}
