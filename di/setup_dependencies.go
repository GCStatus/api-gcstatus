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
) {
	// Create repository instances
	userRepo := db.NewUserRepositoryMySQL(dbConn)
	passwordResetRepo := db.NewPasswordResetRepositoryMySQL(dbConn)

	// Create service instances
	userService := usecases.NewUserService(userRepo)
	authService := usecases.NewAuthService(nil)
	passwordResetService := usecases.NewPasswordResetService(passwordResetRepo)

	return userService, authService, passwordResetService
}
