package ports_admin

import "gcstatus/internal/domain"

type UpdateGenreInterface struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type AdminGenreRepository interface {
	GetAll() ([]domain.Genre, error)
	Create(genre *domain.Genre) error
	Update(id uint, request UpdateGenreInterface) error
	Delete(id uint) error
}
