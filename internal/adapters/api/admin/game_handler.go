package api_admin

import (
	"gcstatus/internal/adapters/api"
	"gcstatus/internal/resources"
	resources_admin "gcstatus/internal/resources/admin"
	usecases_admin "gcstatus/internal/usecases/admin"
	"gcstatus/pkg/s3"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminGameHandler struct {
	gameService *usecases_admin.AdminGameService
}

func NewAdminGameHandler(
	gameService *usecases_admin.AdminGameService,
) *AdminGameHandler {
	return &AdminGameHandler{
		gameService: gameService,
	}
}

func (h *AdminGameHandler) GetAll(c *gin.Context) {
	games, err := h.gameService.GetAll()
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch games: "+err.Error())
		return
	}

	transformedGames := resources_admin.TransformGames(games, s3.GlobalS3Client)

	response := resources.Response{
		Data: transformedGames,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AdminGameHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid game ID: "+err.Error())
		return
	}

	game, err := h.gameService.FindByID(uint(id))
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch game: "+err.Error())
		return
	}

	transformedGame := resources_admin.TransformGame(game, s3.GlobalS3Client)

	response := resources.Response{
		Data: transformedGame,
	}

	c.JSON(http.StatusOK, response)
}
