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

type AdminCategoryHandler struct {
	categoryService *usecases_admin.AdminCategoryService
}

func NewAdminCategoryHandler(
	categoryService *usecases_admin.AdminCategoryService,
) *AdminCategoryHandler {
	return &AdminCategoryHandler{
		categoryService: categoryService,
	}
}

func (h *AdminCategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.categoryService.GetAll()
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch categories: "+err.Error())
		return
	}

	var transformedCategories []resources.CategoryResource

	if len(categories) > 0 {
		transformedCategories = resources.TransformCategories(categories)
	} else {
		transformedCategories = []resources.CategoryResource{}
	}

	response := resources.Response{
		Data: transformedCategories,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AdminCategoryHandler) Create(c *gin.Context) {
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		api.RespondWithError(c, http.StatusUnprocessableEntity, "Please, provide a category name.")
		return
	}

	category := &domain.Category{
		Name: request.Name,
	}

	if err := h.categoryService.Create(category); err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to create category: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "The category was successfully created!"})
}

func (h *AdminCategoryHandler) Update(c *gin.Context) {
	categoryIdStr := c.Param("id")

	categoryID, err := strconv.ParseUint(categoryIdStr, 10, 32)
	if err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid category ID: "+err.Error())
		return
	}

	var request ports_admin.UpdateCategoryInterface

	if err := c.ShouldBindJSON(&request); err != nil {
		api.RespondWithError(c, http.StatusUnprocessableEntity, "Please, provide a category name.")
		return
	}

	if err := h.categoryService.Update(uint(categoryID), request); err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to update category: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The category was successfully updated!"})
}

func (h *AdminCategoryHandler) Delete(c *gin.Context) {
	categoryIdStr := c.Param("id")

	categoryID, err := strconv.ParseUint(categoryIdStr, 10, 32)
	if err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid category ID: "+err.Error())
		return
	}

	if err := h.categoryService.Delete(uint(categoryID)); err != nil {
		if httpErr, ok := err.(*errors.HttpError); ok {
			api.RespondWithError(c, httpErr.Code, httpErr.Error())
		} else {
			api.RespondWithError(c, http.StatusInternalServerError, "Failed to delete category: "+err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "The category was successfully removed!"})
}
