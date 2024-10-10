package api

import (
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/pkg/s3"
	"gcstatus/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService *usecases.GameService
	userService *usecases.UserService
}

func NewGameHandler(
	gameService *usecases.GameService,
	userService *usecases.UserService,
) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		userService: userService,
	}
}

func (h *GameHandler) FindBySlug(c *gin.Context) {
	slug := c.Param("slug")
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	game, err := h.gameService.FindBySlug(slug, user.ID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch game: "+err.Error())
		return
	}

	transformedGame := resources.TransformGame(game, s3.GlobalS3Client, user.ID)

	response := resources.Response{
		Data: transformedGame,
	}

	c.JSON(http.StatusOK, response)
}
