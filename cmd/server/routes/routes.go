package routes

import (
	"gcstatus/config"
	"gcstatus/internal/middlewares"
	"gcstatus/internal/usecases"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the routes for the API
func SetupRouter(
	userService *usecases.UserService,
	authService *usecases.AuthService,
	passwordResetService *usecases.PasswordResetService,
	levelService *usecases.LevelService,
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
	authHandler, passwordResetHandler, levelHandler := InitHandlers(
		authService,
		userService,
		passwordResetService,
		levelService,
	)

	// Define the middlewares
	r.Use(middlewares.LimitThrottleMiddleware())
	protected := r.Group("/")

	// Define the routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/logout", authHandler.Logout)
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/password/email/send", middlewares.LimitResetRequestMiddleware(), passwordResetHandler.RequestPasswordReset)
		authRoutes.POST("/password/reset", passwordResetHandler.ResetUserPassword)
	}

	// Define the auth protected routes
	protected.Use(middlewares.JWTAuthMiddleware(userService))
	{
		protected.GET("/me", authHandler.Me)
		protected.GET("/levels", levelHandler.GetAll)
		protected.PUT("/profile/password", passwordResetHandler.ResetPasswordProfile)
	}

	// Common routes

	return r
}
