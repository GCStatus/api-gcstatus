package api

import (
	"gcstatus/config"
	"gcstatus/internal/domain"
	"gcstatus/internal/usecases"
	"gcstatus/pkg/cache"
	"gcstatus/pkg/email"
	"gcstatus/pkg/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PasswordResetHandler struct {
	passwordResetService *usecases.PasswordResetService
	userService          *usecases.UserService
	authService          *usecases.AuthService
}

func NewPasswordResetHandler(
	passwordResetService *usecases.PasswordResetService,
	userService *usecases.UserService,
	authService *usecases.AuthService,
) *PasswordResetHandler {
	return &PasswordResetHandler{
		passwordResetService: passwordResetService,
		userService:          userService,
		authService:          authService,
	}
}

func (h *PasswordResetHandler) RequestPasswordReset(c *gin.Context) {
	var requestPasswordResetData struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&requestPasswordResetData); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Please, provide an email address.")
		return
	}

	_, err := h.userService.FindUserByEmailOrNickname(requestPasswordResetData.Email)
	if err != nil {
		RespondWithError(c, http.StatusNotFound, "We could not find an user with that email. Please, double check it and try again!")
		return
	}

	token, err := utils.GenerateResetToken()
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Unable to generate a reset token. Please, contact the support!")
		return
	}

	expiresAt := time.Now().Add(1 * time.Hour)

	passwordReset := domain.PasswordReset{
		Email:     requestPasswordResetData.Email,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.passwordResetService.CreatePasswordReset(&passwordReset); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Unable to save your reset token. Please, contact the support!")
		return
	}

	if err := email.SendPasswordResetEmail(requestPasswordResetData.Email, token, email.Send); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "We could not send you a reset email. Please, try again or contact the support.")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset link has been sent! Check your mailbox."})
}

func (h *PasswordResetHandler) ResetUserPassword(c *gin.Context) {
	var request struct {
		Email                string `json:"email" binding:"required"`
		Password             string `json:"password" binding:"required"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
		Token                string `json:"token"  binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid payload data. Please, provide a valid payload.")
		return
	}

	passwordReset, err := h.passwordResetService.FindPasswordResetByToken(request.Token)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "We could not find your password reset request. Please, try again.")
		return
	}

	if time.Now().After(passwordReset.ExpiresAt) {
		RespondWithError(c, http.StatusBadRequest, "The provided token has already expired. Please, try again.")
		return
	}

	if passwordReset.Email != request.Email {
		RespondWithError(c, http.StatusBadRequest, "Something wrong happened. Please, try again later.")
		return
	}

	if request.Password != request.PasswordConfirmation {
		RespondWithError(c, http.StatusBadRequest, "The password and password confirmation do not match.")
		return
	}

	user, err := h.userService.FindUserByEmailOrNickname(request.Email)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "We could not find an user with that email. Please, double check it and try again!")
		return
	}

	if !utils.ValidatePassword(request.Password) {
		RespondWithError(c, http.StatusBadRequest, "Password must be at least 8 characters long and include a lowercase letter, an uppercase letter, a number, and a symbol.")
		return
	}

	err = h.userService.UpdateUserPassword(user.ID, request.Password)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to reset password: "+err.Error())
		return
	}

	err = h.passwordResetService.DeletePasswordReset(passwordReset.ID)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}

	if err := email.SendPasswordResetConfirmationEmail(user.Email, user.Name, email.Send); err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Unable to send the email reset confirmation.")
		return
	}

	err = cache.GlobalCache.RemovePasswordThrottleCache(user.Email)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Something went wrong. Please, if this affects you directly, contact support!")
		return
	}

	cache.GlobalCache.RemoveUserFromCache(user.ID)

	c.JSON(http.StatusOK, gin.H{"message": "You password was successfully reseted!"})
}

func (h *PasswordResetHandler) ResetPasswordProfile(c *gin.Context) {
	var request struct {
		Password                string `json:"password" binding:"required"`
		NewPassword             string `json:"new_password" binding:"required"`
		NewPasswordConfirmation string `json:"new_password_confirmation" binding:"required"`
	}

	env := config.LoadConfig()

	if err := c.ShouldBindJSON(&request); err != nil {
		errorMessages := utils.FormatValidationError(err)
		RespondWithError(c, http.StatusBadRequest, "Invalid payload data. Please, provide a valid payload. Errors: "+strings.Join(errorMessages, " "))
		return
	}

	if request.NewPassword != request.NewPasswordConfirmation {
		RespondWithError(c, http.StatusBadRequest, "The password and password confirmation do not match.")
		return
	}

	if !utils.ValidatePassword(request.NewPassword) {
		RespondWithError(c, http.StatusBadRequest, "Password must be at least 8 characters long and include a lowercase letter, an uppercase letter, a number, and a symbol.")
		return
	}

	authUser, err := utils.ExtractAuthenticatedUser(c, h.userService.GetUserByID)
	if err != nil {
		RespondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}

	userID, ok := authUser.(uint)
	if !ok {
		RespondWithError(c, http.StatusInternalServerError, "Invalid user ID format.")
		return
	}

	user, found := cache.GlobalCache.GetUserFromCache(userID)
	if !found {
		user, err = h.userService.GetUserByID(userID)
		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if isEqual, err := utils.IsHashEqualsValue(user.Password, request.Password); err != nil || !isEqual {
		if !isEqual {
			RespondWithError(c, http.StatusBadRequest, "Your current password is invalid.")
			return
		}

		if err != nil {
			RespondWithError(c, http.StatusInternalServerError, "Something went wrong on you password verification. Please, try again.")
			return
		}
	}

	err = h.userService.UpdateUserPassword(user.ID, request.NewPassword)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to reset password: "+err.Error())
		return
	}

	h.authService.ClearAuthCookies(c, env.AccessTokenKey, env.IsAuthKey, env.Domain)

	cache.GlobalCache.RemoveUserFromCache(userID)

	c.JSON(http.StatusOK, gin.H{"message": "Your password was successfully changed!"})
}
