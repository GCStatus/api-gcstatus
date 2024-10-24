package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
)

type ReviewResource struct {
	ID        uint                `json:"id"`
	Rate      uint                `json:"rate"`
	Review    string              `json:"review"`
	Played    bool                `json:"played"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
	User      MinimalUserResource `json:"user"`
}

func TransformReview(review domain.Reviewable, s3Client s3.S3ClientInterface) ReviewResource {
	reviewResource := ReviewResource{
		ID:        review.ID,
		Rate:      review.Rate,
		Review:    review.Review,
		Played:    review.Played,
		CreatedAt: utils.FormatTimestamp(review.CreatedAt),
		UpdatedAt: utils.FormatTimestamp(review.UpdatedAt),
	}

	if review.User.ID != 0 {
		reviewResource.User = TransformMinimalUser(review.User, s3Client)
	}

	return reviewResource
}

func TransformReviews(reviews []domain.Reviewable, s3Client s3.S3ClientInterface) []ReviewResource {
	var resources []ReviewResource

	for _, review := range reviews {
		resources = append(resources, TransformReview(review, s3Client))
	}

	return resources
}
