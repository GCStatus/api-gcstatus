package api

import (
	"gcstatus/internal/ports"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	userService    *usecases.UserService
	commentService *usecases.CommentService
}

func NewCommentHandler(
	userService *usecases.UserService,
	commentService *usecases.CommentService,
) *CommentHandler {
	return &CommentHandler{
		userService:    userService,
		commentService: commentService,
	}
}

func (h *CommentHandler) Create(c *gin.Context) {
	var request ports.CommentStorePayload

	if err := c.ShouldBindJSON(&request); err != nil {
		RespondWithError(c, http.StatusUnprocessableEntity, "Invalid request data")
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Failed to create comment: could not authenticate user.")
		return
	}

	response, httpErr := h.commentService.Create(user, request)
	if httpErr != nil {
		RespondWithError(c, httpErr.Code, httpErr.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid comment ID: "+err.Error())
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	response, httpErr := h.commentService.Delete(uint(commentID), user.ID)
	if httpErr != nil {
		RespondWithError(c, httpErr.Code, httpErr.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
