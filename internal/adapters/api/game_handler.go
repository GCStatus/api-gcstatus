package api

import (
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
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
	authUserID := utils.GetAuthenticatedUserID(c, h.userService.GetUserByID)

	var userID uint
	if authUserID != nil {
		userID = *authUserID
	}

	game, err := h.gameService.FindBySlug(slug, userID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch game: "+err.Error())
		return
	}

	transformedGame := resources.TransformGame(game, s3.GlobalS3Client, userID)

	response := resources.Response{
		Data: transformedGame,
	}

	c.JSON(http.StatusOK, response)
}
