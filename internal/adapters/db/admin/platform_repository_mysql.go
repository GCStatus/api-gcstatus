package db_admin

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	ports_admin "gcstatus/internal/ports/admin"
	"net/http"

	"gorm.io/gorm"
)

type AdminPlatformRepositoryMySQL struct {
	db *gorm.DB
}

func NewAdminPlatformRepositoryMySQL(db *gorm.DB) ports_admin.AdminPlatformRepository {
	return &AdminPlatformRepositoryMySQL{
		db: db,
	}
}

func (h *AdminPlatformRepositoryMySQL) GetAll() ([]domain.Platform, error) {
	var platforms []domain.Platform
	err := h.db.Model(&domain.Platform{}).
		Find(&platforms).
		Error

	return platforms, err
}

func (h *AdminPlatformRepositoryMySQL) Create(platform *domain.Platform) error {
	return h.db.Create(&platform).Error
}

func (h *AdminPlatformRepositoryMySQL) Update(id uint, request ports_admin.UpdatePlatformInterface) error {
	updateFields := map[string]any{
		"name": request.Name,
		"slug": request.Slug,
	}

	if err := h.db.Model(&domain.Platform{}).Where("id = ?", id).Updates(updateFields).Error; err != nil {
		return fmt.Errorf("failed to update platform: %+s", err.Error())
	}

	return nil
}

func (h *AdminPlatformRepositoryMySQL) Delete(id uint) error {
	if err := h.db.Delete(&domain.Platform{}, id).Error; err != nil {
		return errors.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
