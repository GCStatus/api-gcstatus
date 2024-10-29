package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
)

type DLCResource struct {
	ID               uint                  `json:"id"`
	Name             string                `json:"name"`
	Cover            string                `json:"cover"`
	About            string                `json:"about"`
	Description      string                `json:"description"`
	ShortDescription string                `json:"short_description"`
	ReleaseDate      string                `json:"release_date"`
	CreatedAt        string                `json:"created_at"`
	UpdatedAt        string                `json:"updated_at"`
	Galleries        []GalleriableResource `json:"galleries"`
	Platforms        []PlatformResource    `json:"platforms"`
	Stores           []DLCStoreResource    `json:"stores"`
}

func TransformDLC(DLC domain.DLC, s3Client s3.S3ClientInterface) DLCResource {
	resource := DLCResource{
		ID:               DLC.ID,
		Name:             DLC.Name,
		Cover:            DLC.Cover,
		About:            DLC.About,
		Description:      DLC.Description,
		ShortDescription: DLC.ShortDescription,
		CreatedAt:        utils.FormatTimestamp(DLC.CreatedAt),
		UpdatedAt:        utils.FormatTimestamp(DLC.UpdatedAt),
		ReleaseDate:      utils.FormatTimestamp(DLC.ReleaseDate),
	}

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
