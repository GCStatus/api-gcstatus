package usecases

import (
	"gcstatus/config"
	"gcstatus/internal/ports"
	"gcstatus/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	repo ports.AuthRepository
}

func NewAuthService(repo ports.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(c *gin.Context) {
	s.repo.Login(c)
}

func (s *AuthService) Register(c *gin.Context) {
	s.repo.Register(c)
}

func (s *AuthService) Me(c *gin.Context) {
	s.repo.Me(c)
}

// Parse expiration time
func (s *AuthService) GetExpirationSeconds(jwtTtl string) (int, error) {
	jwtTtlDays, err := strconv.Atoi(jwtTtl)
	if err != nil {
		return 0, err
	}
	return jwtTtlDays * 86400, nil // Convert days to seconds
}

// Parse cookie security settings
func (s *AuthService) GetCookieSettings(httpSecureStr, httpOnlyStr string) (bool, bool, error) {
	httpSecure, err := strconv.ParseBool(httpSecureStr)
	if err != nil {
		return true, true, err // default to secure and httpOnly cookies
	}
	httpOnly, err := strconv.ParseBool(httpOnlyStr)
	if err != nil {
		return true, true, err // default to secure and httpOnly cookies
	}
	return httpSecure, httpOnly, nil
}

// Handle creating JWT tokens
func (s *AuthService) CreateJWTToken(userID uint, expirationSeconds int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(expirationSeconds) * time.Second).Unix(),
	})

	return token.SignedString(config.JWTSecret)
}

// Encrypt the token string
func (s *AuthService) EncryptToken(tokenString, jwtSecret string) (string, error) {
	return utils.Encrypt(tokenString, jwtSecret)
}

// Set cookies
func (s *AuthService) SetAuthCookies(c *gin.Context, tokenKey, tokenValue, authKey string, expirationSeconds int, secure, httpOnly bool, domain string) {
	// Set JWT token cookie
	c.SetCookie(tokenKey, tokenValue, expirationSeconds, "/", domain, secure, httpOnly)

	// Set "is authenticated" cookie for client
	c.SetCookie(authKey, "1", expirationSeconds, "/", domain, secure, false)
}

// Clear cookies
func (s *AuthService) ClearAuthCookies(c *gin.Context, tokenKey, authKey string, secure, httpOnly bool, domain string) {
	c.SetCookie(tokenKey, "", -1, "/", domain, secure, httpOnly)
	c.SetCookie(authKey, "", -1, "/", domain, secure, false)
}
