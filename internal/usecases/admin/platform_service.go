package usecases_admin

import (
	"gcstatus/internal/domain"
	ports_admin "gcstatus/internal/ports/admin"
	"gcstatus/internal/utils"
)

type AdminPlatformService struct {
	repo ports_admin.AdminPlatformRepository
}

func NewAdminPlatformService(repo ports_admin.AdminPlatformRepository) *AdminPlatformService {
	return &AdminPlatformService{
		repo: repo,
	}
}

func (h *AdminPlatformService) GetAll() ([]domain.Platform, error) {
	return h.repo.GetAll()
}

func (h *AdminPlatformService) Create(platform *domain.Platform) error {
	platform.Slug = utils.Slugify(platform.Name)

	return h.repo.Create(platform)
}

func (h *AdminPlatformService) Update(id uint, request ports_admin.UpdatePlatformInterface) error {
	request.Slug = utils.Slugify(request.Name)

	return h.repo.Update(id, request)
}

func (h *AdminPlatformService) Delete(id uint) error {
	return h.repo.Delete(id)
}
