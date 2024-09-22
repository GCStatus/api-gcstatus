package routes

import (
	"gcstatus/internal/middlewares"
	"gcstatus/internal/usecases"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the routes for the API
func SetupRouter(
	userService *usecases.UserService,
	authService *usecases.AuthService,
	passwordResetService *usecases.PasswordResetService,
) *gin.Engine {
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*.gcstatus.tech", "https://gcstatus.tech"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "User-Agent", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowWildcard:    true,
		AllowCredentials: true,
	}))

	// Initialize the handlers
	authHandler, passwordResetHandler := InitHandlers(authService, userService, passwordResetService)

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
	}

	// Common routes

	return r
}
