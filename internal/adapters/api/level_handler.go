package api

import (
	"gcstatus/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LevelHandler struct {
	levelService *usecases.LevelService
}

func NewLevelHandler(levelService *usecases.LevelService) *LevelHandler {
	return &LevelHandler{levelService: levelService}
}

func (h *LevelHandler) GetAll(c *gin.Context) {
	response, err := h.levelService.GetAll()
	if err != nil {
		RespondWithError(c, err.Code, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
