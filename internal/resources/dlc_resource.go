package resources

import (
	"context"
	"gcstatus/internal/domain"
	"gcstatus/pkg/s3"
	"gcstatus/pkg/utils"
	"log"
	"time"
)

type DLCResource struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Cover       string                `json:"cover"`
	ReleaseDate string                `json:"release_date"`
	Galleries   []GalleriableResource `json:"galleries"`
	Platforms   []PlatformResource    `json:"platforms"`
	Stores      []DLCStoreResource    `json:"stores"`
}

func TransformDLC(DLC domain.DLC, s3Client s3.S3ClientInterface) DLCResource {
	resource := DLCResource{
		ID:          DLC.ID,
		Name:        DLC.Name,
		ReleaseDate: utils.FormatTimestamp(DLC.ReleaseDate),
	}

	url, err := s3Client.GetPresignedURL(context.TODO(), DLC.Cover, time.Hour*3)
	if err != nil {
		log.Printf("Error generating presigned URL: %v", err)
	}

	resource.Cover = url
	resource.Platforms = transformPlatforms(DLC.Platforms)
	resource.Galleries = transformGalleries(DLC.Galleries, s3Client)
	resource.Stores = transformDLCStores(DLC.Stores)

	return resource
}

func transformDLCStores(stores []domain.DLCStore) []DLCStoreResource {
	storeResources := make([]DLCStoreResource, 0)
	for _, s := range stores {
		if s.ID != 0 {
			storeResources = append(storeResources, TransformDLCtore(s))
		}
	}

	return storeResources
}
