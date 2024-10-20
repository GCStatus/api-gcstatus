package ports_admin

import "gcstatus/internal/domain"

type UpdateCategoryInterface struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type AdminCategoryRepository interface {
	GetAll() ([]domain.Category, error)
	Create(category *domain.Category) error
	Update(id uint, request UpdateCategoryInterface) error
	Delete(id uint) error
}
