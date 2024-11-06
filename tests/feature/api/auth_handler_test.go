package feature_tests

import (
	"encoding/json"
	"fmt"
	"gcstatus/config"
	"gcstatus/internal/adapters/api"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	"gcstatus/internal/usecases"
	"gcstatus/internal/utils"
	test_mocks "gcstatus/tests/data/mocks"
	testutils "gcstatus/tests/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var authTruncateModels = []any{
	&domain.User{},
	&domain.Wallet{},
	&domain.Profile{},
}

func setupAuthHandler(dbConn *gorm.DB) *api.AuthHandler {
	userService := usecases.NewUserService(db.NewUserRepositoryMySQL(dbConn))
	authService := usecases.NewAuthService(*config.LoadConfig(), userService)
	return api.NewAuthHandler(authService, userService)
}

func TestAuthHandler_Login(t *testing.T) {
	dummyUser, err := test_mocks.CreateDummyUser(t, dbConn, &domain.User{})
	if err != nil {
		t.Fatalf("failed to create dummy user: %+v", err)
	}

	authHandler := setupAuthHandler(dbConn)

	tests := map[string]struct {
		payload        string
		expectCode     int
		expectResponse string
	}{
		"successful login": {
			payload:        fmt.Sprintf(`{"identifier": "%s", "password": "admin1234"}`, dummyUser.Email),
			expectCode:     200,
			expectResponse: "Logged in successfully",
		},
		"invalid identifier": {
			payload:        `{"identifier": "invalid@example.com", "password": "admin1234"}`,
			expectCode:     401,
			expectResponse: "Invalid credentials. Please, double check it and try again!",
		},
		"invalid password": {
			payload:        fmt.Sprintf(`{"identifier": "%s", "password": "invalidpass"}`, dummyUser.Email),
			expectCode:     401,
			expectResponse: "Invalid credentials. Please, double check it and try again!",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tc.payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			authHandler.Login(c)

			assert.Equal(t, tc.expectCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectResponse)
		})
	}

	t.Cleanup(func() {
		testutils.RefreshDatabase(t, dbConn, authTruncateModels)
	})
}

func TestAuthHandler_Me(t *testing.T) {
	user, err := test_mocks.CreateDummyUser(t, dbConn, &domain.User{
		Email:    "default@example.com",
		Password: "admin1234",
	})
	if err != nil {
		t.Fatalf("failed to create dummy user: %+v", err)
	}

	authHandler := setupAuthHandler(dbConn)

	tests := map[string]struct {
		authToken      string
		expectCode     int
		expectResponse map[string]any
	}{
		"successful authentication": {
			authToken:  testutils.GenerateAuthTokenForUser(t, user),
			expectCode: http.StatusOK,
			expectResponse: map[string]any{
				"id":         float64(user.ID),
				"name":       user.Name,
				"email":      user.Email,
				"level":      float64(user.LevelID),
				"experience": float64(user.Experience),
				"nickname":   user.Nickname,
				"birthdate":  utils.FormatTimestamp(user.Birthdate),
				"created_at": utils.FormatTimestamp(user.CreatedAt),
				"updated_at": utils.FormatTimestamp(user.UpdatedAt),
				"wallet":     nil,
			},
		},
		"missing token": {
			authToken:      "",
			expectCode:     http.StatusUnauthorized,
			expectResponse: map[string]any{"message": "Failed to authenticate user."},
		},
		"invalid token": {
			authToken:      "invalidtoken",
			expectCode:     http.StatusUnauthorized,
			expectResponse: map[string]any{"message": "Failed to authenticate user."},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			if tc.authToken != "" {
				req.AddCookie(&http.Cookie{
					Name:     env.AccessTokenKey,
					Value:    tc.authToken,
					Path:     "/",
					Domain:   env.Domain,
					HttpOnly: true,
					Secure:   false,
					MaxAge:   86400,
				})
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			authHandler.Me(c)

			assert.Equal(t, tc.expectCode, w.Code)

			var responseBody map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("failed to parse JSON response: %+v", err)
			}

			if tc.expectCode == http.StatusOK {
				data, ok := responseBody["data"].(map[string]any)
				if assert.True(t, ok, "response should contain 'data' field") {
					for key, expectedValue := range tc.expectResponse {
						assert.Equal(t, expectedValue, data[key], "unexpected value for '%s'", key)
					}
				}
			} else {
				assert.Contains(t, responseBody, "message")
				assert.Equal(t, tc.expectResponse["message"], responseBody["message"])
			}
		})
	}

	t.Cleanup(func() {
		testutils.RefreshDatabase(t, dbConn, authTruncateModels)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	authHandler := setupAuthHandler(dbConn)

	tests := map[string]struct {
		expectCode     int
		expectResponse string
	}{
		"successful logout": {
			expectCode:     200,
			expectResponse: "Logged out successfully",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/logout", nil)
			req.Header.Set("Content-Type", "application/json")

			_, err := test_mocks.ActingAsDummyUser(t, dbConn, &domain.User{}, req, env)
			if err != nil {
				t.Fatalf("failed to create dummy user: %+v", err)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			authHandler.Logout(c)

			assert.Equal(t, tc.expectCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectResponse)
		})
	}

	t.Cleanup(func() {
		testutils.RefreshDatabase(t, dbConn, authTruncateModels)
	})
}

func TestAuthHandler_Register(t *testing.T) {
	authHandler := setupAuthHandler(dbConn)

	tests := map[string]struct {
		payload        string
		expectCode     int
		expectResponse string
	}{
		"successful registration": {
			payload: `{
				"name": "John Doe",
				"email": "johndoe@example.com",
				"nickname": "johnd",
				"birthdate": "2000-01-01",
				"password": "Password@123",
				"password_confirmation": "Password@123"
			}`,
			expectCode:     http.StatusOK,
			expectResponse: "User registered successfully",
		},
		"password mismatch": {
			payload: `{
				"name": "John Doe",
				"email": "johndoe@example.com",
				"nickname": "johnd",
				"birthdate": "2000-01-01",
				"password": "Password@123",
				"password_confirmation": "Mismatch123"
			}`,
			expectCode:     http.StatusBadRequest,
			expectResponse: "Password confirmation does not match.",
		},
		"invalid birthdate format": {
			payload: `{
				"name": "John Doe",
				"email": "johndoe@example.com",
				"nickname": "johnd",
				"birthdate": "01-01-2000",
				"password": "Password@123",
				"password_confirmation": "Password@123"
			}`,
			expectCode:     http.StatusBadRequest,
			expectResponse: "Invalid birthdate format.",
		},
		"underage user": {
			payload: fmt.Sprintf(`{
				"name": "Young User",
				"email": "younguser@example.com",
				"nickname": "youngie",
				"birthdate": "%s",
				"password": "Password@123",
				"password_confirmation": "Password@123"
			}`, time.Now().Format("2006-01-02")),
			expectCode:     http.StatusBadRequest,
			expectResponse: "You must be at least 14 years old to register.",
		},
		"duplicate email": {
			payload: `{
				"name": "John Doe",
				"email": "existing@example.com",
				"nickname": "newnick",
				"birthdate": "2000-01-01",
				"password": "Password@123",
				"password_confirmation": "Password@123"
			}`,
			expectCode:     http.StatusConflict,
			expectResponse: "Email already in use.",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if name == "duplicate email" {
				_, err := test_mocks.CreateDummyUser(t, dbConn, &domain.User{
					Email: "existing@example.com",
				})
				if err != nil {
					t.Fatalf("failed to create dummy user: %+v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tc.payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			authHandler.Register(c)

			assert.Equal(t, tc.expectCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectResponse)
		})
	}

	t.Cleanup(func() {
		testutils.RefreshDatabase(t, dbConn, authTruncateModels)
	})
}
