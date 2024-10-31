package api

import (
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HeartHandler struct {
	userService  *usecases.UserService
	heartService *usecases.HeartService
}

func NewHeartHandler(
	userService *usecases.UserService,
	heartService *usecases.HeartService,
) *HeartHandler {
	return &HeartHandler{
		userService:  userService,
		heartService: heartService,
	}
}

func (h *HeartHandler) ToggleHeartable(c *gin.Context) {
	var request struct {
		HeartableID   uint   `json:"heartable_id" binding:"required"`
		HeartableType string `json:"heartable_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}

	if err := h.heartService.ToggleHeartable(request.HeartableID, request.HeartableType, user.ID); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to save heart.")
		log.Printf("failed to save user heart: %+v", err)
		return
	}
}
