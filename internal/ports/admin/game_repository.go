package ports_admin

import "gcstatus/internal/domain"

type AdminGameRepository interface {
	GetAll() ([]domain.Game, error)
	FindByID(id uint) (domain.Game, error)
}
