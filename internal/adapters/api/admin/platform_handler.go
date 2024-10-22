package api_admin

import (
	"gcstatus/internal/adapters/api"
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	ports_admin "gcstatus/internal/ports/admin"
	"gcstatus/internal/resources"
	usecases_admin "gcstatus/internal/usecases/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminPlatformHandler struct {
	platformService *usecases_admin.AdminPlatformService
}

func NewAdminPlatformHandler(
	platformService *usecases_admin.AdminPlatformService,
) *AdminPlatformHandler {
	return &AdminPlatformHandler{
		platformService: platformService,
	}
}

func (h *AdminPlatformHandler) GetAll(c *gin.Context) {
	platforms, err := h.platformService.GetAll()
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch platforms: "+err.Error())
		return
	}

	var transformedPlatforms []resources.PlatformResource

	if len(platforms) > 0 {
		transformedPlatforms = resources.TransformPlatforms(platforms)
	} else {
		transformedPlatforms = []resources.PlatformResource{}
	}

	response := resources.Response{
		Data: transformedPlatforms,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AdminPlatformHandler) Create(c *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		api.RespondWithError(c, http.StatusUnprocessableEntity, "Please, provide a platform name.")
		return
	}

	platform := &domain.Platform{
		Name: request.Name,
	}

	if err := h.platformService.Create(platform); err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to create platform: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "The platform was successfully created!"})
}

func (h *AdminPlatformHandler) Update(c *gin.Context) {
	platformIdStr := c.Param("id")

	platformID, err := strconv.ParseUint(platformIdStr, 10, 32)
	if err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid platform ID: "+err.Error())
		return
	}

	var request ports_admin.UpdatePlatformInterface

	if err := c.ShouldBindJSON(&request); err != nil {
		api.RespondWithError(c, http.StatusUnprocessableEntity, "Please, provide a platform name.")
		return
	}

	if err := h.platformService.Update(uint(platformID), request); err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to update platform: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The platform was successfully updated!"})
}

func (h *AdminPlatformHandler) Delete(c *gin.Context) {
	platformIdStr := c.Param("id")

	platformID, err := strconv.ParseUint(platformIdStr, 10, 32)
	if err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid platform ID: "+err.Error())
		return
	}

	if err := h.platformService.Delete(uint(platformID)); err != nil {
		if httpErr, ok := err.(*errors.HttpError); ok {
			api.RespondWithError(c, httpErr.Code, httpErr.Error())
		} else {
			api.RespondWithError(c, http.StatusInternalServerError, "Failed to delete platform: "+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The platform was successfully removed!"})
}
