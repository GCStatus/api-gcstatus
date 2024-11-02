package routes

import (
	"gcstatus/config"
	"gcstatus/internal/middlewares"
	"gcstatus/internal/usecases"
	usecases_admin "gcstatus/internal/usecases/admin"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(
	userService *usecases.UserService,
	authService *usecases.AuthService,
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
	adminPlatformService *usecases_admin.AdminPlatformService,
	adminTagService *usecases_admin.AdminTagService,
	adminGameService *usecases_admin.AdminGameService,
	heartService *usecases.HeartService,
	commentService *usecases.CommentService,
	db *gorm.DB,
) *gin.Engine {
	r := gin.Default()
	env := config.LoadConfig()

	corsDomains := []string{}
	for _, domain := range strings.Split(env.CorsDomains, ",") {
		corsDomains = append(corsDomains, strings.TrimSpace(domain))
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     corsDomains,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "User-Agent", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowWildcard:    true,
		AllowCredentials: true,
	}))

	handlers, adminHandlers := InitHandlers(
		authService,
		userService,
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
		db,
	)

	r.Use(middlewares.LimitThrottleMiddleware())

	RegisterCommonRoutes(r, handlers)
	RegisterAuthRoutes(r.Group("/auth"), handlers)
	RegisterProtectedRoutes(r.Group("/"), userService, handlers)
	RegisterAdminRoutes(r.Group("/admin"), userService, adminHandlers)

	return r
}
