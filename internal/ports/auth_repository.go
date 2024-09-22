// Package ports provides interfaces for repository actions.
package ports

import "github.com/gin-gonic/gin"

type AuthRepository interface {
	Me(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Register(c *gin.Context)
	EncryptToken(tokenString, jwtSecret string) (string, error)
	CreateJWTToken(userID uint, expirationSeconds int) (string, error)
	GetCookieSettings(httpSecureStr, httpOnlyStr string) (bool, bool, error)
	ClearAuthCookies(c *gin.Context, tokenKey, authKey string, secure, httpOnly bool, domain string)
	SetAuthCookies(c *gin.Context, tokenKey, tokenValue, authKey string, expirationSeconds int, secure, httpOnly bool, domain string)
}
