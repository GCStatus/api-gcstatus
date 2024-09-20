package routes

import (
	"gcstatus/internal/adapters/api"
	"gcstatus/internal/usecases"
)

func InitHandlers(
	authService *usecases.AuthService,
	userService *usecases.UserService,
	passwordResetService *usecases.PasswordResetService,
) (
	authHandler *api.AuthHandler,
	passwordResetHandler *api.PasswordResetHandler,
) {
	authHandler = api.NewAuthHandler(authService, userService)
	passwordResetHandler = api.NewPasswordResetHandler(passwordResetService, userService)

	return authHandler, passwordResetHandler
}
