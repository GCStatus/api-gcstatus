package middlewares

import (
	"gcstatus/internal/adapters/api"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserServiceInterface interface {
	GetUserByIDForAdmin(userID uint) (*domain.User, error)
}

func NewPermissionMiddleware(userService UserServiceInterface) func(requiredScopes ...string) gin.HandlerFunc {
	return func(requiredScopes ...string) gin.HandlerFunc {
		return func(c *gin.Context) {
			user, err := utils.Auth(c, userService.GetUserByIDForAdmin)
			if err != nil {
				api.RespondWithError(c, http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}

			hasFullAccess := false
			for _, role := range user.Roles {
				if role.Role.Name == "Technology" {
					hasFullAccess = true
					break
				}
			}

			if hasFullAccess {
				c.Next()
				return
			}

			userPermissions := collectPermissions(user)

			for _, requiredScope := range requiredScopes {
				if !userHasPermission(userPermissions, requiredScope) {
					api.RespondWithError(c, http.StatusForbidden, "insufficient permissions")
					c.Abort()
					return
				}
			}

			c.Next()
		}
	}
}

func collectPermissions(user *domain.User) map[string]bool {
	permissions := make(map[string]bool)

	for _, userPerm := range user.Permissions {
		permissions[userPerm.Permission.Scope] = true
	}

	for _, roleable := range user.Roles {
		for _, rolePerm := range roleable.Role.Permissions {
			permissions[rolePerm.Permission.Scope] = true
		}
	}

	return permissions
}

func userHasPermission(userPermissions map[string]bool, requiredScope string) bool {
	_, hasPermission := userPermissions[requiredScope]
	return hasPermission
}
