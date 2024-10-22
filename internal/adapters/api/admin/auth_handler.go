package api_admin

import (
	"fmt"
	"gcstatus/config"
	"gcstatus/internal/adapters/api"
	"gcstatus/internal/errors"
	"gcstatus/internal/resources"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *usecases.AuthService
	userService *usecases.UserService
}

func NewAuthHandler(authService *usecases.AuthService, userService *usecases.UserService) *AuthHandler {
	return &AuthHandler{authService: authService, userService: userService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request struct {
		Email    string `json:"identifier" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	env := config.LoadConfig()

	if err := c.ShouldBindJSON(&request); err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Please provide valid credentials.")
		return
	}

	user, err := h.userService.AuthenticateUserForAdmin(request.Email, request.Password)
	if err != nil {
		log.Printf("failed to authenticate user admin: %s", err.Error())
		if httpErr, ok := err.(*errors.HttpError); ok {
			api.RespondWithError(c, httpErr.Code, httpErr.Error())
		} else {
			api.RespondWithError(c, http.StatusInternalServerError, "Failed to authenticate user: "+err.Error())
		}
		return
	}

	expirationSeconds, err := h.authService.GetExpirationSeconds(env.JwtTtl)
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Could not parse token expiration.")
		return
	}

	httpSecure, httpOnly, err := h.authService.GetCookieSettings(env.HttpSecure, env.HttpOnly)
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Could not parse cookie settings.")
		return
	}

	tokenString, err := h.authService.CreateJWTToken(user.ID, expirationSeconds)
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Could not create token.")
		return
	}

	encryptedToken, err := h.authService.EncryptToken(tokenString, env.JwtSecret)
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Encryption error: %v", err))
		return
	}

	h.authService.SetAuthCookies(c, env.AccessTokenKey, encryptedToken, env.IsAuthKey, expirationSeconds, httpSecure, httpOnly, env.Domain)

	c.JSON(http.StatusOK, resources.Response{
		Data: gin.H{"message": "Logged in successfully"},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	env := config.LoadConfig()

	h.authService.ClearAuthCookies(c, env.AccessTokenKey, env.IsAuthKey, env.Domain)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByIDForAdmin)
	if err != nil {
		api.RespondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}

	transformedUser := resources_admin.TransformUser(*user, s3.GlobalS3Client)

	c.JSON(http.StatusOK, resources.Response{
		Data: transformedUser,
	})
}
