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

type AdminGenreHandler struct {
	genreService *usecases_admin.AdminGenreService
}

func NewAdminGenreHandler(
	genreService *usecases_admin.AdminGenreService,
) *AdminGenreHandler {
	return &AdminGenreHandler{
		genreService: genreService,
	}
}

func (h *AdminGenreHandler) GetAll(c *gin.Context) {
	genres, err := h.genreService.GetAll()
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch genres: "+err.Error())
		return
	}

	var transformedGenres []resources.GenreResource

	if len(genres) > 0 {
		transformedGenres = resources.TransformGenres(genres)
	} else {
		transformedGenres = []resources.GenreResource{}
	}

	response := resources.Response{
		Data: transformedGenres,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AdminGenreHandler) Create(c *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		api.RespondWithError(c, http.StatusUnprocessableEntity, "Please, provide a genre name.")
		return
	}

	genre := &domain.Genre{
		Name: request.Name,
	}

	if err := h.genreService.Create(genre); err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to create genre: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "The genre was successfully created!"})
}

func (h *AdminGenreHandler) Update(c *gin.Context) {
	genreIdStr := c.Param("id")

	genreID, err := strconv.ParseUint(genreIdStr, 10, 32)
	if err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid genre ID: "+err.Error())
		return
	}

	var request ports_admin.UpdateGenreInterface

	if err := c.ShouldBindJSON(&request); err != nil {
		api.RespondWithError(c, http.StatusUnprocessableEntity, "Please, provide a genre name.")
		return
	}

	if err := h.genreService.Update(uint(genreID), request); err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to update genre: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The genre was successfully updated!"})
}

func (h *AdminGenreHandler) Delete(c *gin.Context) {
	genreIdStr := c.Param("id")

	genreID, err := strconv.ParseUint(genreIdStr, 10, 32)
	if err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid genre ID: "+err.Error())
		return
	}

	if err := h.genreService.Delete(uint(genreID)); err != nil {
		if httpErr, ok := err.(*errors.HttpError); ok {
			api.RespondWithError(c, httpErr.Code, httpErr.Error())
		} else {
			api.RespondWithError(c, http.StatusInternalServerError, "Failed to delete genre: "+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The genre was successfully removed!"})
}
