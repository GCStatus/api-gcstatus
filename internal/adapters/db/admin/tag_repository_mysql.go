package db_admin

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	ports_admin "gcstatus/internal/ports/admin"
	"net/http"

	"gorm.io/gorm"
)

type AdminTagRepositoryMySQL struct {
	db *gorm.DB
}

func NewAdminTagRepositoryMySQL(db *gorm.DB) ports_admin.AdminTagRepository {
	return &AdminTagRepositoryMySQL{
		db: db,
	}
}

func (h *AdminTagRepositoryMySQL) GetAll() ([]domain.Tag, error) {
	var tags []domain.Tag
	err := h.db.Model(&domain.Tag{}).
		Find(&tags).
		Error

	return tags, err
}

func (h *AdminTagRepositoryMySQL) Create(tag *domain.Tag) error {
	return h.db.Create(&tag).Error
}

func (h *AdminTagRepositoryMySQL) Update(id uint, request ports_admin.UpdateTagInterface) error {
	updateFields := map[string]any{
		"name": request.Name,
		"slug": request.Slug,
	}

	if err := h.db.Model(&domain.Tag{}).Where("id = ?", id).Updates(updateFields).Error; err != nil {
		return fmt.Errorf("failed to update tag: %+s", err.Error())
	}

	return nil
}

func (h *AdminTagRepositoryMySQL) Delete(id uint) error {
	if err := h.db.Delete(&domain.Tag{}, id).Error; err != nil {
		return errors.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
