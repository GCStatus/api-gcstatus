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
	userService   *usecases.UserService
	gameService   *usecases.GameService
	bannerService *usecases.BannerService
}

func NewHomeHandler(
	userService *usecases.UserService,
	gameService *usecases.GameService,
	bannerService *usecases.BannerService,
) *HomeHandler {
	return &HomeHandler{
		userService:   userService,
		gameService:   gameService,
		bannerService: bannerService,
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

	banners, err := h.bannerService.GetBannersForHome()
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch banners: "+err.Error())
		return
	}

	var transformedNextGreatRelease any
	if nextGreatRelease != nil {
		transformedNextGreatRelease = resources.TransformGame(*nextGreatRelease, s3.GlobalS3Client, userID)
	} else {
		transformedNextGreatRelease = nil
	}

	transformedBanners := resources.TransformBanners(banners, s3.GlobalS3Client, userID)
	transformedHotGames := resources.TransformGames(hotGames, s3.GlobalS3Client, userID)
	transformedPopularGames := resources.TransformGames(popularGames, s3.GlobalS3Client, userID)
	transformedUpcomingGames := resources.TransformGames(upcomingGames, s3.GlobalS3Client, userID)
	transformedMostLikedGames := resources.TransformGames(mostHeartedGames, s3.GlobalS3Client, userID)

	response := resources.Response{
		Data: map[string]any{
			"banners":          transformedBanners,
			"hot":              transformedHotGames,
			"popular":          transformedPopularGames,
			"upcoming_games":   transformedUpcomingGames,
			"most_liked_games": transformedMostLikedGames,
			"next_release":     transformedNextGreatRelease,
		},
	}

	c.JSON(http.StatusOK, response)
}
