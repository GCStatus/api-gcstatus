package routes

import (
	"gcstatus/internal/adapters/api"
	"gcstatus/internal/usecases"
)

func InitHandlers(
	authService *usecases.AuthService,
	userService *usecases.UserService,
	passwordResetService *usecases.PasswordResetService,
	levelService *usecases.LevelService,
) (
	authHandler *api.AuthHandler,
	passwordResetHandler *api.PasswordResetHandler,
	levelHandler *api.LevelHandler,
) {
	authHandler = api.NewAuthHandler(authService, userService)
	passwordResetHandler = api.NewPasswordResetHandler(passwordResetService, userService)
	levelHandler = api.NewLevelHandler(levelService)

	return authHandler, passwordResetHandler, levelHandler
}
