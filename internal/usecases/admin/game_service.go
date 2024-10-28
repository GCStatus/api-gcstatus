package usecases_admin

import (
	"gcstatus/internal/domain"
	ports_admin "gcstatus/internal/ports/admin"
)

type AdminGameService struct {
	repo ports_admin.AdminGameRepository
}

func NewAdminGameService(repo ports_admin.AdminGameRepository) *AdminGameService {
	return &AdminGameService{repo: repo}
}

func (h *AdminGameService) GetAll() ([]domain.Game, error) {
	return h.repo.GetAll()
}

func (h *AdminGameService) FindByID(id uint) (domain.Game, error) {
	return h.repo.FindByID(id)
}
