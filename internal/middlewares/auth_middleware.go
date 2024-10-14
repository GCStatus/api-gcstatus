package middlewares

import (
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware checks if a valid JWT token is present in the request cookies
func JWTAuthMiddleware(userService *usecases.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.ExtractAuthenticatedUser(c, userService.GetUserByID)
		if err != nil {
			if err.Error() == "user is blocked" {
				c.JSON(http.StatusForbidden, gin.H{"message": "Your account is blocked. Please, contact support."})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			}

			c.Abort()
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
