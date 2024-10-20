package usecases_admin

import (
	"gcstatus/internal/domain"
	ports_admin "gcstatus/internal/ports/admin"
	"gcstatus/internal/utils"
)

type AdminCategoryService struct {
	repo ports_admin.AdminCategoryRepository
}

func NewAdminCategoryService(repo ports_admin.AdminCategoryRepository) *AdminCategoryService {
	return &AdminCategoryService{
		repo: repo,
	}
}

func (h *AdminCategoryService) GetAll() ([]domain.Category, error) {
	return h.repo.GetAll()
}

func (h *AdminCategoryService) Create(category *domain.Category) error {
	category.Slug = utils.Slugify(category.Name)

	return h.repo.Create(category)
}

func (h *AdminCategoryService) Update(id uint, request ports_admin.UpdateCategoryInterface) error {
	request.Slug = utils.Slugify(request.Name)

	return h.repo.Update(id, request)
}

func (h *AdminCategoryService) Delete(id uint) error {
	return h.repo.Delete(id)
}
