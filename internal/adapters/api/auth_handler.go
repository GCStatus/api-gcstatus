package api

import (
	"errors"
	"fmt"
	"gcstatus/config"
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/internal/usecases"
	"gcstatus/pkg/cache"
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

	// Validate input
	if err := c.ShouldBindJSON(&loginData); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Please provide valid credentials.")
		return
	}

	// Authenticate user (email or nickname)
	user, err := h.userService.AuthenticateUser(loginData.Identifier, loginData.Password)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, "Invalid credentials. Please try again.")
		return
	}

	if user.Blocked {
		RespondWithError(c, http.StatusForbidden, "You are blocked on GCStatus platform. If you think this is an error, please, contact support!")
		return
	}

	// Parse JWT TTL and other configurations from the service
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

	// Generate JWT token
	tokenString, err := h.authService.CreateJWTToken(user.ID, expirationSeconds)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Could not create token.")
		return
	}

	// Encrypt token string
	encryptedToken, err := h.authService.EncryptToken(tokenString, env.JwtSecret)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Encryption error: %v", err))
		return
	}

	// Set the JWT token and auth cookies
	h.authService.SetAuthCookies(c, env.AccessTokenKey, encryptedToken, env.IsAuthKey, expirationSeconds, httpSecure, httpOnly, env.Domain)

	// Respond with success
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

	// Check if new_password matches password_confirmation
	if registrationData.Password != registrationData.PasswordConfirmation {
		RespondWithError(c, http.StatusBadRequest, "Password confirmation does not match.")
		return
	}

	// Validate birthdate (must be at least 14 years old)
	birthdate, err := time.Parse("2006-01-02", registrationData.Birthdate)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid birthdate format.")
		return
	}

	// Check if the user is at least 14 years old
	if time.Since(birthdate).Hours() < 14*365*24 {
		RespondWithError(c, http.StatusBadRequest, "You must be at least 14 years old to register.")
		return
	}

	// Validate password
	if !utils.ValidatePassword(registrationData.Password) {
		RespondWithError(c, http.StatusBadRequest, "Password must be at least 8 characters long and include a lowercase letter, an uppercase letter, a number, and a symbol.")
		return
	}

	// Check if email already exists
	existingUserByEmail, err := h.userService.FindUserByEmailOrNickname(registrationData.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// If there is an error and it's NOT "record not found", handle it as an actual error
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// If the email exists, respond with a conflict error
	if existingUserByEmail != nil {
		RespondWithError(c, http.StatusConflict, "Email already in use.")
		return
	}

	// Check if nickname already exists
	existingUserByNickname, err := h.userService.FindUserByEmailOrNickname(registrationData.Nickname)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// If there is an error and it's NOT "record not found", handle it as an actual error
		RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// If the nickname exists, respond with a conflict error
	if existingUserByNickname != nil {
		RespondWithError(c, http.StatusConflict, "Nickname already in use.")
		return
	}

	// Proceed with user creation
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
	}

	if err := h.userService.CreateUser(&user); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Error creating user.")
		return
	}

	// Load config
	env := config.LoadConfig()

	// Parse JWT TTL and other configurations
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

	// Generate JWT token
	tokenString, err := h.authService.CreateJWTToken(user.ID, expirationSeconds)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Could not create token.")
		return
	}

	// Encrypt token string
	encryptedToken, err := utils.Encrypt(tokenString, env.JwtSecret)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Encryption error: %v", err))
		return
	}

	// Set the JWT token and auth cookies
	h.authService.SetAuthCookies(c, env.AccessTokenKey, encryptedToken, env.IsAuthKey, expirationSeconds, httpSecure, httpOnly, env.Domain)

	// Respond with success
	c.JSON(http.StatusOK, resources.Response{
		Data: gin.H{"message": "User registered successfully"},
	})
}

// Logout implements the AuthRepository interface
func (h *AuthHandler) Logout(c *gin.Context) {
	env := config.LoadConfig()

	httpSecure, httpOnly, err := h.authService.GetCookieSettings(env.HttpSecure, env.HttpOnly)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Could not parse cookie settings.")
		return
	}

	// Clear the JWT token and auth cookies
	h.authService.ClearAuthCookies(c, env.AccessTokenKey, env.IsAuthKey, httpSecure, httpOnly, env.Domain)

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	authUser, err := utils.ExtractAuthenticatedUser(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Assert that authUser is of type uint
	userID, ok := authUser.(uint)
	if !ok {
		RespondWithError(c, http.StatusInternalServerError, "Invalid user ID format.")
		return
	}

	// Try to get the user from the cache
	user, found := cache.GetUserFromCache(userID)
	if !found {
		// If user is not in the cache, retrieve from the service
		user, err = h.userService.GetUserByID(userID)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		// Store the user in the cache
		cache.SetUserInCache(userID, user)
	}

	transformedUser := resources.TransformUser(*user)

	c.JSON(http.StatusOK, resources.Response{
		Data: transformedUser,
	})
}
