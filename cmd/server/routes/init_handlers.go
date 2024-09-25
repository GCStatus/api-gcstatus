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
	profileService *usecases.ProfileService,
) (
	authHandler *api.AuthHandler,
	passwordResetHandler *api.PasswordResetHandler,
	levelHandler *api.LevelHandler,
	profileHandler *api.ProfileHandler,
	userHandler *api.UserHandler,
) {
	userHandler = api.NewUserHandler(userService)
	authHandler = api.NewAuthHandler(authService, userService)
	passwordResetHandler = api.NewPasswordResetHandler(passwordResetService, userService, authService)
	levelHandler = api.NewLevelHandler(levelService)
	profileHandler = api.NewProfileHandler(profileService, userService)

	return authHandler, passwordResetHandler, levelHandler, profileHandler, userHandler
}
