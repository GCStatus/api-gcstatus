package routes

import (
	"gcstatus/internal/middlewares"
	"gcstatus/internal/usecases"

	"github.com/gin-gonic/gin"
)

func RegisterProtectedRoutes(
	r *gin.RouterGroup,
	userService *usecases.UserService,
	handlers *Handlers,
) {
	r.Use(middlewares.JWTAuthMiddleware(userService))

	r.GET("/me", handlers.AuthHandler.Me)
	r.GET("/levels", handlers.LevelHandler.GetAll)
	r.POST("/auth/logout", handlers.AuthHandler.Logout)

	r.GET("/titles", handlers.TitleHandler.GetAllForUser)
	r.PUT("/titles/:id/toggle", handlers.TitleHandler.ToggleEnableTitle)
	r.POST("/titles/:id/buy", handlers.TitleHandler.BuyTitle)

	r.PUT("/profile/password", handlers.PasswordResetHandler.ResetPasswordProfile)
	r.PUT("/profile/picture", handlers.ProfileHandler.UpdatePicture)
	r.PUT("/profile/socials", handlers.ProfileHandler.UpdateSocials)

	r.PUT("/user/update/basics", handlers.UserHandler.UpdateUserBasics)
	r.PUT("/user/update/sensitive", handlers.UserHandler.UpdateUserNickAndEmail)

	r.GET("/transactions", handlers.TransactionHandler.GetAllForUser)

	r.GET("/notifications", handlers.NotificationHandler.GetAllForUser)
	r.PUT("/notifications/:id/read", handlers.NotificationHandler.MarkAsRead)
	r.PUT("/notifications/:id/unread", handlers.NotificationHandler.MarkAsUnread)
	r.PUT("/notifications/all/read", handlers.NotificationHandler.MarkAllAsRead)
	r.PUT("/notifications/all/unread", handlers.NotificationHandler.MarkAllAsUnread)
	r.DELETE("/notifications/:id", handlers.NotificationHandler.DeleteNotification)
	r.DELETE("/notifications/all", handlers.NotificationHandler.DeleteAllNotifications)

	r.GET("/missions", handlers.MissionHandler.GetAllForUser)
	r.POST("/missions/:id/complete", handlers.MissionHandler.CompleteMission)

	r.POST("/hearts", handlers.HeartHandler.ToggleHeartable)
}
