package usecases_admin

import (
	"gcstatus/internal/domain"
	ports_admin "gcstatus/internal/ports/admin"
	"gcstatus/internal/utils"
)

type AdminGenreService struct {
	repo ports_admin.AdminGenreRepository
}

func NewAdminGenreService(repo ports_admin.AdminGenreRepository) *AdminGenreService {
	return &AdminGenreService{
		repo: repo,
	}
}

func (h *AdminGenreService) GetAll() ([]domain.Genre, error) {
	return h.repo.GetAll()
}

func (h *AdminGenreService) Create(genre *domain.Genre) error {
	genre.Slug = utils.Slugify(genre.Name)

	return h.repo.Create(genre)
}

func (h *AdminGenreService) Update(id uint, request ports_admin.UpdateGenreInterface) error {
	request.Slug = utils.Slugify(request.Name)

	return h.repo.Update(id, request)
}

func (h *AdminGenreService) Delete(id uint) error {
	return h.repo.Delete(id)
}
