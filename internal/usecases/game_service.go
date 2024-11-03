package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	"gcstatus/internal/ports"
	"net/http"
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

func (h *GameService) ExistsForStore(storeID uint, appID uint) (bool, error) {
	return h.repo.ExistsForStore(storeID, appID)
}

func (h *GameService) Search(input string) ([]domain.Game, error) {
	return h.repo.Search(input)
}

func (h *GameService) FindByClassification(classification string, filterable string) ([]domain.Game, error) {
	validClassifications := map[string]bool{
		"categories":  true,
		"genres":      true,
		"tags":        true,
		"platforms":   true,
		"crackers":    true,
		"developers":  true,
		"publishers":  true,
		"protections": true,
		"cracks":      true,
	}

	if !validClassifications[classification] {
		return nil, errors.NewHttpError(
			http.StatusBadRequest,
			"The given classification is not valid. The valid classifications are: crackers, cracks, publishers, protections, developers, platforms, genres, tags and categories.",
		)
	}

	return h.repo.FindByClassification(classification, filterable)
}
