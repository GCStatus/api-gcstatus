package routes

import (
	"gcstatus/internal/adapters/api"
	api_admin "gcstatus/internal/adapters/api/admin"
	"gcstatus/internal/usecases"
	usecases_admin "gcstatus/internal/usecases/admin"

	"gorm.io/gorm"
)

type Handlers struct {
	AuthHandler          *api.AuthHandler
	PasswordResetHandler *api.PasswordResetHandler
	LevelHandler         *api.LevelHandler
	ProfileHandler       *api.ProfileHandler
	UserHandler          *api.UserHandler
	TitleHandler         *api.TitleHandler
	TransactionHandler   *api.TransactionHandler
	NotificationHandler  *api.NotificationHandler
	MissionHandler       *api.MissionHandler
	GameHandler          *api.GameHandler
	HomeHandler          *api.HomeHandler
	HeartHandler         *api.HeartHandler
	CommentHandler       *api.CommentHandler
}

type AdminHandlers struct {
	AdminAuthHandler     *api_admin.AuthHandler
	AdminCategoryHandler *api_admin.AdminCategoryHandler
	AdminGenreHandler    *api_admin.AdminGenreHandler
	AdminPlatformHandler *api_admin.AdminPlatformHandler
	AdminTagHandler      *api_admin.AdminTagHandler
	AdminGameHandler     *api_admin.AdminGameHandler
	AdminSteamHandler    *api_admin.SteamHandler
}

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
	adminCategoryService *usecases_admin.AdminCategoryService,
	adminGenreService *usecases_admin.AdminGenreService,
	AdminPlatformService *usecases_admin.AdminPlatformService,
	adminTagService *usecases_admin.AdminTagService,
	adminGameService *usecases_admin.AdminGameService,
	heartService *usecases.HeartService,
	commentService *usecases.CommentService,
	db *gorm.DB,
) (*Handlers, *AdminHandlers) {
	return &Handlers{
			AuthHandler:          api.NewAuthHandler(authService, userService),
			PasswordResetHandler: api.NewPasswordResetHandler(passwordResetService, userService, authService),
			LevelHandler:         api.NewLevelHandler(levelService),
			ProfileHandler:       api.NewProfileHandler(profileService, userService),
			UserHandler:          api.NewUserHandler(userService),
			TitleHandler:         api.NewTitleHandler(titleService, userService, walletService, taskService, transactionService, notificationService),
			TransactionHandler:   api.NewTransactionHandler(transactionService, userService),
			NotificationHandler:  api.NewNotificationHandler(notificationService, userService),
			MissionHandler:       api.NewMissionHandler(missionService, userService),
			GameHandler:          api.NewGameHandler(gameService, userService),
			HomeHandler:          api.NewHomeHandler(userService, gameService, bannerService),
			HeartHandler:         api.NewHeartHandler(userService, heartService),
			CommentHandler:       api.NewCommentHandler(userService, commentService),
		},
		&AdminHandlers{
			AdminAuthHandler:     api_admin.NewAuthHandler(authService, userService),
			AdminCategoryHandler: api_admin.NewAdminCategoryHandler(adminCategoryService),
			AdminGenreHandler:    api_admin.NewAdminGenreHandler(adminGenreService),
			AdminPlatformHandler: api_admin.NewAdminPlatformHandler(AdminPlatformService),
			AdminTagHandler:      api_admin.NewAdminTagHandler(adminTagService),
			AdminGameHandler:     api_admin.NewAdminGameHandler(adminGameService),
			AdminSteamHandler:    api_admin.NewSteamHandler(gameService, db),
		}
}
