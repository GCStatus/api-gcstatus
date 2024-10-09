package ports

import "gcstatus/internal/domain"

type GameRepository interface {
	FindBySlug(slug string, userID uint) (domain.Game, error)
}
