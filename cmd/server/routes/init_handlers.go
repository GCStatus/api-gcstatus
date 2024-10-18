package routes

import (
	"gcstatus/internal/adapters/api"
	api_admin "gcstatus/internal/adapters/api/admin"
	"gcstatus/internal/usecases"

	"gorm.io/gorm"
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
	bannerService *usecases.BannerService,
	db *gorm.DB,
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
	homeHandler *api.HomeHandler,
	steamHandler *api_admin.SteamHandler,
	adminAuthHandler *api_admin.AuthHandler,
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
	gameHandler = api.NewGameHandler(gameService, userService)
	homeHandler = api.NewHomeHandler(userService, gameService, bannerService)
	steamHandler = api_admin.NewSteamHandler(gameService, db)
	adminAuthHandler = api_admin.NewAuthHandler(authService, userService)

	return authHandler,
		passwordResetHandler,
		levelHandler,
		profileHandler,
		userHandler,
		titleHandler,
		transactionHandler,
		notificationHandler,
		missionHandler,
		gameHandler,
		homeHandler,
		steamHandler,
		adminAuthHandler
}
