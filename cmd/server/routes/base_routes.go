package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterCommonRoutes(
	r *gin.Engine,
	handlers *Handlers,
) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Everything is ok!"})
	})
	r.GET("/home", handlers.HomeHandler.Home)
	r.GET("/games/search", handlers.GameHandler.Search)
	r.GET("/games/:slug", handlers.GameHandler.FindBySlug)
}
