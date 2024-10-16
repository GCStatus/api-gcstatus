package ports

import "gcstatus/internal/domain"

type GameRepository interface {
	FindBySlug(slug string, userID uint) (domain.Game, error)
	FindGamesByCondition(condition string, limit *uint) ([]domain.Game, error)
	HomeGames() ([]domain.Game, []domain.Game, []domain.Game, *domain.Game, []domain.Game, error)
}
