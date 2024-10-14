package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
)

type CommentableResource struct {
	ID          uint                  `json:"id"`
	Comment     string                `json:"comment"`
	IsHearted   bool                  `json:"is_hearted"`
	HeartsCount uint                  `json:"hearts_count"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
	By          MinimalUserResource   `json:"by"`
	Replies     []CommentableResource `json:"replies"`
}

func TransformCommentable(commentable domain.Commentable, s3Client s3.S3ClientInterface, userID uint) CommentableResource {
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
		replies[i] = TransformCommentable(reply, s3Client, userID)
	}

	resource.Replies = replies

	heartsMap := make(map[uint]bool)
	for _, heart := range commentable.Hearts {
		heartsMap[heart.UserID] = true
	}

	resource.IsHearted = heartsMap[userID]

	return resource
}
