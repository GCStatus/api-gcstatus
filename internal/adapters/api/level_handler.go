package api

import (
	"gcstatus/internal/resources"
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
	levels, err := h.levelService.GetAll()
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	transformedLevels := resources.TransformLevels(levels)

	response := resources.Response{
		Data: transformedLevels,
	}

	c.JSON(http.StatusOK, response)
}
