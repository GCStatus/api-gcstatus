package ports

import "gcstatus/internal/domain"

type CommentStorePayload struct {
	ParentID        *uint  `json:"parent_id"`
	Comment         string `json:"comment" binding:"required"`
	CommentableID   uint   `json:"commentable_id" binding:"required"`
	CommentableType string `json:"commentable_type" binding:"required"`
}

type CommentRepository interface {
	FindByID(id uint) (*domain.Commentable, error)
	Create(commentable domain.Commentable) (*domain.Commentable, error)
	Delete(id uint) error
}
