package usecases

import (
	goErr "errors"
	"gcstatus/config"
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	"gcstatus/internal/ports"
	"gcstatus/internal/resources"
	"gcstatus/internal/utils"
	"gcstatus/pkg/s3"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserAuthenticator interface {
	GetUserByID(id uint) (*domain.User, error)
	CreateWithProfile(user *domain.User) error
	AuthenticateUser(identifier, password string) (*domain.User, error)
	FindUserByEmailOrNickname(emailOrNickname string) (*domain.User, error)
}

type AuthService struct {
	repo        ports.AuthRepository
	env         config.Config
	userService UserAuthenticator
}

type LoginPayload struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type RegisterPayload struct {
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"required,email"`
	Nickname             string `json:"nickname" binding:"required"`
	Birthdate            string `json:"birthdate" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

func (h *AuthService) SetUserService(userService *UserService) {
	h.userService = userService
}

func NewAuthService(
	repo ports.AuthRepository,
	env config.Config,
	userService UserAuthenticator,
) *AuthService {
	return &AuthService{
		repo:        repo,
		env:         env,
		userService: userService,
	}
}

func (s *AuthService) Login(c *gin.Context, payload LoginPayload) (resources.Response, *errors.HttpError) {
	errGeneric := errors.NewHttpError(http.StatusInternalServerError, "Failed to authenticate user. Please, try again later.")

	user, err := s.userService.AuthenticateUser(payload.Identifier, payload.Password)
	if err != nil {
		return resources.Response{}, errors.NewHttpError(http.StatusUnauthorized, "Invalid credentials. Please, double check it and try again!")
	}

	if user.Blocked {
		return resources.Response{}, errors.NewHttpError(http.StatusForbidden, "Your user has been blocked on GCStatus platform. If you think this is a mistake, please, contact support.")
	}

	expirationSeconds, err := s.GetExpirationSeconds(s.env.JwtTtl)
	if err != nil {
		log.Printf("failed to get expiration seconds: %+v", err)
		return resources.Response{}, errGeneric
	}

	httpSecure, httpOnly, err := s.GetCookieSettings(s.env.HttpSecure, s.env.HttpOnly)
	if err != nil {
		log.Printf("could not parse cookie settings: %+v", err)
		return resources.Response{}, errGeneric
	}

	tokenString, err := s.CreateJWTToken(user.ID, expirationSeconds)
	if err != nil {
		log.Printf("could not create token: %+v", err)
		return resources.Response{}, errGeneric
	}

	encryptedToken, err := s.EncryptToken(tokenString, s.env.JwtSecret)
	if err != nil {
		log.Printf("encryption error: %+v", err)
		return resources.Response{}, errGeneric
	}

	s.SetAuthCookies(c, s.env.AccessTokenKey, encryptedToken, s.env.IsAuthKey, expirationSeconds, httpSecure, httpOnly, s.env.Domain)

	return resources.Response{
		Data: gin.H{"message": "Logged in successfully"},
	}, nil
}

func (s *AuthService) Register(c *gin.Context, payload RegisterPayload) (resources.Response, *errors.HttpError) {
	errGeneric := errors.NewHttpError(http.StatusInternalServerError, "Failed to create user. Please, try again later.")

	if payload.Password != payload.PasswordConfirmation {
		return resources.Response{}, errors.NewHttpError(http.StatusBadRequest, "Password confirmation does not match.")
	}

	birthdate, err := time.Parse("2006-01-02", payload.Birthdate)
	if err != nil {
		return resources.Response{}, errors.NewHttpError(http.StatusBadRequest, "Invalid birthdate format.")
	}

	if time.Since(birthdate).Hours() < 14*365*24 {
		return resources.Response{}, errors.NewHttpError(http.StatusBadRequest, "You must be at least 14 years old to register.")
	}

	if !utils.ValidatePassword(payload.Password) {
		return resources.Response{}, errors.NewHttpError(http.StatusBadRequest, "Password must be at least 8 characters long and include a lowercase letter, an uppercase letter, a number, and a symbol.")
	}

	existingUserByEmail, err := s.userService.FindUserByEmailOrNickname(payload.Email)
	if err != nil && !goErr.Is(err, gorm.ErrRecordNotFound) {
		return resources.Response{}, errors.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	if existingUserByEmail != nil {
		return resources.Response{}, errors.NewHttpError(http.StatusConflict, "Email already in use.")
	}

	existingUserByNickname, err := s.userService.FindUserByEmailOrNickname(payload.Nickname)
	if err != nil && !goErr.Is(err, gorm.ErrRecordNotFound) {
		return resources.Response{}, errors.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	if existingUserByNickname != nil {
		return resources.Response{}, errors.NewHttpError(http.StatusConflict, "Nickname already in use.")
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return resources.Response{}, errors.NewHttpError(http.StatusInternalServerError, "Error hashing password.")
	}

	user := domain.User{
		Name:      payload.Name,
		Email:     payload.Email,
		Nickname:  payload.Nickname,
		Birthdate: birthdate,
		Password:  string(hashedPassword),
		LevelID:   1,
		Profile:   domain.Profile{Share: false},
		Wallet:    domain.Wallet{Amount: 0},
	}

	if err := s.userService.CreateWithProfile(&user); err != nil {
		log.Printf("failed to create user with profile: %+v", err)
		return resources.Response{}, errors.NewHttpError(http.StatusInternalServerError, "Error creating user.")
	}

	expirationSeconds, err := s.GetExpirationSeconds(s.env.JwtTtl)
	if err != nil {
		log.Printf("failed to get expiration seconds: %+v", err)
		return resources.Response{}, errGeneric
	}

	httpSecure, httpOnly, err := s.GetCookieSettings(s.env.HttpSecure, s.env.HttpOnly)
	if err != nil {
		log.Printf("could not parse cookie settings: %+v", err)
		return resources.Response{}, errGeneric
	}

	tokenString, err := s.CreateJWTToken(user.ID, expirationSeconds)
	if err != nil {
		log.Printf("could not create token: %+v", err)
		return resources.Response{}, errGeneric
	}

	encryptedToken, err := s.EncryptToken(tokenString, s.env.JwtSecret)
	if err != nil {
		log.Printf("encryption error: %+v", err)
		return resources.Response{}, errGeneric
	}

	s.SetAuthCookies(c, s.env.AccessTokenKey, encryptedToken, s.env.IsAuthKey, expirationSeconds, httpSecure, httpOnly, s.env.Domain)

	return resources.Response{
		Data: gin.H{"message": "User registered successfully"},
	}, nil
}

func (s *AuthService) Me(c *gin.Context) (resources.Response, *errors.HttpError) {
	user, err := utils.Auth(c, s.userService.GetUserByID)
	if err != nil {
		log.Printf("failed to authenticate user: %+v", err)
		return resources.Response{}, errors.NewHttpError(http.StatusUnauthorized, "Failed to authenticate user.")
	}

	transformedUser := resources.TransformUser(*user, s3.GlobalS3Client)

	return resources.Response{
		Data: transformedUser,
	}, nil
}

func (s *AuthService) Logout(c *gin.Context) (resources.Response, *errors.HttpError) {
	s.ClearAuthCookies(c, s.env.AccessTokenKey, s.env.IsAuthKey, s.env.Domain)

	return resources.Response{
		Data: gin.H{"message": "Logged out successfully"},
	}, nil
}

func (s *AuthService) GetExpirationSeconds(jwtTtl string) (int, error) {
	jwtTtlDays, err := strconv.Atoi(jwtTtl)
	if err != nil {
		return 0, err
	}

	return jwtTtlDays * 86400, nil // Convert days to seconds
}

func (s *AuthService) GetCookieSettings(httpSecureStr, httpOnlyStr string) (bool, bool, error) {
	httpSecure, err := strconv.ParseBool(httpSecureStr)
	if err != nil {
		return true, true, err
	}

	httpOnly, err := strconv.ParseBool(httpOnlyStr)
	if err != nil {
		return true, true, err
	}

	return httpSecure, httpOnly, nil
}

func (s *AuthService) CreateJWTToken(userID uint, expirationSeconds int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(expirationSeconds) * time.Second).Unix(),
	})

	return token.SignedString(config.JWTSecret)
}

func (s *AuthService) EncryptToken(tokenString, jwtSecret string) (string, error) {
	return utils.Encrypt(tokenString, jwtSecret)
}

func (s *AuthService) SetAuthCookies(c *gin.Context, tokenKey, tokenValue, authKey string, expirationSeconds int, secure, httpOnly bool, domain string) {
	c.SetCookie(tokenKey, tokenValue, expirationSeconds, "/", domain, secure, httpOnly)
	c.SetCookie(authKey, "1", expirationSeconds, "/", domain, secure, false)
}

func (s *AuthService) ClearAuthCookies(c *gin.Context, tokenKey, authKey string, domain string) {
	c.SetCookie(tokenKey, "", -1, "/", domain, false, false)
	c.SetCookie(authKey, "", -1, "/", domain, false, false)
}
