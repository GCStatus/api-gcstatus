package db

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type PasswordResetRepositoryMySQL struct {
	db *gorm.DB
}

func NewPasswordResetRepositoryMySQL(db *gorm.DB) ports.PasswordResetRepository {
	return &PasswordResetRepositoryMySQL{db: db}
}

func (h *PasswordResetRepositoryMySQL) CreatePasswordReset(passwordReset *domain.PasswordReset) error {
	return h.db.Create(passwordReset).Error
}

func (h *PasswordResetRepositoryMySQL) FindPasswordResetByToken(token string) (*domain.PasswordReset, error) {
	var passwordReset domain.PasswordReset
	err := h.db.Where("token = ?", token).First(&passwordReset).Error
	return &passwordReset, err
}

func (repo *PasswordResetRepositoryMySQL) DeletePasswordResetByID(id uint) error {
	return repo.db.Delete(&domain.PasswordReset{}, id).Error
}
