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
	titleService *usecases.TitleService,
	taskService *usecases.TaskService,
	walletService *usecases.WalletService,
	transactionService *usecases.TransactionService,
	notificationService *usecases.NotificationService,
	missionService *usecases.MissionService,
	gameService *usecases.GameService,
) (
	authHandler *api.AuthHandler,
	passwordResetHandler *api.PasswordResetHandler,
	levelHandler *api.LevelHandler,
	profileHandler *api.ProfileHandler,
	userHandler *api.UserHandler,
	titleHandler *api.TitleHandler,
	transactionHandler *api.TransactionHandler,
	notificationHandler *api.NotificationHandler,
	missionHandler *api.MissionHandler,
	gameHandler *api.GameHandler,
) {
	userHandler = api.NewUserHandler(userService)
	authHandler = api.NewAuthHandler(authService, userService)
	passwordResetHandler = api.NewPasswordResetHandler(passwordResetService, userService, authService)
	levelHandler = api.NewLevelHandler(levelService)
	profileHandler = api.NewProfileHandler(profileService, userService)
	titleHandler = api.NewTitleHandler(titleService, userService, walletService, taskService, transactionService, notificationService)
	transactionHandler = api.NewTransactionHandler(transactionService, userService)
	notificationHandler = api.NewNotificationHandler(notificationService, userService)
	missionHandler = api.NewMissionHandler(missionService, userService)
	gameHandler = api.NewGameHandler(gameService)

	return authHandler,
		passwordResetHandler,
		levelHandler,
		profileHandler,
		userHandler,
		titleHandler,
		transactionHandler,
		notificationHandler,
		missionHandler,
		gameHandler
}
