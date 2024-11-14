package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	"gcstatus/internal/ports"
	"gcstatus/internal/resources"
	"net/http"
)

type LevelService struct {
	repo ports.LevelRepository
}

func NewLevelService(repo ports.LevelRepository) *LevelService {
	return &LevelService{repo: repo}
}

func (h *LevelService) GetAll() (resources.Response, *errors.HttpError) {
	levels, err := h.repo.GetAll()
	if err != nil {
		return resources.Response{}, errors.NewHttpError(http.StatusInternalServerError, "Failed to fetch platform levels.")
	}

	var transformedLevels []resources.LevelResource

	if len(levels) > 0 {
		transformedLevels = resources.TransformLevels(levels)
	} else {
		transformedLevels = []resources.LevelResource{}
	}

	return resources.Response{
		Data: transformedLevels,
	}, nil
}

func (h *LevelService) FindById(id uint) (*domain.Level, error) {
	return h.repo.FindById(id)
}

func (h *LevelService) FindByLevel(level uint) (*domain.Level, error) {
	return h.repo.FindByLevel(level)
}
