package di

import (
	"gcstatus/config"
	"gcstatus/internal/usecases"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDependencies sets up the database connection, repository, and services
func InitDependencies() (
	*usecases.UserService,
	*usecases.AuthService,
	*usecases.PasswordResetService,
	*gorm.DB,
) {
	// Load config
	cfg := config.LoadConfig()

	// Setup DB connection
	dsn := config.GetDBConnectionURL(cfg)
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the database
	MigrateModels(dbConn)

	// Setup dependencies
	userService, authService, passwordResetService := Setup(dbConn)

	return userService, authService, passwordResetService, dbConn
}
