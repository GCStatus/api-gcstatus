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

// SetupRouter initializes the routes for the API
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
	AdminPlatformService *usecases_admin.AdminPlatformService,
	AdminTagService *usecases_admin.AdminTagService,
	db *gorm.DB,
) *gin.Engine {
	r := gin.Default()
	env := config.LoadConfig()

	corsDomains := []string{}
	for _, domain := range strings.Split(env.CorsDomains, ",") {
		corsDomains = append(corsDomains, strings.TrimSpace(domain))
	}

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     corsDomains,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "User-Agent", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowWildcard:    true,
		AllowCredentials: true,
	}))

	// Initialize the handlers
	authHandler,
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
		adminAuthHandler,
		adminCategoryHandler,
		adminGenreHandler,
		adminPlatformHandler,
		adminTagHandler := InitHandlers(
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
		AdminPlatformService,
		AdminTagService,
		db,
	)

	// Define the middlewares
	r.Use(middlewares.LimitThrottleMiddleware())
	permissionMiddleware := middlewares.NewPermissionMiddleware(userService)
	protected := r.Group("/")
	admin := r.Group("/admin")

	// Define the routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/password/email/send", middlewares.LimitResetRequestMiddleware(), passwordResetHandler.RequestPasswordReset)
		authRoutes.POST("/password/reset", passwordResetHandler.ResetUserPassword)
	}

	// Define the auth protected routes
	protected.Use(middlewares.JWTAuthMiddleware(userService))
	{
		protected.GET("/me", authHandler.Me)
		protected.GET("/levels", levelHandler.GetAll)
		protected.POST("/auth/logout", authHandler.Logout)

		protected.GET("/titles", titleHandler.GetAllForUser)
		protected.PUT("/titles/:id/toggle", titleHandler.ToggleEnableTitle)
		protected.POST("/titles/:id/buy", titleHandler.BuyTitle)

		protected.PUT("/profile/password", passwordResetHandler.ResetPasswordProfile)
		protected.PUT("/profile/picture", profileHandler.UpdatePicture)
		protected.PUT("/profile/socials", profileHandler.UpdateSocials)

		protected.PUT("/user/update/basics", userHandler.UpdateUserBasics)
		protected.PUT("/user/update/sensitive", userHandler.UpdateUserNickAndEmail)

		protected.GET("/transactions", transactionHandler.GetAllForUser)

		protected.GET("/notifications", notificationHandler.GetAllForUser)
		protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
		protected.PUT("/notifications/:id/unread", notificationHandler.MarkAsUnread)
		protected.PUT("/notifications/all/read", notificationHandler.MarkAllAsRead)
		protected.PUT("/notifications/all/unread", notificationHandler.MarkAllAsUnread)
		protected.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
		protected.DELETE("/notifications/all", notificationHandler.DeleteAllNotifications)

		protected.GET("/missions", missionHandler.GetAllForUser)
		protected.POST("/missions/:id/complete", missionHandler.CompleteMission)
	}

	admin.POST("/login", adminAuthHandler.Login)
	admin.Use(middlewares.JWTAuthMiddleware(userService))
	{
		admin.GET("/me", adminAuthHandler.Me)
		admin.POST("/logout", adminAuthHandler.Logout)
		admin.POST("/steam/register/:appID", permissionMiddleware("create:steam-jobs-games"), steamHandler.RegisterByAppID)

		admin.GET("/categories", permissionMiddleware("view:categories"), adminCategoryHandler.GetAll)
		admin.POST("/categories", permissionMiddleware("view:categories", "create:categories"), adminCategoryHandler.Create)
		admin.PUT("/categories/:id", permissionMiddleware("view:categories", "update:categories"), adminCategoryHandler.Update)
		admin.DELETE("/categories/:id", permissionMiddleware("view:categories", "delete:categories"), adminCategoryHandler.Delete)

		admin.GET("/genres", permissionMiddleware("view:genres"), adminGenreHandler.GetAll)
		admin.POST("/genres", permissionMiddleware("view:genres", "create:genres"), adminGenreHandler.Create)
		admin.PUT("/genres/:id", permissionMiddleware("view:genres", "update:genres"), adminGenreHandler.Update)
		admin.DELETE("/genres/:id", permissionMiddleware("view:genres", "delete:genres"), adminGenreHandler.Delete)

		admin.GET("/platforms", permissionMiddleware("view:platforms"), adminPlatformHandler.GetAll)
		admin.POST("/platforms", permissionMiddleware("view:platforms", "create:platforms"), adminPlatformHandler.Create)
		admin.PUT("/platforms/:id", permissionMiddleware("view:platforms", "update:platforms"), adminPlatformHandler.Update)
		admin.DELETE("/platforms/:id", permissionMiddleware("view:platforms", "delete:platforms"), adminPlatformHandler.Delete)

		admin.GET("/tags", permissionMiddleware("view:tags"), adminTagHandler.GetAll)
		admin.POST("/tags", permissionMiddleware("view:tags", "create:tags"), adminTagHandler.Create)
		admin.PUT("/tags/:id", permissionMiddleware("view:tags", "update:tags"), adminTagHandler.Update)
		admin.DELETE("/tags/:id", permissionMiddleware("view:tags", "delete:tags"), adminTagHandler.Delete)
	}

	// Common routes
	r.GET("/home", homeHandler.Home)
	r.GET("/games/:slug", gameHandler.FindBySlug)

	return r
}
