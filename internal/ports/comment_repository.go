package ports

import "gcstatus/internal/domain"

type CommentRepository interface {
	FindByID(id uint) (*domain.Commentable, error)
	Create(commentable domain.Commentable) (*domain.Commentable, error)
	Delete(id uint) error
}
