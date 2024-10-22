package ports_admin

import "gcstatus/internal/domain"

type UpdatePlatformInterface struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type AdminPlatformRepository interface {
	GetAll() ([]domain.Platform, error)
	Create(platform *domain.Platform) error
	Update(id uint, request UpdatePlatformInterface) error
	Delete(id uint) error
}
