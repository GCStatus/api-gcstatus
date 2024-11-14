package api

import (
	"gcstatus/internal/resources"

	"github.com/gin-gonic/gin"
)

// Helper to respond with error
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, resources.Response{
		Data: gin.H{"message": message},
	})
}
