package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
)

type PublisherResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Acting    bool   `json:"acting"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func TransformPublisher(publisher domain.Publisher) PublisherResource {
	return PublisherResource{
		ID:        publisher.ID,
		Name:      publisher.Name,
		Acting:    publisher.Acting,
		CreatedAt: utils.FormatTimestamp(publisher.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(publisher.UpdatedAt),
	}
}

func TransformPublishers(Publishers []domain.Publisher) []PublisherResource {
	var resources []PublisherResource
	for _, Publisher := range Publishers {
		resources = append(resources, TransformPublisher(Publisher))
	}
	return resources
}
