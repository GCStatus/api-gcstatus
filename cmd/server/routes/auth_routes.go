package routes

import (
	"gcstatus/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	r *gin.RouterGroup,
	handlers *Handlers,
) {
	r.POST("/login", handlers.AuthHandler.Login)
	r.POST("/register", handlers.AuthHandler.Register)
	r.POST("/password/email/send", middlewares.LimitResetRequestMiddleware(), handlers.PasswordResetHandler.RequestPasswordReset)
	r.POST("/password/reset", handlers.PasswordResetHandler.ResetUserPassword)
}
