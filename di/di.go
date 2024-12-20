package di

import (
	"context"
	"gcstatus/config"
	"gcstatus/internal/usecases"
	usecases_admin "gcstatus/internal/usecases/admin"
	"gcstatus/pkg/cache"
	"gcstatus/pkg/s3"
	"gcstatus/pkg/sqs"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDependencies() (
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
	*gorm.DB,
) {
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
	userService,
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
		commentService := Setup(dbConn)

	// Setup clients for non-test environment
	if cfg.ENV != "testing" {
		sqsClient := sqs.NewSQSClient()
		cache.GlobalCache = cache.NewRedisCache()
		s3.GlobalS3Client = s3.NewS3Client()
		sqs.GlobalSQSClient = sqsClient

		consumer := sqs.NewSQSConsumer(
			sqsClient.GetAWSClient(),
			cfg.AwsSqsUrl,
			userService,
			transactionService,
			notificationService,
			taskService,
			missionService,
			walletService,
		)

		go consumer.Start(context.Background())
	}

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
		commentService,
		dbConn
}
