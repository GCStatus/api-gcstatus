package api

import (
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	userService *usecases.UserService
	gameService *usecases.GameService
}

func NewHomeHandler(
	userService *usecases.UserService,
	gameService *usecases.GameService,
) *HomeHandler {
	return &HomeHandler{
		userService: userService,
		gameService: gameService,
	}
}

func (h *HomeHandler) Home(c *gin.Context) {
	authUserID := utils.GetAuthenticatedUserID(c, h.userService.GetUserByID)

	var userID uint
	if authUserID != nil {
		userID = *authUserID
	}

	hotGames, popularGames, mostHeartedGames, nextGreatRelease, upcomingGames, err := h.gameService.HomeGames()
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch games: "+err.Error())
		return
	}

	transformedHotGames := resources.TransformGames(hotGames, s3.GlobalS3Client, userID)
	transformedPopularGames := resources.TransformGames(popularGames, s3.GlobalS3Client, userID)
	transformedUpcomingGames := resources.TransformGames(upcomingGames, s3.GlobalS3Client, userID)
	transformedMostLikedGames := resources.TransformGames(mostHeartedGames, s3.GlobalS3Client, userID)
	transformedNextGreatRelease := resources.TransformGame(*nextGreatRelease, s3.GlobalS3Client, userID)

	response := resources.Response{
		Data: map[string]any{
			"hot":              transformedHotGames,
			"popular":          transformedPopularGames,
			"most_liked_games": transformedMostLikedGames,
			"next_release":     transformedNextGreatRelease,
			"upcoming_games":   transformedUpcomingGames,
		},
	}

	c.JSON(http.StatusOK, response)
}
