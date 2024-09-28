package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
)

type TitleService struct {
	repo ports.TitleRepository
}

func NewTitleService(repo ports.TitleRepository) *TitleService {
	return &TitleService{repo: repo}
}

func (h *TitleService) GetAll(userID uint) ([]domain.Title, error) {
	return h.repo.GetAll(userID)
}

func (h *TitleService) ToggleEnableTitle(userID, titleID uint) error {
	return h.repo.ToggleEnableTitle(userID, titleID)
}

func (h *TitleService) FindById(titleID uint) (domain.Title, error) {
	return h.repo.FindById(titleID)
}
