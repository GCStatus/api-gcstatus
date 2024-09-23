package tests

import (
	"errors"
	"gcstatus/tests"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockAuthRepository struct{}

func (m *MockAuthRepository) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"user": "test-user"})
}

func (m *MockAuthRepository) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged in"})
}

func (m *MockAuthRepository) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (m *MockAuthRepository) Register(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func (m *MockAuthRepository) EncryptToken(tokenString, jwtSecret string) (string, error) {
	if tokenString == "" || jwtSecret == "" {
		return "", errors.New("invalid input")
	}

	return "encryptedToken", nil
}

func (m *MockAuthRepository) CreateJWTToken(userID uint, expirationSeconds int) (string, error) {
	if userID == 0 || expirationSeconds <= 0 {
		return "", errors.New("invalid input")
	}

	return "jwtToken", nil
}

func (m *MockAuthRepository) GetCookieSettings(httpSecureStr, httpOnlyStr string) (bool, bool, error) {
	return true, true, nil
}

func (m *MockAuthRepository) ClearAuthCookies(c *gin.Context, tokenKey, authKey string, domain string) {
	c.SetCookie(tokenKey, "", -1, "/", domain, false, false)
	c.SetCookie(authKey, "", -1, "/", domain, false, false)
}

func (m *MockAuthRepository) SetAuthCookies(c *gin.Context, tokenKey, tokenValue, authKey string, expirationSeconds int, secure, httpOnly bool, domain string) {
	c.SetCookie(tokenKey, tokenValue, expirationSeconds, "/", domain, secure, httpOnly)
	c.SetCookie(authKey, tokenValue, expirationSeconds, "/", domain, secure, httpOnly)
}

func TestMockAuthRepository_Me(t *testing.T) {
	testCases := map[string]struct {
		status       int
		expectedBody string
		expectError  bool
	}{
		"success": {
			status:       http.StatusOK,
			expectedBody: `{"user":"test-user"}`,
			expectError:  false,
		},
		"fail": {
			status:       http.StatusInternalServerError,
			expectedBody: `{"error":"something went wrong"}`,
			expectError:  true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &MockAuthRepository{}
			c, w := tests.SetupGinTestContext(http.MethodGet, "/me", "")

			if tc.expectError {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`{"error":"something went wrong"}`))
				if err != nil {
					t.Fatalf("something went wrong on creating error response: %s", err.Error())
				}
			} else {
				mockRepo.Me(c)
			}

			if w.Code != tc.status {
				t.Fatalf("expected status %d, got %d", tc.status, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Fatalf("expected body %s, got %s", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestMockAuthRepository_Login(t *testing.T) {
	testCases := map[string]struct {
		status       int
		expectedBody string
		expectError  bool
	}{
		"success": {
			status:       http.StatusOK,
			expectedBody: `{"message":"logged in"}`,
			expectError:  false,
		},
		"fail": {
			status:       http.StatusInternalServerError,
			expectedBody: `{"error":"something went wrong"}`,
			expectError:  true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &MockAuthRepository{}
			c, w := tests.SetupGinTestContext(http.MethodPost, "/login", "")

			if tc.expectError {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`{"error":"something went wrong"}`))
				if err != nil {
					t.Fatalf("something went wrong on creating error response: %s", err.Error())
				}
			} else {
				mockRepo.Login(c)
			}

			if w.Code != tc.status {
				t.Fatalf("expected status %d, got %d", tc.status, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Fatalf("expected body %s, got %s", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestMockAuthRepository_EncryptToken(t *testing.T) {
	testCases := map[string]struct {
		baseToken   string
		expectToken string
		secret      string
		expectError bool
	}{
		"valid payload": {
			baseToken:   "myToken",
			expectToken: "encryptedToken",
			secret:      "secret",
			expectError: false,
		},
		"invalid payload": {
			baseToken:   "",
			expectToken: "encryptedToken",
			secret:      "",
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &MockAuthRepository{}
			token, err := mockRepo.EncryptToken(tc.baseToken, tc.secret)

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid input")
				assert.Empty(t, token, "expected no encrypted token in case of error")
			} else {
				assert.NoError(t, err)
				assert.Contains(t, token, tc.expectToken)
			}
		})
	}
}

func TestMockAuthRepository_Logout(t *testing.T) {
	testCases := map[string]struct {
		status       int
		expectedBody string
		expectError  bool
	}{
		"success": {
			status:       http.StatusOK,
			expectedBody: `{"message":"logged out"}`,
			expectError:  false,
		},
		"fail": {
			status:       http.StatusInternalServerError,
			expectedBody: `{"error":"something went wrong"}`,
			expectError:  true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &MockAuthRepository{}
			c, w := tests.SetupGinTestContext(http.MethodPost, "/logout", "")

			if tc.expectError {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`{"error":"something went wrong"}`))
				if err != nil {
					t.Fatalf("something went wrong on creating error response: %s", err.Error())
				}
			} else {
				mockRepo.Logout(c)
			}

			if w.Code != tc.status {
				t.Fatalf("expected status %d, got %d", tc.status, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Fatalf("expected body %s, got %s", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestMockAuthRepository_Register(t *testing.T) {
	testCases := map[string]struct {
		status       int
		expectedBody string
		expectError  bool
	}{
		"success": {
			status:       http.StatusCreated,
			expectedBody: `{"message":"user registered"}`,
			expectError:  false,
		},
		"fail": {
			status:       http.StatusInternalServerError,
			expectedBody: `{"error":"something went wrong"}`,
			expectError:  true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &MockAuthRepository{}
			c, w := tests.SetupGinTestContext(http.MethodPost, "/register", "")

			if tc.expectError {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`{"error":"something went wrong"}`))
				if err != nil {
					t.Fatalf("something went wrong on creating error response: %s", err.Error())
				}
			} else {
				mockRepo.Register(c)
			}

			if w.Code != tc.status {
				t.Fatalf("expected status %d, got %d", tc.status, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Fatalf("expected body %s, got %s", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestMockAuthRepository_CreateJWTToken(t *testing.T) {
	testCases := map[string]struct {
		userID            uint
		expirationSeconds uint
		expectError       bool
		expectToken       string
	}{
		"Valid payload": {
			userID:            1,
			expirationSeconds: 14400,
			expectError:       false,
			expectToken:       "jwtToken",
		},
		"Invalid payload": {
			userID:            0,
			expirationSeconds: 0,
			expectError:       true,
			expectToken:       "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &MockAuthRepository{}

			token, err := mockRepo.CreateJWTToken(tc.userID, int(tc.expirationSeconds))

			if tc.expectError {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), "invalid input")
					assert.Empty(t, token, "expected no token in case of error")
				}
			} else {
				assert.NoError(t, err)
				assert.Contains(t, token, tc.expectToken)
			}
		})
	}
}

func TestMockAuthRepository_SetAuthCookies(t *testing.T) {
	testCases := map[string]struct {
		tokenKey          string
		tokenValue        string
		authKey           string
		expirationSeconds int
		secure            bool
		httpOnly          bool
		domain            string
		expectedCookies   int
	}{
		"set cookies success": {
			tokenKey:          "token",
			tokenValue:        "jwtToken",
			authKey:           "auth",
			expirationSeconds: 3600,
			secure:            true,
			httpOnly:          true,
			domain:            "example.com",
			expectedCookies:   2,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &MockAuthRepository{}
			c, w := tests.SetupGinTestContext(http.MethodPost, "/", "")

			mockRepo.SetAuthCookies(c, tc.tokenKey, tc.tokenValue, tc.authKey, tc.expirationSeconds, tc.secure, tc.httpOnly, tc.domain)

			cookies := w.Result().Cookies()

			if len(cookies) != tc.expectedCookies {
				t.Fatalf("expected %d cookies, got %d", tc.expectedCookies, len(cookies))
			}

			tokenCookie := cookies[0]
			authCookie := cookies[1]

			if tokenCookie.Name != tc.tokenKey || tokenCookie.Value != tc.tokenValue {
				t.Errorf("expected token cookie %s=%s, got %s=%s", tc.tokenKey, tc.tokenValue, tokenCookie.Name, tokenCookie.Value)
			}

			if authCookie.Name != tc.authKey || authCookie.Value != tc.tokenValue {
				t.Errorf("expected auth cookie %s=%s, got %s=%s", tc.authKey, tc.tokenValue, authCookie.Name, authCookie.Value)
			}

			if tokenCookie.Secure != tc.secure || authCookie.HttpOnly != tc.httpOnly {
				t.Errorf("expected secure=%v, httpOnly=%v; got secure=%v, httpOnly=%v", tc.secure, tc.httpOnly, tokenCookie.Secure, authCookie.HttpOnly)
			}
		})
	}
}

func TestMockAuthRepository_ClearAuthCookies(t *testing.T) {
	testCases := map[string]struct {
		tokenKey        string
		authKey         string
		secure          bool
		httpOnly        bool
		domain          string
		expectedCookies int
	}{
		"clear cookies success": {
			tokenKey:        "token",
			authKey:         "auth",
			secure:          true,
			httpOnly:        true,
			domain:          "example.com",
			expectedCookies: 2,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &MockAuthRepository{}
			c, w := tests.SetupGinTestContext(http.MethodPost, "/", "")

			mockRepo.ClearAuthCookies(c, tc.tokenKey, tc.authKey, tc.domain)

			cookies := w.Result().Cookies()

			if len(cookies) != tc.expectedCookies {
				t.Fatalf("expected %d cookies, got %d", tc.expectedCookies, len(cookies))
			}

			tokenCookie := cookies[0]
			authCookie := cookies[1]

			if tokenCookie.Name != tc.tokenKey || tokenCookie.Value != "" || tokenCookie.MaxAge != -1 {
				t.Errorf("expected token cookie %s to be cleared, got %s=%s (MaxAge=%d)", tc.tokenKey, tokenCookie.Name, tokenCookie.Value, tokenCookie.MaxAge)
			}

			if authCookie.Name != tc.authKey || authCookie.Value != "" || authCookie.MaxAge != -1 {
				t.Errorf("expected auth cookie %s to be cleared, got %s=%s (MaxAge=%d)", tc.authKey, authCookie.Name, authCookie.Value, authCookie.MaxAge)
			}
		})
	}
}
