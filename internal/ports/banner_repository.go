package ports

import "gcstatus/internal/domain"

type BannerRepository interface {
	GetBannersForHome() ([]domain.Banner, error)
}
