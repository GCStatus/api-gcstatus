package api_admin

import (
	"gcstatus/internal/adapters/api"
	"gcstatus/internal/domain"
	"gcstatus/internal/jobs"
	"gcstatus/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SteamHandler struct {
	gameService *usecases.GameService
	db          *gorm.DB
}

func NewSteamHandler(gameService *usecases.GameService, db *gorm.DB) *SteamHandler {
	return &SteamHandler{gameService: gameService, db: db}
}

func (h *SteamHandler) RegisterByAppID(c *gin.Context) {
	appIDStr := c.Param("appID")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid app ID: "+err.Error())
		return
	}

	exists, err := h.gameService.ExistsForStore(domain.SteamStoreID, uint(appID))
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to check if game already exists: "+err.Error())
		return
	}

	if exists {
		api.RespondWithError(c, http.StatusConflict, "The game you are trying to add already exists from Steam!")
		return
	}

	go jobs.FetchSteamOneByOneApp(h.db, int(appID))

	c.JSON(http.StatusOK, gin.H{"message": "The appID you requested is successfully running in background."})
}
