package ports

import "gcstatus/internal/domain"

type GameRepository interface {
	FindBySlug(slug string) (domain.Game, error)
}
