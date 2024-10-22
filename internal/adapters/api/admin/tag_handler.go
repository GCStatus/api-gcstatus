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

type AdminTagHandler struct {
	tagService *usecases_admin.AdminTagService
}

func NewAdminTagHandler(
	tagService *usecases_admin.AdminTagService,
) *AdminTagHandler {
	return &AdminTagHandler{
		tagService: tagService,
	}
}

func (h *AdminTagHandler) GetAll(c *gin.Context) {
	tags, err := h.tagService.GetAll()
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch tags: "+err.Error())
		return
	}

	var transformedTags []resources.TagResource

	if len(tags) > 0 {
		transformedTags = resources.TransformTags(tags)
	} else {
		transformedTags = []resources.TagResource{}
	}

	response := resources.Response{
		Data: transformedTags,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AdminTagHandler) Create(c *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		api.RespondWithError(c, http.StatusUnprocessableEntity, "Please, provide a tag name.")
		return
	}

	tag := &domain.Tag{
		Name: request.Name,
	}

	if err := h.tagService.Create(tag); err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to create tag: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "The tag was successfully created!"})
}

func (h *AdminTagHandler) Update(c *gin.Context) {
	tagIdStr := c.Param("id")

	tagID, err := strconv.ParseUint(tagIdStr, 10, 32)
	if err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid tag ID: "+err.Error())
		return
	}

	var request ports_admin.UpdateTagInterface

	if err := c.ShouldBindJSON(&request); err != nil {
		api.RespondWithError(c, http.StatusUnprocessableEntity, "Please, provide a tag name.")
		return
	}

	if err := h.tagService.Update(uint(tagID), request); err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to update tag: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The tag was successfully updated!"})
}

func (h *AdminTagHandler) Delete(c *gin.Context) {
	tagIdStr := c.Param("id")

	tagID, err := strconv.ParseUint(tagIdStr, 10, 32)
	if err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid tag ID: "+err.Error())
		return
	}

	if err := h.tagService.Delete(uint(tagID)); err != nil {
		if httpErr, ok := err.(*errors.HttpError); ok {
			api.RespondWithError(c, httpErr.Code, httpErr.Error())
		} else {
			api.RespondWithError(c, http.StatusInternalServerError, "Failed to delete tag: "+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The tag was successfully removed!"})
}
