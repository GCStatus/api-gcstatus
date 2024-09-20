package ports

import "gcstatus/internal/domain"

type PasswordResetRepository interface {
	CreatePasswordReset(passwordReset *domain.PasswordReset) error
	FindPasswordResetByToken(token string) (*domain.PasswordReset, error)
	DeletePasswordResetByID(id uint) error
}
