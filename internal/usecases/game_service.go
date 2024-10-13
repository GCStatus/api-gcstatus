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

func (h *GameService) FindGamesByCondition(condition string, limit *uint) ([]domain.Game, error) {
	return h.repo.FindGamesByCondition(condition, limit)
}

func (h *GameService) HomeGames() ([]domain.Game, []domain.Game, []domain.Game, *domain.Game, []domain.Game, error) {
	return h.repo.HomeGames()
}

func (h *GameService) FindBySlug(slug string, userID uint) (domain.Game, error) {
	return h.repo.FindBySlug(slug, userID)
}
