package middlewares

import (
	"gcstatus/pkg/cache"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	RateLimit  = 200             // Max requests allowed
	TimeWindow = 1 * time.Minute // Time window for rate limiting
)

func LimitThrottleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate-limit:" + ip

		// Increment the request count and set expiration if not set
		count, err := cache.GlobalCache.AddThrottleCache(key)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			c.Abort()
			return
		}

		// Set expiration for the key if it's the first request
		if count == 1 {
			cache.GlobalCache.ExpireThrottleCache(key, TimeWindow)
		}

		// Check if the user exceeded the rate limit
		if count > int64(RateLimit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests. Please try again later."})
			c.Abort()
			return
		}

		c.Next()
	}
}
