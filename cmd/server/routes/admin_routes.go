package routes

import (
	"gcstatus/internal/middlewares"
	"gcstatus/internal/usecases"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(
	r *gin.RouterGroup,
	userService *usecases.UserService,
	handlers *AdminHandlers,
) {
	permissionMiddleware := middlewares.NewPermissionMiddleware(userService)
	r.Use(middlewares.JWTAuthMiddleware(userService))

	r.POST("/login", handlers.AdminAuthHandler.Login)
	r.GET("/me", handlers.AdminAuthHandler.Me)
	r.POST("/logout", handlers.AdminAuthHandler.Logout)
	r.POST("/steam/register/:appID", permissionMiddleware("create:steam-jobs-games"), handlers.AdminSteamHandler.RegisterByAppID)

	r.GET("/categories", permissionMiddleware("view:categories"), handlers.AdminCategoryHandler.GetAll)
	r.POST("/categories", permissionMiddleware("view:categories", "create:categories"), handlers.AdminCategoryHandler.Create)
	r.PUT("/categories/:id", permissionMiddleware("view:categories", "update:categories"), handlers.AdminCategoryHandler.Update)
	r.DELETE("/categories/:id", permissionMiddleware("view:categories", "delete:categories"), handlers.AdminCategoryHandler.Delete)

	r.GET("/genres", permissionMiddleware("view:genres"), handlers.AdminGenreHandler.GetAll)
	r.POST("/genres", permissionMiddleware("view:genres", "create:genres"), handlers.AdminGenreHandler.Create)
	r.PUT("/genres/:id", permissionMiddleware("view:genres", "update:genres"), handlers.AdminGenreHandler.Update)
	r.DELETE("/genres/:id", permissionMiddleware("view:genres", "delete:genres"), handlers.AdminGenreHandler.Delete)

	r.GET("/platforms", permissionMiddleware("view:platforms"), handlers.AdminPlatformHandler.GetAll)
	r.POST("/platforms", permissionMiddleware("view:platforms", "create:platforms"), handlers.AdminPlatformHandler.Create)
	r.PUT("/platforms/:id", permissionMiddleware("view:platforms", "update:platforms"), handlers.AdminPlatformHandler.Update)
	r.DELETE("/platforms/:id", permissionMiddleware("view:platforms", "delete:platforms"), handlers.AdminPlatformHandler.Delete)

	r.GET("/tags", permissionMiddleware("view:tags"), handlers.AdminTagHandler.GetAll)
	r.POST("/tags", permissionMiddleware("view:tags", "create:tags"), handlers.AdminTagHandler.Create)
	r.PUT("/tags/:id", permissionMiddleware("view:tags", "update:tags"), handlers.AdminTagHandler.Update)
	r.DELETE("/tags/:id", permissionMiddleware("view:tags", "delete:tags"), handlers.AdminTagHandler.Delete)

	r.GET("/games", permissionMiddleware("view:games"), handlers.AdminGameHandler.GetAll)
	r.GET("/games/:id", permissionMiddleware("view:games"), handlers.AdminGameHandler.FindByID)
}
