package api

import (
	"gcstatus/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *usecases.AuthService
	userService *usecases.UserService
}

func NewAuthHandler(
	authService *usecases.AuthService,
	userService *usecases.UserService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request usecases.LoginPayload

	if err := c.ShouldBindJSON(&request); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Please provide a valid data.")
		return
	}

	response, err := h.authService.Login(c, request)
	if err != nil {
		RespondWithError(c, err.Code, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request usecases.RegisterPayload

	if err := c.ShouldBindJSON(&request); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Please, provide a valid data to proceed.")
		return
	}

	response, err := h.authService.Register(c, request)
	if err != nil {
		RespondWithError(c, err.Code, err.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response, err := h.authService.Logout(c)
	if err != nil {
		RespondWithError(c, err.Code, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Me(c *gin.Context) {
	response, err := h.authService.Me(c)
	if err != nil {
		RespondWithError(c, err.Code, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
