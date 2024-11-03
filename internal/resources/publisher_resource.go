package resources

import "gcstatus/internal/domain"

type PublisherResource struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Acting bool   `json:"acting"`
}

func TransformPublisher(publisher domain.Publisher) PublisherResource {
	return PublisherResource{
		ID:     publisher.ID,
		Name:   publisher.Name,
		Slug:   publisher.Slug,
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
