package resources

import "gcstatus/internal/domain"

type PublisherResource struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Acting bool   `json:"acting"`
}

func TransformPublisher(publisher domain.Publisher) PublisherResource {
	return PublisherResource{
		ID:     publisher.ID,
		Name:   publisher.Name,
		Acting: publisher.Acting,
	}
}

func TransformPublishers(Publishers []domain.Publisher) []PublisherResource {
	var resources []PublisherResource

	for _, Publisher := range Publishers {
		resources = append(resources, TransformPublisher(Publisher))
	}

	return resources
}
