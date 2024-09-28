package di

import (
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/usecases"

	"gorm.io/gorm"
)

func Setup(dbConn *gorm.DB) (
	*usecases.UserService,
	*usecases.AuthService,
	*usecases.PasswordResetService,
	*usecases.LevelService,
	*usecases.ProfileService,
	*usecases.TitleService,
	*usecases.TaskService,
	*usecases.WalletService,
) {
	// Create repository instances
	userRepo := db.NewUserRepositoryMySQL(dbConn)
	passwordResetRepo := db.NewPasswordResetRepositoryMySQL(dbConn)
	levelRepo := db.NewLevelRepositoryMySQL(dbConn)
	profileRepo := db.NewProfileRepositoryMySQL(dbConn)
	titleRepo := db.NewTitleRepositoryMySQL(dbConn)
	taskRepo := db.NewTaskRepositoryMySQL(dbConn)
	walletRepo := db.NewWalletRepositoryMySQL(dbConn)

	// Create service instances
	userService := usecases.NewUserService(userRepo)
	authService := usecases.NewAuthService(nil)
	passwordResetService := usecases.NewPasswordResetService(passwordResetRepo)
	levelService := usecases.NewLevelService(levelRepo)
	profileService := usecases.NewProfileService(profileRepo)
	titleService := usecases.NewTitleService(titleRepo)
	taskService := usecases.NewTaskService(taskRepo)
	walletService := usecases.NewWalletService(walletRepo)

	return userService,
		authService,
		passwordResetService,
		levelService,
		profileService,
		titleService,
		taskService,
		walletService
}
