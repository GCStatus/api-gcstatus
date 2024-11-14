package di

import (
	"gcstatus/config"
	"gcstatus/internal/adapters/db"
	db_admin "gcstatus/internal/adapters/db/admin"
	"gcstatus/internal/usecases"
	usecases_admin "gcstatus/internal/usecases/admin"

	"gorm.io/gorm"
)

func Setup(dbConn *gorm.DB, env config.Config) (
	*usecases.UserService,
	*usecases.AuthService,
	*usecases.PasswordResetService,
	*usecases.LevelService,
	*usecases.ProfileService,
	*usecases.TitleService,
	*usecases.TaskService,
	*usecases.WalletService,
	*usecases.TransactionService,
	*usecases.NotificationService,
	*usecases.MissionService,
	*usecases.GameService,
	*usecases.BannerService,
	*usecases_admin.AdminCategoryService,
	*usecases_admin.AdminGenreService,
	*usecases_admin.AdminPlatformService,
	*usecases_admin.AdminTagService,
	*usecases_admin.AdminGameService,
	*usecases.HeartService,
	*usecases.CommentService,
) {
	// Create repository instances
	userRepo := db.NewUserRepositoryMySQL(dbConn)
	passwordResetRepo := db.NewPasswordResetRepositoryMySQL(dbConn)
	levelRepo := db.NewLevelRepositoryMySQL(dbConn)
	profileRepo := db.NewProfileRepositoryMySQL(dbConn)
	titleRepo := db.NewTitleRepositoryMySQL(dbConn)
	taskRepo := db.NewTaskRepositoryMySQL(dbConn)
	walletRepo := db.NewWalletRepositoryMySQL(dbConn)
	transactionRepo := db.NewTransactionRepositoryMySQL(dbConn)
	notificationRepo := db.NewNotificationRepositoryMySQL(dbConn)
	missionRepo := db.NewMissionRepositoryMySQL(dbConn)
	gameRepo := db.NewGameRepositoryMySQL(dbConn)
	bannerRepo := db.NewBannerRepositoryMySQL(dbConn)
	adminCategoryRepo := db_admin.NewAdminCategoryRepositoryMySQL(dbConn)
	adminGenreRepo := db_admin.NewAdminGenreRepositoryMySQL(dbConn)
	adminPlatformRepo := db_admin.NewAdminPlatformRepositoryMySQL(dbConn)
	adminTagRepo := db_admin.NewAdminTagRepositoryMySQL(dbConn)
	adminGameRepo := db_admin.NewAdminGameRepositoryMySQL(dbConn)
	heartRepo := db.NewHeartRepositoryMySQL(dbConn)
	commentRepo := db.NewCommentRepositoryMySQL(dbConn)

	// Create service instances
	userService := usecases.NewUserService(userRepo)
	authService := usecases.NewAuthService(env, nil)
	passwordResetService := usecases.NewPasswordResetService(passwordResetRepo)
	levelService := usecases.NewLevelService(levelRepo)
	profileService := usecases.NewProfileService(profileRepo)
	titleService := usecases.NewTitleService(titleRepo)
	taskService := usecases.NewTaskService(taskRepo)
	walletService := usecases.NewWalletService(walletRepo)
	transactionService := usecases.NewTransactionService(transactionRepo)
	notificationService := usecases.NewNotificationService(notificationRepo)
	missionService := usecases.NewMissionService(missionRepo)
	gameService := usecases.NewGameService(gameRepo)
	bannerService := usecases.NewBannerService(bannerRepo)
	adminCategoryService := usecases_admin.NewAdminCategoryService(adminCategoryRepo)
	adminGenreService := usecases_admin.NewAdminGenreService(adminGenreRepo)
	adminPlatformService := usecases_admin.NewAdminPlatformService(adminPlatformRepo)
	adminTagService := usecases_admin.NewAdminTagService(adminTagRepo)
	adminGameService := usecases_admin.NewAdminGameService(adminGameRepo)
	heartService := usecases.NewHeartService(heartRepo)
	commentService := usecases.NewCommentService(commentRepo)

	// Set dependencies that require service instances to avoid circular references
	authService.SetUserService(userService)

	return userService,
		authService,
		passwordResetService,
		levelService,
		profileService,
		titleService,
		taskService,
		walletService,
		transactionService,
		notificationService,
		missionService,
		gameService,
		bannerService,
		adminCategoryService,
		adminGenreService,
		adminPlatformService,
		adminTagService,
		adminGameService,
		heartService,
		commentService
}
