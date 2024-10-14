package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
)

type BannerService struct {
	repo ports.BannerRepository
}

func NewBannerService(repo ports.BannerRepository) *BannerService {
	return &BannerService{repo: repo}
}

func (h *BannerService) GetBannersForHome() ([]domain.Banner, error) {
	return h.repo.GetBannersForHome()
}
