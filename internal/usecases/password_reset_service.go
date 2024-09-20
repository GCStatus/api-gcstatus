package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
)

type PasswordResetService struct {
	repo ports.PasswordResetRepository
}

func NewPasswordResetService(repo ports.PasswordResetRepository) *PasswordResetService {
	return &PasswordResetService{repo: repo}
}

func (h *PasswordResetService) CreatePasswordReset(passwordReset *domain.PasswordReset) error {
	return h.repo.CreatePasswordReset(passwordReset)
}

func (h *PasswordResetService) FindPasswordResetByToken(token string) (*domain.PasswordReset, error) {
	passwordReset, err := h.repo.FindPasswordResetByToken(token)
	if err != nil {
		return nil, err
	}

	return passwordReset, nil
}

func (h *PasswordResetService) DeletePasswordReset(id uint) error {
	err := h.repo.DeletePasswordResetByID(id)
	if err != nil {
		return err
	}
	return nil
}
