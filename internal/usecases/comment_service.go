package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	"gcstatus/internal/ports"
	"net/http"
)

type CommentService struct {
	repo ports.CommentRepository
}

func NewCommentService(repo ports.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (h *CommentService) Create(commentable domain.Commentable) (*domain.Commentable, error) {
	return h.repo.Create(commentable)
}

func (h *CommentService) Delete(id uint, userID uint) error {
	comment, err := h.repo.FindByID(id)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.NewHttpError(http.StatusForbidden, "This comment does not belongs to you user!")
	}

	if err := h.repo.Delete(id); err != nil {
		return err
	}

	return nil
}
