package api

import "github.com/gin-gonic/gin"

// Helper to respond with error
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"message": message})
}
