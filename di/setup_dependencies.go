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
) {
	// Create repository instances
	userRepo := db.NewUserRepositoryMySQL(dbConn)
	passwordResetRepo := db.NewPasswordResetRepositoryMySQL(dbConn)
	levelRepo := db.NewLevelRepositoryMySQL(dbConn)
	profileRepo := db.NewProfileRepositoryMySQL(dbConn)

	// Create service instances
	userService := usecases.NewUserService(userRepo)
	authService := usecases.NewAuthService(nil)
	passwordResetService := usecases.NewPasswordResetService(passwordResetRepo)
	levelService := usecases.NewLevelService(levelRepo)
	profileService := usecases.NewProfileService(profileRepo)

	return userService, authService, passwordResetService, levelService, profileService
}
