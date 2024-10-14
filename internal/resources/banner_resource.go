package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/s3"
)

type BannerResource struct {
	ID             uint          `json:"id"`
	BannerableType string        `json:"bannerable_type"`
	Game           *GameResource `json:"game,omitempty"`
}

func TransformBanner(banner domain.Banner, s3Client s3.S3ClientInterface, userID uint) BannerResource {
	resource := BannerResource{
		ID:             banner.ID,
		BannerableType: banner.BannerableType,
	}

	switch banner.BannerableType {
	case "games":
		if game, ok := banner.Bannerable.(domain.Game); ok {
			gameResource := TransformGame(game, s3Client, userID)
			resource.Game = &gameResource
		}
	}

	return resource
}

func TransformBanners(banners []domain.Banner, s3Client s3.S3ClientInterface, userID uint) []BannerResource {
	resources := []BannerResource{}

	for _, banner := range banners {
		resources = append(resources, TransformBanner(banner, s3Client, userID))
	}

	return resources
}
