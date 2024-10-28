package resources_admin

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
)

type CommentableResource struct {
	ID          uint                  `json:"id"`
	Comment     string                `json:"comment"`
	HeartsCount uint                  `json:"hearts_count"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
	By          MinimalUserResource   `json:"by"`
	Replies     []CommentableResource `json:"replies"`
}

func TransformCommentable(commentable domain.Commentable, s3Client s3.S3ClientInterface) CommentableResource {
	resource := CommentableResource{
		ID:          commentable.ID,
		Comment:     commentable.Comment,
		CreatedAt:   utils.FormatTimestamp(commentable.CreatedAt),
		UpdatedAt:   utils.FormatTimestamp(commentable.UpdatedAt),
		By:          TransformMinimalUser(commentable.User, s3Client),
		HeartsCount: uint(len(commentable.Hearts)),
	}

	replies := make([]CommentableResource, len(commentable.Replies))
	for i, reply := range commentable.Replies {
		replies[i] = TransformCommentable(reply, s3Client)
	}

	resource.Replies = replies

	return resource
}
