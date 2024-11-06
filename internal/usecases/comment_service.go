package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	"gcstatus/internal/ports"
	"gcstatus/internal/resources"
	"gcstatus/pkg/s3"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentService struct {
	repo ports.CommentRepository
}

func NewCommentService(repo ports.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (h *CommentService) Create(user *domain.User, payload ports.CommentStorePayload) (resources.Response, *errors.HttpError) {
	commentable := domain.Commentable{
		UserID:          user.ID,
		Comment:         payload.Comment,
		CommentableID:   payload.CommentableID,
		CommentableType: payload.CommentableType,
		ParentID:        payload.ParentID,
	}

	comment, err := h.repo.Create(commentable)
	if err != nil {
		log.Printf("failed to create comment: %+v.\n err: %+v", commentable, err)
		return resources.Response{}, errors.NewHttpError(http.StatusInternalServerError, "Failed to create comment. Please, try again later.")
	}

	transformedComment := resources.TransformCommentable(*comment, s3.GlobalS3Client, user.ID)

	response := resources.Response{
		Data: transformedComment,
	}

	return response, nil
}

func (h *CommentService) Delete(id uint, userID uint) (resources.Response, *errors.HttpError) {
	comment, err := h.repo.FindByID(id)
	if err != nil {
		return resources.Response{}, errors.NewHttpError(http.StatusNotFound, "Could not found the given comment!")
	}

	if comment.UserID != userID {
		return resources.Response{}, errors.NewHttpError(http.StatusForbidden, "This comment does not belongs to you user!")
	}

	if err := h.repo.Delete(id); err != nil {
		log.Printf("failed to delete comment: %+v.\n err: %+v", comment, err)
		return resources.Response{}, errors.NewHttpError(http.StatusInternalServerError, "We could not delete the given comment. Please, try again later.")
	}

	return resources.Response{
		Data: gin.H{"message": "Your comment was successfully removed!"},
	}, nil
}
