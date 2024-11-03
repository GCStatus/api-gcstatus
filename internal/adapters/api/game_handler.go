package api

import (
	"gcstatus/internal/errors"
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
	"log"
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

func (h *GameHandler) Search(c *gin.Context) {
	authUserID := utils.GetAuthenticatedUserID(c, h.userService.GetUserByID)

	var userID uint
	if authUserID != nil {
		userID = *authUserID
	}
	searchQuery := c.Query("search")
	if searchQuery == "" {
		RespondWithError(c, http.StatusUnprocessableEntity, "Search query parameter is required")
		return
	}

	games, err := h.gameService.Search(searchQuery)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to search games")
		log.Printf("failed to search games: %+v", err)
		return
	}

	var transformedGames []resources.GameResource
	if len(games) > 0 {
		transformedGames = resources.TransformGames(games, s3.GlobalS3Client, userID)
	} else {
		transformedGames = []resources.GameResource{}
	}

	response := resources.Response{
		Data: transformedGames,
	}

	c.JSON(http.StatusOK, response)
}

func (h *GameHandler) FindByClassification(c *gin.Context) {
	filterable := c.Param("filterable")
	classification := c.Param("classification")
	authUserID := utils.GetAuthenticatedUserID(c, h.userService.GetUserByID)

	var userID uint
	if authUserID != nil {
		userID = *authUserID
	}

	games, err := h.gameService.FindByClassification(classification, filterable)
	if err != nil {
		if httpErr, ok := err.(*errors.HttpError); ok {
			RespondWithError(c, httpErr.Code, httpErr.Error())
		} else {
			RespondWithError(c, http.StatusInternalServerError, "Failed to fetch the games: "+err.Error())
		}
		return
	}

	transformedGames := resources.TransformGames(games, s3.GlobalS3Client, userID)

	response := resources.Response{
		Data: transformedGames,
	}

	c.JSON(http.StatusOK, response)
}
