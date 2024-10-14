package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"gcstatus/config"
	"gcstatus/internal/domain"
	"gcstatus/pkg/cache"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(encryptable string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(encryptable), nil)
	return hex.EncodeToString(ciphertext), nil
}

func Decrypt(decryptable string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	data, err := hex.DecodeString(decryptable)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("malformed ciphertext")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

type UserFetcher func(id uint) (*domain.User, error)

func ExtractAuthenticatedUser(c *gin.Context, fetchUser UserFetcher) (any, error) {
	env := config.LoadConfig()

	encryptedToken, err := c.Cookie(env.AccessTokenKey)
	if err != nil {
		return nil, errors.New("user is not authenticated")
	}

	tokenString, err := Decrypt(encryptedToken, env.JwtSecret)
	if err != nil {
		return nil, errors.New("failed to decrypt token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(env.JwtSecret), nil
	})

	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"]
		if userIDFloat, ok := userID.(float64); ok {
			userIDUint := uint(userIDFloat)

			user, err := fetchUser(userIDUint)
			if err != nil {
				return nil, errors.New("user not found")
			}

			if user.Blocked {
				return nil, errors.New("user is blocked")
			}

			return userIDUint, nil
		}

		return nil, errors.New("invalid user ID format")
	}

	return nil, errors.New("invalid token claims")
}

func GetAuthenticatedUserID(c *gin.Context, fetchUser UserFetcher) *uint {
	env := config.LoadConfig()

	encryptedToken, err := c.Cookie(env.AccessTokenKey)
	if err != nil {
		return nil
	}

	tokenString, err := Decrypt(encryptedToken, env.JwtSecret)
	if err != nil {
		return nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(env.JwtSecret), nil
	})

	if err != nil {
		return nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, exists := claims["user_id"]
		if !exists {
			return nil
		}

		userIDFloat, ok := userID.(float64)
		if !ok {
			return nil
		}
		userIDUint := uint(userIDFloat)

		user, err := fetchUser(userIDUint)
		if err != nil {
			return nil
		}

		if user.Blocked {
			return nil
		}

		return &userIDUint
	}

	return nil
}

func ValidatePassword(password string) bool {
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSymbol := regexp.MustCompile(`[!@#$%^&*()_+]`).MatchString(password)
	isLongEnough := len(password) >= 8

	return hasLower && hasUpper && hasDigit && hasSymbol && isLongEnough
}

func HashPassword(password string) (string, error) {
	if len(password) == 0 || isWhitespace(password) {
		return "", errors.New("no whitespaces allowed")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func isWhitespace(s string) bool {
	for _, r := range s {
		if r != ' ' {
			return false
		}
	}

	return true
}

func VarDump(myVar ...any) {
	fmt.Printf("%+v\n", myVar)
}

func DD(myVar ...any) {
	VarDump(myVar...)
	os.Exit(1)
}

func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func GetFirstAndLastName(fullName string) (string, string) {
	words := strings.Fields(fullName)

	if len(words) == 0 {
		return "", ""
	}

	firstName := words[0]

	var lastName string
	if len(words) > 1 {
		lastName = words[len(words)-1]
	} else {
		lastName = ""
	}

	return firstName, lastName
}

func FormatValidationError(err error) []string {
	if err == nil {
		return []string{}
	}

	var errorMessages []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errorMessage := fieldError.Field() + " is required and cannot be empty."
			errorMessages = append(errorMessages, errorMessage)
		}
	} else {
		errorMessages = append(errorMessages, err.Error())
	}

	return errorMessages
}

func IsHashEqualsValue(hash string, value string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value)); err != nil {
		return false
	}

	return true
}

func Auth(c *gin.Context, fetchUser UserFetcher) (*domain.User, error) {
	authUser, err := ExtractAuthenticatedUser(c, fetchUser)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %v", err)
	}

	userID, ok := authUser.(uint)
	if !ok {
		return nil, fmt.Errorf("invalid user ID format")
	}

	user, found := cache.GlobalCache.GetUserFromCache(userID)
	if found {
		return user, nil
	}

	user, err = fetchUser(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	cache.GlobalCache.SetUserInCache(user)

	return user, nil
}

func NullString(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}

	return *s
}

func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string {
	return &s
}

func UintPtr(u uint) *uint {
	return &u
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func FormatTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05")
}

func NormalizeWhitespace(str string) string {
	re := regexp.MustCompile(`\s+`)

	return strings.TrimSpace(re.ReplaceAllString(str, " "))
}

func Slugify(s string) string {
	s = strings.ToLower(s)

	reg, _ := regexp.Compile(`[^a-z0-9\s]+`)
	s = reg.ReplaceAllString(s, "")

	s = strings.ReplaceAll(s, " ", "-")

	s = strings.Trim(s, "-")

	return s
}
