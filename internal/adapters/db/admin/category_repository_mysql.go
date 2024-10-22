package db_admin

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	ports_admin "gcstatus/internal/ports/admin"
	"net/http"

	"gorm.io/gorm"
)

type AdminCategoryRepositoryMySQL struct {
	db *gorm.DB
}

func NewAdminCategoryRepositoryMySQL(db *gorm.DB) ports_admin.AdminCategoryRepository {
	return &AdminCategoryRepositoryMySQL{
		db: db,
	}
}

func (h *AdminCategoryRepositoryMySQL) GetAll() ([]domain.Category, error) {
	var categories []domain.Category
	err := h.db.Model(&domain.Category{}).
		Find(&categories).
		Error

	return categories, err
}

func (h *AdminCategoryRepositoryMySQL) Create(category *domain.Category) error {
	return h.db.Create(&category).Error
}

func (h *AdminCategoryRepositoryMySQL) Update(id uint, request ports_admin.UpdateCategoryInterface) error {
	updateFields := map[string]any{
		"name": request.Name,
		"slug": request.Slug,
	}

	if err := h.db.Model(&domain.Category{}).Where("id = ?", id).Updates(updateFields).Error; err != nil {
		return fmt.Errorf("failed to update category: %+s", err.Error())
	}

	return nil
}

func (h *AdminCategoryRepositoryMySQL) Delete(id uint) error {
	if err := h.db.Delete(&domain.Category{}, id).Error; err != nil {
		return errors.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
