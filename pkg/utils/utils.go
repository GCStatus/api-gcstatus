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
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Encrypts using AES encryption
func Encrypt(encryptable string, key string) (string, error) {
	// Create a new AES cipher using the secret key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Create a GCM mode instance with the block cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a nonce for AES-GCM encryption
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the encryptable
	ciphertext := gcm.Seal(nonce, nonce, []byte(encryptable), nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypts using AES encryption
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

	// Parse and validate the JWT token
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

func IsHashEqualsValue(hash string, value string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value)); err != nil {
		return false, errors.New("failed to compare hash")
	}

	return true, nil
}
