package resources

import (
	"context"
	"gcstatus/internal/domain"
	"gcstatus/pkg/s3"
	"log"
	"time"
)

type GalleriableResource struct {
	ID        uint              `json:"id"`
	Path      string            `json:"path"`
	MediaType MediaTypeResource `json:"media_type"`
}

func TransformGalleriable(galleriable domain.Galleriable, s3Client s3.S3ClientInterface) GalleriableResource {
	resource := GalleriableResource{
		ID:        galleriable.ID,
		MediaType: TransformMediaType(galleriable.MediaType),
	}

	if galleriable.S3 {
		url, err := s3Client.GetPresignedURL(context.TODO(), galleriable.Path, time.Hour*3)
		if err != nil {
			log.Printf("Error generating presigned URL: %v", err)
		}

		resource.Path = url
	} else {
		resource.Path = galleriable.Path
	}

	return resource
}
