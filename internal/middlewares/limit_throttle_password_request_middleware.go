package middlewares

import (
	"context"
	"gcstatus/pkg/cache"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	oneMinute = 1 * time.Minute // Limit time window
)

// LimitResetRequestMiddleware is middleware that allows one request per minute
func LimitResetRequestMiddleware() gin.HandlerFunc {
	var ctx = context.Background()

	return func(c *gin.Context) {
		var input struct {
			Email string `json:"email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			c.Abort()
			return
		}

		emailKey := "password-reset:" + input.Email

		_, err := cache.GetPasswordThrottleCache(ctx, emailKey)
		if err == redis.Nil {
			cache.SetPasswordThrottleCache(ctx, emailKey, oneMinute)

			c.Next()
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "You must wait for 60 seconds before sending the email again!"})
			c.Abort()
		}
	}
}
