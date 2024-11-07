package api

import (
	"gcstatus/internal/ports"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
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
	var request ports.HeartTogglePayload

	if err := c.ShouldBindJSON(&request); err != nil {
		RespondWithError(c, http.StatusUnprocessableEntity, "Invalid request data")
		return
	}

	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}

	response, httpErr := h.heartService.ToggleHeartable(request.HeartableID, request.HeartableType, user.ID)
	if err != nil {
		RespondWithError(c, httpErr.Code, httpErr.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
