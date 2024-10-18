package middlewares

import (
	"gcstatus/internal/adapters/api"
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
				api.RespondWithError(c, http.StatusForbidden, "Your account is blocked. Please, contact support.")
			} else {
				api.RespondWithError(c, http.StatusInternalServerError, err.Error())
			}

			c.Abort()
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
