package db_admin

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	ports_admin "gcstatus/internal/ports/admin"
	"net/http"

	"gorm.io/gorm"
)

type AdminGenreRepositoryMySQL struct {
	db *gorm.DB
}

func NewAdminGenreRepositoryMySQL(db *gorm.DB) ports_admin.AdminGenreRepository {
	return &AdminGenreRepositoryMySQL{
		db: db,
	}
}

func (h *AdminGenreRepositoryMySQL) GetAll() ([]domain.Genre, error) {
	var genres []domain.Genre
	err := h.db.Model(&domain.Genre{}).
		Find(&genres).
		Error

	return genres, err
}

func (h *AdminGenreRepositoryMySQL) Create(genre *domain.Genre) error {
	return h.db.Create(&genre).Error
}

func (h *AdminGenreRepositoryMySQL) Update(id uint, request ports_admin.UpdateGenreInterface) error {
	updateFields := map[string]any{
		"name": request.Name,
		"slug": request.Slug,
	}

	if err := h.db.Model(&domain.Genre{}).Where("id = ?", id).Updates(updateFields).Error; err != nil {
		return fmt.Errorf("failed to update genre: %+s", err.Error())
	}

	return nil
}

func (h *AdminGenreRepositoryMySQL) Delete(id uint) error {
	if err := h.db.Delete(&domain.Genre{}, id).Error; err != nil {
		return errors.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
