package ports

import "gcstatus/internal/domain"

type GameRepository interface {
	FindBySlug(slug string, userID uint) (domain.Game, error)
	FindGamesByCondition(condition string, limit *uint) ([]domain.Game, error)
	FindByClassification(classification string, filterable string) ([]domain.Game, error)
	HomeGames() ([]domain.Game, []domain.Game, []domain.Game, *domain.Game, []domain.Game, error)
	ExistsForStore(storeID uint, appID uint) (bool, error)
	Search(input string) ([]domain.Game, error)
	CalendarGames() ([]domain.Game, error)
}
