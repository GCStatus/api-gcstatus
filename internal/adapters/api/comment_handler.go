package api

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	userService    *usecases.UserService
	commentService *usecases.CommentService
}

func NewCommentHandler(
	userServuce *usecases.UserService,
	commentService *usecases.CommentService,
) *CommentHandler {
	return &CommentHandler{
		userService:    userServuce,
		commentService: commentService,
	}
}

func (h *CommentHandler) Create(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	var request struct {
		ParentID        *uint  `json:"parent_id"`
		Comment         string `json:"comment" binding:"required"`
		CommentableID   uint   `json:"commentable_id" binding:"required"`
		CommentableType string `json:"commentable_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		RespondWithError(c, http.StatusUnprocessableEntity, "Invalid request data")
		return
	}

	commentable := domain.Commentable{
		UserID:          user.ID,
		Comment:         request.Comment,
		CommentableID:   request.CommentableID,
		CommentableType: request.CommentableType,
		ParentID:        request.ParentID,
	}

	comment, err := h.commentService.Create(commentable)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to create comment.")
		return
	}

	transformedComment := resources.TransformCommentable(*comment, s3.GlobalS3Client, user.ID)

	response := resources.Response{
		Data: transformedComment,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	commentIDStr := c.Param("id")
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid comment ID: "+err.Error())
		return
	}

	if err := h.commentService.Delete(uint(commentID), user.ID); err != nil {
		if httpErr, ok := err.(*errors.HttpError); ok {
			RespondWithError(c, httpErr.Code, httpErr.Error())
		} else {
			RespondWithError(c, http.StatusInternalServerError, "Failed to delete comment: "+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your comment was successfully removed!"})
}
