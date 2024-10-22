package ports_admin

import "gcstatus/internal/domain"

type UpdateTagInterface struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type AdminTagRepository interface {
	GetAll() ([]domain.Tag, error)
	Create(tag *domain.Tag) error
	Update(id uint, request UpdateTagInterface) error
	Delete(id uint) error
}
