package di

import (
	"gcstatus/config"
	"gcstatus/internal/usecases"
	"gcstatus/pkg/cache"
	"gcstatus/pkg/s3"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDependencies sets up the database connection, repository, and services
func InitDependencies() (
	*usecases.UserService,
	*usecases.AuthService,
	*usecases.PasswordResetService,
	*usecases.LevelService,
	*usecases.ProfileService,
	*usecases.TitleService,
	*usecases.TaskService,
	*usecases.WalletService,
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

	// Setup clients for non-test environment
	if cfg.ENV != "testing" {
		cache.GlobalCache = cache.NewRedisCache()
		s3.GlobalS3Client = s3.NewS3Client()
	}

	// Setup dependencies
	userService,
		authService,
		passwordResetService,
		levelService,
		profileService,
		titleService,
		taskService,
		walletService := Setup(dbConn)

	return userService,
		authService,
		passwordResetService,
		levelService,
		profileService,
		titleService,
		taskService,
		walletService,
		dbConn
}
