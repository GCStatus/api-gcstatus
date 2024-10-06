package api

import (
	"gcstatus/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService *usecases.GameService
}

func NewGameHandler(
	gameService *usecases.GameService,
) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

func (h *GameHandler) FindBySlug(c *gin.Context) {
	slug := c.Param("slug")

	game, err := h.gameService.FindBySlug(slug)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to fetch game: "+err.Error())
	}

	c.JSON(http.StatusOK, game)
}
