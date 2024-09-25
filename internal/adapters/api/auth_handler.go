package api

import (
	"errors"
	"fmt"
	"gcstatus/config"
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/pkg/s3"
	"gcstatus/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService *usecases.AuthService
	userService *usecases.UserService
}

func NewAuthHandler(authService *usecases.AuthService, userService *usecases.UserService) *AuthHandler {
	return &AuthHandler{authService: authService, userService: userService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginData struct {
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}

	env := config.LoadConfig()

	if err := c.ShouldBindJSON(&loginData); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Please provide valid credentials.")
		return
	}

	user, err := h.userService.AuthenticateUser(loginData.Identifier, loginData.Password)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Invalid credentials. Please try again.")
		return
	}

	if user.Blocked {
		RespondWithError(c, http.StatusForbidden, "You are blocked on GCStatus platform. If you think this is an error, please, contact support!")
		return
	}

	expirationSeconds, err := h.authService.GetExpirationSeconds(env.JwtTtl)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Could not parse token expiration.")
		return
	}

	httpSecure, httpOnly, err := h.authService.GetCookieSettings(env.HttpSecure, env.HttpOnly)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Could not parse cookie settings.")
		return
	}

	tokenString, err := h.authService.CreateJWTToken(user.ID, expirationSeconds)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Could not create token.")
		return
	}

	encryptedToken, err := h.authService.EncryptToken(tokenString, env.JwtSecret)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Encryption error: %v", err))
		return
	}

	h.authService.SetAuthCookies(c, env.AccessTokenKey, encryptedToken, env.IsAuthKey, expirationSeconds, httpSecure, httpOnly, env.Domain)

	c.JSON(http.StatusOK, resources.Response{
		Data: gin.H{"message": "Logged in successfully"},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var registrationData struct {
		Name                 string `json:"name" binding:"required"`
		Email                string `json:"email" binding:"required,email"`
		Nickname             string `json:"nickname" binding:"required"`
		Birthdate            string `json:"birthdate" binding:"required"`
		Password             string `json:"password" binding:"required"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registrationData); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Please, provide a valid data to proceed.")
		return
	}

	if registrationData.Password != registrationData.PasswordConfirmation {
		RespondWithError(c, http.StatusBadRequest, "Password confirmation does not match.")
		return
	}

	birthdate, err := time.Parse("2006-01-02", registrationData.Birthdate)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid birthdate format.")
		return
	}

	if time.Since(birthdate).Hours() < 14*365*24 {
		RespondWithError(c, http.StatusBadRequest, "You must be at least 14 years old to register.")
		return
	}

	if !utils.ValidatePassword(registrationData.Password) {
		RespondWithError(c, http.StatusBadRequest, "Password must be at least 8 characters long and include a lowercase letter, an uppercase letter, a number, and a symbol.")
		return
	}

	existingUserByEmail, err := h.userService.FindUserByEmailOrNickname(registrationData.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if existingUserByEmail != nil {
		RespondWithError(c, http.StatusConflict, "Email already in use.")
		return
	}

	existingUserByNickname, err := h.userService.FindUserByEmailOrNickname(registrationData.Nickname)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if existingUserByNickname != nil {
		RespondWithError(c, http.StatusConflict, "Nickname already in use.")
		return
	}

	hashedPassword, err := utils.HashPassword(registrationData.Password)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Error hashing password.")
		return
	}

	user := domain.User{
		Name:      registrationData.Name,
		Email:     registrationData.Email,
		Nickname:  registrationData.Nickname,
		Birthdate: birthdate,
		Password:  string(hashedPassword),
		LevelID:   1,
		Profile:   domain.Profile{Share: false},
		Wallet:    domain.Wallet{Amount: 0},
	}

	if err := h.userService.CreateWithProfile(&user); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Error creating user.")
		return
	}

	env := config.LoadConfig()

	expirationSeconds, err := h.authService.GetExpirationSeconds(env.JwtTtl)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Could not parse token expiration.")
		return
	}

	httpSecure, httpOnly, err := h.authService.GetCookieSettings(env.HttpSecure, env.HttpOnly)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Could not parse cookie settings.")
		return
	}

	tokenString, err := h.authService.CreateJWTToken(user.ID, expirationSeconds)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Could not create token.")
		return
	}

	encryptedToken, err := utils.Encrypt(tokenString, env.JwtSecret)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Encryption error: %v", err))
		return
	}

	h.authService.SetAuthCookies(c, env.AccessTokenKey, encryptedToken, env.IsAuthKey, expirationSeconds, httpSecure, httpOnly, env.Domain)

	c.JSON(http.StatusOK, resources.Response{
		Data: gin.H{"message": "User registered successfully"},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	env := config.LoadConfig()

	h.authService.ClearAuthCookies(c, env.AccessTokenKey, env.IsAuthKey, env.Domain)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	user, err := utils.Auth(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	transformedUser := resources.TransformUser(*user, s3.GlobalS3Client)

	c.JSON(http.StatusOK, resources.Response{
		Data: transformedUser,
	})
}
