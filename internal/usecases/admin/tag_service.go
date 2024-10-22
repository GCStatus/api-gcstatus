package usecases_admin

import (
	"gcstatus/internal/domain"
	ports_admin "gcstatus/internal/ports/admin"
	"gcstatus/internal/utils"
)

type AdminTagService struct {
	repo ports_admin.AdminTagRepository
}

func NewAdminTagService(repo ports_admin.AdminTagRepository) *AdminTagService {
	return &AdminTagService{
		repo: repo,
	}
}

func (h *AdminTagService) GetAll() ([]domain.Tag, error) {
	return h.repo.GetAll()
}

func (h *AdminTagService) Create(tag *domain.Tag) error {
	tag.Slug = utils.Slugify(tag.Name)

	return h.repo.Create(tag)
}

func (h *AdminTagService) Update(id uint, request ports_admin.UpdateTagInterface) error {
	request.Slug = utils.Slugify(request.Name)

	return h.repo.Update(id, request)
}

func (h *AdminTagService) Delete(id uint) error {
	return h.repo.Delete(id)
}
